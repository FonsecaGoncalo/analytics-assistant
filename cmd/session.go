package cmd

import (
	"analytics/config"
	"analytics/db"
	"analytics/openai"
	"bufio"
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
)

var (
	dbURL       string
	dbUsername  string
	dbPassword  string
	contextFile string
	logSql      bool
)

type MessageData struct {
	Ddl          string
	Query        string
	QueryResults string
	Prompt       string
	Context      string
}

func init() {
	rootCmd.AddCommand(sessionCmd)
	sessionCmd.Flags().StringVarP(&dbURL, "db-url", "u", "", "MySQL database URL")
	sessionCmd.Flags().StringVarP(&dbUsername, "db-username", "n", "", "MySQL database username")
	sessionCmd.Flags().StringVarP(&dbPassword, "db-password", "p", "", "MySQL database password")
	sessionCmd.Flags().StringVarP(&contextFile, "context-file", "c", "", "Path to a file containing business context")
	sessionCmd.Flags().BoolVarP(&logSql, "log-sql", "s", false, "Log SQL")
	sessionCmd.MarkFlagRequired("db-url")
	sessionCmd.MarkFlagRequired("db-username")
	sessionCmd.MarkFlagRequired("db-password")
}

var sessionCmd = &cobra.Command{
	Use:   "session",
	Short: "Initiate a session with the analyst",
	Run: func(cmd *cobra.Command, args []string) {
		//Connect to the MySQL database
		dbc, err := db.Connect(dbURL, dbUsername, dbPassword)
		if err != nil {
			fmt.Println("Error connecting to the MySQL database:", err)
			return
		}

		defer dbc.Close()

		ddl, err := dbc.GetDDL()
		if err != nil {
			fmt.Println("Error getting the database schema:", err)
			return
		}

		systemMessages := renderSystemMessages(config.GetAnalystSystemMessages(), ddl)

		if contextFile != "" {
			c, err := readFileContents(contextFile)
			if err != nil {
				fmt.Println("Error reading file", err)
				return
			}
			context := renderTemplate(config.GetAnalystContextMessages(), &MessageData{Context: c})

			systemMessages = append(systemMessages, context)
		}

		analyst := openai.NewOpenAISession(systemMessages, config.GetAnalystTemperature())
		queryParser := openai.NewOpenAISession(config.GetQueryParserSystemMessages(), config.GetQueryParserTemperature())

		reader := bufio.NewReader(os.Stdin)
		for {
			fmt.Print("🧍 Prompt: ")
			input, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("Error reading input:", err)
				continue
			}

			input = strings.TrimSpace(input)

			if input == "exit" {
				break
			}

			answer, err := handlePrompt(input, analyst, queryParser, dbc)
			if err != nil {
				fmt.Println("Error handling prompt", err)
				return
			}
			fmt.Println("🤖 Analytics assistant: " + answer)
		}
	},
}

func handlePrompt(input string, analyst *openai.Session, queryParser *openai.Session, dbc *db.DBConnection) (string, error) {
	response := analyst.UserPrompt(input)

	m := renderTemplate(config.GetQueryParserMessage(), &MessageData{
		Query: response,
	})

	query := queryParser.UserPrompt(m)
	if query == "No query was found." {
		return response, nil
	}

	queryResult, err := dbc.ExecuteQuery(query, logSql)
	if err != nil {
		//TODO: implement retry
		return "", err
	}

	m = renderTemplate(config.GetAnalystQueryResultsMessage(), &MessageData{
		QueryResults: queryResult,
		Query:        query,
		Prompt:       input,
	})

	return analyst.SystemPrompt(m), nil
}

func renderSystemMessages(messageTemplates []string, ddl string) []string {
	m := make([]string, len(messageTemplates))

	data := &MessageData{
		Ddl: ddl,
	}

	for i := range messageTemplates {
		renderedTemplate := renderTemplate(messageTemplates[i], data)

		m[i] = renderedTemplate
	}

	return m
}

func renderTemplate(tmpl string, data *MessageData) string {
	t, err := template.New("message").Parse(tmpl)
	if err != nil {
		panic("Error parsing template")
	}

	var buf strings.Builder
	err = t.Execute(&buf, data)
	if err != nil {
		panic("Error rendering the template")
	}

	return buf.String()
}

func readFileContents(filename string) (string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

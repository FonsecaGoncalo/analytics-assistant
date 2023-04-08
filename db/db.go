package db

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type DBConnection struct {
	db *sql.DB
}

func Connect(url, username, password string) (*DBConnection, error) {
	dataSourceName := fmt.Sprintf("%s:%s@%s", username, password, url)
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return &DBConnection{db: db}, nil
}

func (dbc *DBConnection) Close() error {
	return dbc.db.Close()
}

func (dbc *DBConnection) GetDDL() (string, error) {
	rows, err := dbc.db.Query("SHOW TABLES")
	if err != nil {
		return "", err
	}
	defer rows.Close()

	var ddl []string
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return "", err
		}

		var createTable string
		row := dbc.db.QueryRow("SHOW CREATE TABLE " + tableName)
		if err := row.Scan(&tableName, &createTable); err != nil {
			return "", err
		}
		ddl = append(ddl, createTable)
	}

	return strings.Join(ddl, "\n\n"), nil
}

func (dbc *DBConnection) ExecuteQuery(query string, logQuery bool) (string, error) {
	if logQuery {
		fmt.Println("QUERY: \n " + query)
	}

	rows, err := dbc.db.Query(query)
	if err != nil {
		return "", err
	}

	queryResult, err := formatQueryResults(rows)
	if err != nil {
		return "", err
	}

	return queryResult, nil
}

func formatQueryResults(rows *sql.Rows) (string, error) {
	columns, err := rows.Columns()
	if err != nil {
		return "", err
	}
	values := make([]interface{}, len(columns))
	valuePtrs := make([]interface{}, len(columns))

	for i := range columns {
		valuePtrs[i] = &values[i]
	}

	var results []map[string]interface{}
	for rows.Next() {
		err := rows.Scan(valuePtrs...)
		if err != nil {
			return "", err
		}

		result := make(map[string]interface{})
		for i, col := range columns {
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				val = string(b)
			}
			result[col] = val
		}
		results = append(results, result)
	}

	var formattedResults strings.Builder
	for _, result := range results {
		for col, val := range result {
			formattedResults.WriteString(fmt.Sprintf("%s: %v\n", col, val))
		}
		formattedResults.WriteString("\n")
	}

	return formattedResults.String(), nil
}

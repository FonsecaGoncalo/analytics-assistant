#!/bin/bash

set -e

TOOL_NAME="analytics-assistant"
VERSION="v0.0.1-alpha"
RELEASE_URL="https://github.com/FonsecaGoncalo/analytics-assistant/releases/download/${VERSION}"

# Get the operating system and architecture information
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)
if [ "${ARCH}" = "x86_64" ]; then
  ARCH="amd64"
elif [ "${ARCH}" = "armv7l" ]; then
  ARCH="arm"
fi

# Download the package
PACKAGE_NAME="${TOOL_NAME}-${VERSION}-${OS}-${ARCH}.tar.gz"
curl -L -o "${PACKAGE_NAME}" "${RELEASE_URL}/${PACKAGE_NAME}"

# Extract the package
tar -xf "${PACKAGE_NAME}"
rm "${PACKAGE_NAME}"

# Move the tool to a system-wide location (assuming /usr/local/bin is in the PATH)
sudo mv "${TOOL_NAME}" "/usr/local/bin/${TOOL_NAME}"

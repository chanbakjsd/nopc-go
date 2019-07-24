#!/bin/bash

if [[ -x "$(command -v antlr4)" ]]; then
	antlr4 -Dlanguage=Go nop.g4 -o ./src/parser/
elif [[ -f /usr/local/bin/antlr4.jar ]]; then
	java -jar /usr/local/bin/antlr4.jar -Dlanguage=Go nop.g4 -o ./src/parser/
else
	echo "You must install antlr4, see README.md"
	exit 1
fi

if ! go build -o bin/nopc src/main.go; then
    echo "Error has occurred"
fi

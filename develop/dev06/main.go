package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

// Реализовать утилиту аналог консольной команды cut (man cut).
// Утилита должна принимать строки через STDIN, разбивать по разделителю (TAB) на колонки и выводить запрошенные.

// Реализовать поддержку утилитой следующих ключей:
// -f - "fields" - выбрать поля (колонки)
// -d - "delimiter" - использовать другой разделитель
// -s - "separated" - только строки с разделителем

type cutParameters struct {
	Fields    string
	Delimiter string
	Separated bool
}

func main() {
	parameters := parceCmdArgs()

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		line := scanner.Text()
		if parameters.Separated && !strings.Contains(line, parameters.Delimiter) {
			continue
		}
		fields := strings.Split(line, parameters.Delimiter)
		outputFields := fieldsSelection(fields, parameters.Fields)
		fmt.Println(strings.Join(outputFields, parameters.Delimiter))
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
		os.Exit(1)
	}
}

func parceCmdArgs() cutParameters {
	parameters := cutParameters{}

	flag.StringVar(&parameters.Fields, "f", "", "Select fields")
	flag.StringVar(&parameters.Delimiter, "d", "\t", "Use a different separator")
	flag.BoolVar(&parameters.Separated, "s", false, "Delimited lines only")

	flag.Parse()

	return parameters
}

func fieldsSelection(fields []string, fieldNumbers string) []string {
	if fieldNumbers == "" {
		return fields
	}

	var outputFields []string
	selectedFields := strings.Split(fieldNumbers, ",")

	for _, field := range selectedFields {
		index := parseIndex(field, len(fields))
		if index != -1 {
			outputFields = append(outputFields, fields[index])
		}
	}

	return outputFields
}

func parseIndex(field string, maxIndex int) int {
	index := parsePosInt(field)
	if index == 0 || index > maxIndex {
		return -1
	}
	return index - 1
}

func parsePosInt(str string) int {
	num := 0
	fmt.Sscanf(str, "%d", &num)
	if num <= 0 {
		return 0
	}
	return num
}

package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
)

// Реализовать утилиту фильтрации по аналогии с консольной утилитой (man grep — смотрим описание и основные параметры).

// Реализовать поддержку утилитой следующих ключей:
// -A - "after" печатать +N строк после совпадения
// -B - "before" печатать +N строк до совпадения
// -C - "context" (A+B) печатать ±N строк вокруг совпадения
// -c - "count" (количество строк)
// -i - "ignore-case" (игнорировать регистр)
// -v - "invert" (вместо совпадения, исключать)
// -F - "fixed", точное совпадение со строкой, не паттерн
// -n - "line num", напечатать номер строки

type grepParameters struct {
	After      int
	Before     int
	Context    int
	Count      bool
	IgnoreCase bool
	Invert     bool
	Fixed      bool
	LineNum    bool
	Pattern    string
	FilePaths  []string
}

func main() {
	options := parceCmdArgs()

	if len(options.FilePaths) == 0 {
		fmt.Println("Specify the path to the input file.")
		return
	}

	for _, filePath := range options.FilePaths {
		lines, err := readLines(filePath)
		if err != nil {
			fmt.Printf("Error reading the file %s: %v\n", filePath, err)
			return
		}

		matchedLines := grepFiltering(lines, options)
		conclusion(matchedLines, options)
	}
}

func parceCmdArgs() grepParameters {
	parameters := grepParameters{}

	flag.IntVar(&parameters.After, "A", 0, "Печатать +N строк после совпадения")
	flag.IntVar(&parameters.Before, "B", 0, "Печатать +N строк до совпадения")
	flag.IntVar(&parameters.Context, "C", 0, "Печатать ±N строк вокруг совпадения")
	flag.BoolVar(&parameters.Count, "c", false, "Количество строк")
	flag.BoolVar(&parameters.IgnoreCase, "i", false, "Игнорировать регистр")
	flag.BoolVar(&parameters.Invert, "v", false, "Вместо совпадения, исключать")
	flag.BoolVar(&parameters.Fixed, "F", false, "Точное совпадение со строкой, не паттерн")
	flag.BoolVar(&parameters.LineNum, "n", false, "Печатать номера строк")

	flag.Parse()

	parameters.FilePaths = flag.Args()

	if len(parameters.FilePaths) > 0 {
		parameters.Pattern = parameters.FilePaths[0]
		parameters.FilePaths = parameters.FilePaths[1:]
	}

	return parameters
}

func readLines(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, scanner.Err()
}

func grepFiltering(lines []string, options grepParameters) []string {
	var result []string
	pattern := options.Pattern

	if options.IgnoreCase {
		pattern = "(?i)" + pattern
	}

	if options.Fixed {
		pattern = regexp.QuoteMeta(pattern)
	}

	re, err := regexp.Compile(pattern)
	if err != nil {
		fmt.Printf("Error in the regular expression: %v\n", err)
		return result
	}

	for i, line := range lines {
		matched := re.MatchString(line)

		if options.Invert {
			matched = !matched
		}

		if matched {
			result = append(result, line)

			if options.After > 0 && i+options.After < len(lines) {
				result = append(result, lines[i+1:i+1+options.After]...)
			}

			if options.Before > 0 && i-options.Before >= 0 {
				result = append(result, lines[i-options.Before:i]...)
			}

			if options.Context > 0 && i-options.Context >= 0 && i+options.Context < len(lines) {
				result = append(result, lines[i-options.Context:i]...)
				result = append(result, lines[i+1:i+1+options.Context]...)
			}
		}
	}

	return result
}

func conclusion(lines []string, options grepParameters) {
	if options.Count {
		fmt.Printf("Number of matches: %d\n", len(lines))
	} else {
		for i, line := range lines {
			if options.LineNum {
				fmt.Printf("%d: ", i+1)
			}

			fmt.Println(line)
		}
	}
}

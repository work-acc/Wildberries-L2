package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

//Отсортировать строки в файле по аналогии с консольной утилитой sort
// (man sort — смотрим описание и основные параметры): на входе подается файл из несортированными строками, на выходе — файл с отсортированными.
// Реализовать поддержку утилитой следующих ключей:

// -k — указание колонки для сортировки (слова в строке могут выступать в качестве колонок, по умолчанию разделитель — пробел)
// -n — сортировать по числовому значению
// -r — сортировать в обратном порядке
// -u — не выводить повторяющиеся строки

// Дополнительно
// Реализовать поддержку утилитой следующих ключей:

// -M — сортировать по названию месяца
// -b — игнорировать хвостовые пробелы
// -c — проверять отсортированы ли данные
// -h — сортировать по числовому значению с учетом суффиксов

type SortParameters struct {
	ColumnKey           int
	Numeric             bool
	Reverse             bool
	Unique              bool
	Month               bool
	IgnoreLeadingBlanks bool
	Check               bool
	NumericSuffix       bool
}

func parceCmdArgs() SortParameters {
	columnKey := flag.Int("k", 0, "Индекс столбца для сортировки (по умолчанию 0)")
	numeric := flag.Bool("n", false, "Сортировать по числовому значению")
	reverse := flag.Bool("r", false, "Сортировать в обратном порядке")
	unique := flag.Bool("u", false, "Убрать повторяющиеся строки")
	month := flag.Bool("M", false, "Сортировать по названию месяца")
	ignoreLeadingBlanks := flag.Bool("b", false, "Игнорировать начальные пробелы")
	check := flag.Bool("c", false, "Проверить, отсортированы ли данные")
	numericSuffix := flag.Bool("h", false, "Сортировка по числовому значению с суффиксами")

	flag.Parse()

	return SortParameters{
		ColumnKey:           *columnKey,
		Numeric:             *numeric,
		Reverse:             *reverse,
		Unique:              *unique,
		Month:               *month,
		IgnoreLeadingBlanks: *ignoreLeadingBlanks,
		Check:               *check,
		NumericSuffix:       *numericSuffix,
	}
}

func main() {
	filePath := flag.String("file", "", "Path to the input file")
	parameters := parceCmdArgs()
	flag.Parse()

	if *filePath == "" {
		fmt.Println("Specify the path to the input file using the flag -file")
		return
	}

	lines, err := readFile(*filePath)
	if err != nil {
		fmt.Println("Error reading the file:", err)
		return
	}

	copyLines := make([]string, len(lines))
	copy(copyLines, lines)

	if parameters.Unique {
		lines = removeDuplicates(lines)
	}

	sort.Sort(sortLines{lines, parameters})

	if parameters.Check && isSorted(copyLines, parameters) {
		fmt.Println("The data is sorted")
		return
	} else if parameters.Check && !isSorted(copyLines, parameters) {
		fmt.Println("The data is not sorted")
		return
	}

	outputFilePath := "sorted_" + *filePath
	err = writeFile(outputFilePath, lines)
	if err != nil {
		fmt.Println("Error writing the file: ", err)
		return
	}

	fmt.Println("The data is written to a file:", outputFilePath)
}

func readFile(filePath string) ([]string, error) {
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

func writeFile(filePath string, lines []string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, line := range lines {
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			return err
		}
	}

	return writer.Flush()
}

type sortLines struct {
	lines      []string
	parameters SortParameters
}

func (s sortLines) Len() int      { return len(s.lines) }
func (s sortLines) Swap(i, j int) { s.lines[i], s.lines[j] = s.lines[j], s.lines[i] }

func (s sortLines) Less(i, j int) bool {
	line1 := s.lines[i]
	line2 := s.lines[j]

	if s.parameters.Numeric {
		num1, err1 := extractNumericValue(line1, s.parameters.ColumnKey, s.parameters.NumericSuffix)
		num2, err2 := extractNumericValue(line2, s.parameters.ColumnKey, s.parameters.NumericSuffix)

		if err1 == nil && err2 == nil {
			if num1 < num2 {
				return !s.parameters.Reverse
			} else if num1 > num2 {
				return s.parameters.Reverse
			}
		}
	}

	if s.parameters.Month {
		month1, err1 := parseMonth(line1, s.parameters.ColumnKey)
		month2, err2 := parseMonth(line2, s.parameters.ColumnKey)

		if err1 == nil && err2 == nil {
			if month1 < month2 {
				return !s.parameters.Reverse
			} else if month1 > month2 {
				return s.parameters.Reverse
			}
		}
	}

	if s.parameters.IgnoreLeadingBlanks {
		line1 = strings.TrimSpace(line1)
		line2 = strings.TrimSpace(line2)
	}

	if !s.parameters.Reverse {
		return line1 < line2
	} else {
		return line1 > line2
	}
}

func extractNumericValue(line string, columnIndex int, numericSuffix bool) (float64, error) {
	fields := strings.Fields(line)
	if columnIndex >= len(fields) {
		return 0, fmt.Errorf("the index of the column is out of range")
	}

	value := fields[columnIndex]
	if numericSuffix {
		return parceNumericSuffix(value)
	}

	return strconv.ParseFloat(value, 64)
}

func parceNumericSuffix(value string) (float64, error) {
	suffixes := map[string]float64{
		"K": 1e3,
		"M": 1e6,
		"G": 1e9,
		"T": 1e12,
	}

	for suffix, multiplier := range suffixes {
		if strings.HasSuffix(value, suffix) {
			numStr := strings.TrimSuffix(value, suffix)
			num, err := strconv.ParseFloat(numStr, 64)
			if err != nil {
				return 0, err
			}
			return num * multiplier, nil
		}
	}

	return strconv.ParseFloat(value, 64)
}

func isSorted(lines []string, parameters SortParameters) bool {
	sorter := sortLines{lines, parameters}
	for i := 1; i < len(lines); i++ {
		if sorter.Less(i, i-1) {
			return false
		}
	}
	return true
}

func removeDuplicates(lines []string) []string {
	seen := make(map[string]struct{})
	result := make([]string, 0, len(lines))

	for _, line := range lines {
		if _, exists := seen[line]; !exists {
			result = append(result, line)
			seen[line] = struct{}{}
		}
	}

	return result
}

func parseMonth(line string, columnIndex int) (time.Month, error) {
	fields := strings.Fields(line)
	if columnIndex >= len(fields) {
		return 0, fmt.Errorf("the index of the column is out of range")
	}

	monthStr := strings.ToLower(fields[columnIndex])
	switch monthStr {
	case "January":
		return time.January, nil
	case "February":
		return time.February, nil
	case "March":
		return time.March, nil
	case "April":
		return time.April, nil
	case "May":
		return time.May, nil
	case "June":
		return time.June, nil
	case "July":
		return time.July, nil
	case "August":
		return time.August, nil
	case "September":
		return time.September, nil
	case "October":
		return time.October, nil
	case "November":
		return time.November, nil
	case "December":
		return time.December, nil
	default:
		return 0, fmt.Errorf("wrong month: %s", monthStr)
	}
}

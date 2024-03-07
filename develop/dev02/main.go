package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

//Создать Go-функцию, осуществляющую примитивную распаковку строки, содержащую повторяющиеся символы/руны, например:
// "a4bc2d5e" => "aaaabccddddde"
// "abcd" => "abcd"
// "45" => "" (некорректная строка)
// "" => ""

// Дополнительно
// Реализовать поддержку escape-последовательностей.
// Например:
// qwe\4\5 => qwe45 (*)
// qwe\45 => qwe44444 (*)
// qwe\\5 => qwe\\\\\ (*)

// В случае если была передана некорректная строка, функция должна возвращать ошибку. Написать unit-тесты.

func main() {
	result, err := unpack("qwe\\4\\5")
	if err != nil {
		fmt.Println(err)

		return
	}

	fmt.Println("result = ", result)
}

func unpack(data string) (string, error) {
	if len(data) == 0 {
		return "", nil
	}

	var ans string
	if exist := strings.Contains(data, "\\"); exist {
		ans, err := unpackEscape(data)
		return ans, err

	}

	ans, err := unpackPrimitive(data)
	if err != nil {
		return "", err
	}

	return ans, nil
}

func unpackPrimitive(data string) (string, error) {
	var result []rune

	for i := 0; i < len(data); i++ {
		if unicode.IsDigit(rune(data[i])) {
			return "", errors.New("invalid data")
		}

		if i < len(data)-1 && unicode.IsDigit(rune(data[i+1])) {
			count, _ := strconv.Atoi(string(data[i+1]))

			for j := 0; j < count; j++ {
				result = append(result, rune(data[i]))
			}

			i++

			continue
		}

		result = append(result, rune(data[i]))
	}

	return string(result), nil
}

func unpackEscape(data string) (string, error) {
	var result []rune

	for i := 0; i < len(data); i++ {

		if i-1 >= 0 && data[i-1] == '\\' && i+1 < len(data) && data[i+1] == '\\' {
			result = append(result, rune(data[i]))
			i++

		} else if i-1 >= 0 && data[i-1] == '\\' && i+1 < len(data) && unicode.IsDigit(rune(data[i+1])) {
			counter, _ := strconv.Atoi(string(data[i+1]))

			for j := 0; j < counter; j++ {
				result = append(result, rune(data[i]))
			}
			i++

		} else if i-1 >= 0 && data[i-1] == '\\' && i+1 == len(data) {
			result = append(result, rune(data[i]))

		} else if data[i] == '\\' {
			continue

		} else if unicode.IsDigit(rune(data[i])) && i-1 >= 0 {
			counter, _ := strconv.Atoi(string(data[i]))

			if unicode.IsDigit(rune(data[i-1])) {
				return "", errors.New("invalid data")
			}

			for j := 0; j < counter-1; j++ {
				result = append(result, rune(data[i]))
			}

		} else {
			result = append(result, rune(data[i]))
		}

	}

	return string(result), nil
}

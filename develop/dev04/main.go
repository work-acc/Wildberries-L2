package main

import (
	"fmt"
	"sort"
	"strings"
)

// Написать функцию поиска всех множеств анаграмм по словарю.

// Например:
// 'пятак', 'пятка' и 'тяпка' - принадлежат одному множеству,
// 'листок', 'слиток' и 'столик' - другому.

// Требования:
// Входные данные для функции: ссылка на массив, каждый элемент которого - слово на русском языке в кодировке utf8
// Выходные данные: ссылка на мапу множеств анаграмм
// Ключ - первое встретившееся в словаре слово из множества. Значение - ссылка на массив, каждый элемент которого,
// слово из множества.
// Массив должен быть отсортирован по возрастанию.
// Множества из одного элемента не должны попасть в результат.
// Все слова должны быть приведены к нижнему регистру.
// В результате каждое слово должно встречаться только один раз.

func main() {
	words := []string{"пятак", "пятка", "тяпка", "листок", "слиток", "столик"}
	result := searchAnagrams(&words)

	for key, value := range result {
		if len(value) > 1 {
			fmt.Printf("A lot of anagrams for %s: %v\n", key, value)
		}
	}
}

func sortString(str string) string {
	sortedRunes := []rune(str)
	sort.Slice(sortedRunes, func(i, j int) bool {
		return sortedRunes[i] < sortedRunes[j]
	})
	return string(sortedRunes)
}

func searchAnagrams(words *[]string) map[string][]string {
	anagrams := make(map[string][]string)

	for _, word := range *words {
		sortedWord := sortString(strings.ToLower(word))

		if set, found := anagrams[sortedWord]; found {
			anagrams[sortedWord] = append(set, word)
		} else {
			anagrams[sortedWord] = []string{word}
		}
	}

	for key, value := range anagrams {
		if len(value) <= 1 {
			delete(anagrams, key)
		} else {
			sort.Strings(anagrams[key])
		}
	}

	return anagrams
}

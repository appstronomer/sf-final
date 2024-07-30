package checker

import "regexp"

var re *regexp.Regexp = regexp.MustCompile(`qwerty|йцукен|zxvbnm`)

// Проверка комментария на вхождение стоп-слов
// Возвращает true, если стоп-слово обнаружено
func CheckIfIncorrect(comment Comment) bool {
	return re.Match([]byte(comment.Content))
}

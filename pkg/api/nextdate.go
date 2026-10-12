package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const dateFormat = "20060102"

func afterNow(date, now time.Time) bool {
	return date.After(now)
}

// NextDate вычисляет следующую дату выполнения задачи
func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if repeat == "" {
		return "", fmt.Errorf("не указано правило повторения")
	}

	date, err := time.Parse(dateFormat, dstart)
	if err != nil {
		return "", fmt.Errorf("некорректная дата: %w", err)
	}

	parts := strings.Split(repeat, " ")
	switch parts[0] {
	case "d":
		if err := applyDays(&date, now, parts); err != nil {
			return "", err
		}
	case "y":
		if len(parts) != 1 {
			return "", fmt.Errorf("неверный формат правила y")
		}
		for {
			date = date.AddDate(1, 0, 0)
			if afterNow(date, now) {
				break
			}
		}
	case "w":
		if err := applyWeekdays(&date, now, parts); err != nil {
			return "", err
		}
	case "m":
		if err := applyMonthDays(&date, now, parts); err != nil {
			return "", err
		}
	default:
		return "", fmt.Errorf("неподдерживаемый формат правила повторения: %q", repeat)
	}

	return date.Format(dateFormat), nil
}

func applyDays(date *time.Time, now time.Time, parts []string) error {
	if len(parts) != 2 {
		return fmt.Errorf("не указан интервал в днях")
	}
	days, err := strconv.Atoi(parts[1])
	if err != nil {
		return fmt.Errorf("некорректный интервал в днях: %w", err)
	}
	if days < 1 || days > 400 {
		return fmt.Errorf("интервал в днях должен быть от 1 до 400")
	}
	for {
		*date = date.AddDate(0, 0, days)
		if afterNow(*date, now) {
			break
		}
	}
	return nil
}

func applyWeekdays(date *time.Time, now time.Time, parts []string) error {
	if len(parts) != 2 {
		return fmt.Errorf("не указаны дни недели")
	}
	weekdays, err := parseIntListRange(parts[1], 1, 7)
	if err != nil {
		return fmt.Errorf("некорректные дни недели: %w", err)
	}
	allowed := make(map[int]bool, len(weekdays))
	for _, d := range weekdays {
		allowed[d] = true
	}
	for {
		*date = date.AddDate(0, 0, 1)
		weekday := int(date.Weekday())
		if weekday == 0 { // в time.Weekday воскресенье это 0, а у нас 7
			weekday = 7
		}
		if allowed[weekday] && afterNow(*date, now) {
			break
		}
	}
	return nil
}

func applyMonthDays(date *time.Time, now time.Time, parts []string) error {
	if len(parts) < 2 || len(parts) > 3 {
		return fmt.Errorf("неверный формат правила m")
	}

	days, err := parseMonthDays(parts[1])
	if err != nil {
		return err
	}

	var months []int
	if len(parts) == 3 {
		months, err = parseIntListRange(parts[2], 1, 12)
		if err != nil {
			return fmt.Errorf("некорректные месяцы: %w", err)
		}
	}
	allowedMonths := make(map[int]bool, len(months))
	for _, m := range months {
		allowedMonths[m] = true
	}

	matches := func(d time.Time) bool {
		if len(allowedMonths) > 0 && !allowedMonths[int(d.Month())] {
			return false
		}
		last := lastDayOfMonth(d)
		for _, day := range days {
			switch {
			case day > 0 && d.Day() == day:
				return true
			case day == -1 && d.Day() == last:
				return true
			case day == -2 && d.Day() == last-1:
				return true
			}
		}
		return false
	}

	for {
		*date = date.AddDate(0, 0, 1)
		if matches(*date) && afterNow(*date, now) {
			break
		}
	}
	return nil
}

func lastDayOfMonth(t time.Time) int {
	firstOfNextMonth := time.Date(t.Year(), t.Month()+1, 1, 0, 0, 0, 0, t.Location())
	return firstOfNextMonth.AddDate(0, 0, -1).Day()
}

func parseIntListRange(s string, min, max int) ([]int, error) {
	items := strings.Split(s, ",")
	result := make([]int, 0, len(items))
	for _, item := range items {
		n, err := strconv.Atoi(item)
		if err != nil {
			return nil, fmt.Errorf("некорректное число %q", item)
		}
		if n < min || n > max {
			return nil, fmt.Errorf("значение %d вне диапазона %d..%d", n, min, max)
		}
		result = append(result, n)
	}
	return result, nil
}

// parseMonthDays разбирает список дней месяца: допустимы значения 1-31,
// а также -1 (последний день месяца) и -2 (предпоследний)
func parseMonthDays(s string) ([]int, error) {
	items := strings.Split(s, ",")
	result := make([]int, 0, len(items))
	for _, item := range items {
		n, err := strconv.Atoi(item)
		if err != nil {
			return nil, fmt.Errorf("некорректный день месяца %q", item)
		}
		if n == 0 || n < -2 || n > 31 {
			return nil, fmt.Errorf("недопустимый день месяца %d", n)
		}
		result = append(result, n)
	}
	return result, nil
}

func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	nowParam := r.FormValue("now")
	dateParam := r.FormValue("date")
	repeat := r.FormValue("repeat")

	now := time.Now()
	if nowParam != "" {
		parsed, err := time.Parse(dateFormat, nowParam)
		if err != nil {
			http.Error(w, "некорректный параметр now", http.StatusBadRequest)
			return
		}
		now = parsed
	}

	next, err := NextDate(now, dateParam, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Write([]byte(next))
}

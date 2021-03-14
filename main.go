package main

import (
	"encoding/json"
	"fmt"
	"github.com/ajstarks/svgo"
	"io"
	"log"
	"net/http"
	"time"
)

// Структура квадратов
type SquareCalendar struct {
	size   int
	margin int
}

// Структура поле для квадратов
type Calendar struct {
	days           int
	margin         int
	squareCalendar SquareCalendar
	data           [][]int
}

// Запус слушателя
func main() {
	http.HandleFunc("/gitlab/", handlerRoute)

	err := http.ListenAndServe(":2004", nil)

	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

// Генерация SVG в виде конечного календаря
func handlerRoute(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/svg+xml")

	//убрать "/gitlab/" символы
	trimmedURL := r.URL.Path[8:]

	if trimmedURL != "" {
		var title string

		if string(trimmedURL[len(trimmedURL)-1]) == "/" {
			title = trimmedURL[0:(len(trimmedURL) - 1)]
		} else {
			title = trimmedURL
		}

		calendar := Calendar{days: 7, margin: 20}
		calendar.squareCalendar = SquareCalendar{size: 15, margin: 2}
		calendar.data = getCalendarCommitDataGitlab(title)

		s := svg.New(w)
		s.Start(
			len(calendar.data)*calendar.squareCalendar.size+calendar.margin,
			calendar.days*calendar.squareCalendar.size+calendar.margin,
		)

		renderCalendar(s, &calendar)

		s.End()
	} else {
		handleErrorCode(404, "Path not found.", w)
	}
}

// Генерация ошибочной страницы
func handleErrorCode(errorCode int, description string, w http.ResponseWriter) {
	w.WriteHeader(errorCode) // изменить HTTP status code (пример 404, 500)
	w.Header().Set("Content-Type", "text/html")
	_, _ = w.Write([]byte(fmt.Sprintf(
		"<html><body><h1>Error %d</h1><p>%s</p></body></html>",
		errorCode,
		description)))
}

// Получить готовые данные в двумерном массиве с количеством коммитов
func getCalendarCommitDataGitlab(username string) [][]int {
	response, err := http.Get("https://gitlab.com/users/" + username + "/calendar.json")

	if err != nil {
		log.Fatal("ErrRequest:", err)
	}

	defer response.Body.Close()

	var body map[string]interface{}
	bodyResponse, _ := io.ReadAll(response.Body)
	_ = json.Unmarshal(bodyResponse, &body)

	calendarDates := getCalendarDates()
	calendarCommitData := make([][]int, len(calendarDates))

	for weekIndex, weekValue := range calendarDates {
		days := make([]int, len(weekValue))

		for dayIndex, dayValue := range weekValue {
			if body[dayValue] != nil {
				days[dayIndex] = int(body[dayValue].(float64))
			} else {
				days[dayIndex] = 0
			}
		}

		calendarCommitData[weekIndex] = days
	}

	return calendarCommitData
}

// Получение от-до дат в двумерном массиве
func getCalendarDates() [][]string {
	timeNow := time.Now()
	timeAfter := timeNow.AddDate(0, -9, 0)
	timeAfter = timeAfter.AddDate(0, 0, -int(timeAfter.Weekday()+0))
	timeDiff := int(timeNow.Sub(timeAfter) / (24 * time.Hour))

	var calendarDates []string

	for i := timeDiff; i >= 0; i-- {
		calendarDates = append(calendarDates, timeNow.AddDate(0, 0, -i).Format("2006-01-02"))
	}

	return chunkBy(calendarDates, 7)
}

// Генерация календаря
func renderCalendar(s *svg.SVG, calendar *Calendar) {
	for weekIndex, weekValue := range calendar.data {
		for dayIndex, dayValue := range weekValue {
			var style string

			switch {
			case dayValue == 0:
				style = "fill:#ededed"
			case dayValue < 10:
				style = "fill:rgb(172, 213, 242)"
			case dayValue < 20:
				style = "fill:rgb(127, 168, 201)"
			case dayValue < 30:
				style = "fill:rgb(82, 123, 160)"
			default:
				style = "fill:rgb(37, 78, 119)"
			}

			s.Square(
				calendar.margin/2+calendar.squareCalendar.size*weekIndex,
				calendar.margin/2+calendar.squareCalendar.size*dayIndex,
				calendar.squareCalendar.size-calendar.squareCalendar.margin,
				style,
			)
		}
	}
}

// Разбиение []int на чанки по n размеру чанка
func chunkBy(items []string, chunkSize int) (chunks [][]string) {
	for chunkSize < len(items) {
		items, chunks = items[chunkSize:], append(chunks, items[0:chunkSize:chunkSize])
	}

	return append(chunks, items)
}

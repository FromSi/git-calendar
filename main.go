package main

import (
	"github.com/ajstarks/svgo"
	"log"
	"math/rand"
	"net/http"
)

type SquareCalendar struct {
	size   int
	margin int
}

type Calendar struct {
	days           int
	margin         int
	squareCalendar SquareCalendar
	data           [][]int
}

// Запус слушателя
func main() {
	http.Handle("/", http.HandlerFunc(render))

	err := http.ListenAndServe(":2004", nil)

	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

// Генерация SVG в виде конечного календаря
func render(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "image/svg+xml")

	calendar := Calendar{days: 7, margin: 20}
	calendar.squareCalendar = SquareCalendar{size: 15, margin: 2}
	calendar.data = chunkBy(make([]int, 336), calendar.days) //TODO: пока в разработке

	s := svg.New(w)
	s.Start(
		len(calendar.data)*calendar.squareCalendar.size+calendar.margin,
		calendar.days*calendar.squareCalendar.size+calendar.margin,
	)

	renderCalendar(s, &calendar)

	s.End()
}

func renderCalendar(s *svg.SVG, calendar *Calendar) {
	for i, week := range calendar.data {
		for j, day := range week {
			println(day)                     //TODO: пока в разработке
			randomNumber := rand.Intn(5) - 1 //TODO: пока в разработке

			style := ""

			switch randomNumber { //TODO: пока в разработке
			case 0:
				style = "fill:#ededed"
			case 1:
				style = "fill:rgb(172, 213, 242)"
			case 2:
				style = "fill:rgb(127, 168, 201)"
			case 3:
				style = "fill:rgb(82, 123, 160)"
			case 4:
				style = "fill:rgb(37, 78, 119)"
			}

			s.Square(
				calendar.margin/2+calendar.squareCalendar.size*i,
				calendar.margin/2+calendar.squareCalendar.size*j,
				calendar.squareCalendar.size-calendar.squareCalendar.margin,
				style,
			)
		}
	}
}

// Разбиение []int на чанки по n размеру чанка
func chunkBy(items []int, chunkSize int) (chunks [][]int) {
	for chunkSize < len(items) {
		items, chunks = items[chunkSize:], append(chunks, items[0:chunkSize:chunkSize])
	}

	return append(chunks, items)
}

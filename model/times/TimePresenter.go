package times

import (
	// "fmt"
	"math"
	"strconv"
)

type TimePresenter struct {}

func (t *TimePresenter)plural(count int, singular string) (result string) {
	if (count == 1) || (count == 0) {
			result = strconv.Itoa(count) + " " + singular + " "
	} else {
			result = strconv.Itoa(count) + " " + singular + "s "
	}
	return
}

func (t *TimePresenter)SecondsToHuman(input int) (result string) {
	years := math.Floor(float64(input) / 60 / 60 / 24 / 7 / 30 / 12)
	segundos := input % (60 * 60 * 24 * 7 * 30 * 12)
	months := math.Floor(float64(segundos) / 60 / 60 / 24 / 7 / 30)
	segundos = input % (60 * 60 * 24 * 7 * 30)
	weeks := math.Floor(float64(segundos) / 60 / 60 / 24 / 7)
	segundos = input % (60 * 60 * 24 * 7)
	days := math.Floor(float64(segundos) / 60 / 60 / 24)
	segundos = input % (60 * 60 * 24)
	horas := math.Floor(float64(segundos) / 60 / 60)
	segundos = input % (60 * 60)
	minutos := math.Floor(float64(segundos) / 60)
	segundos = input % 60

	if years > 0 {
			result = t.plural(int(years), "years") + t.plural(int(months), "month") + t.plural(int(weeks), "week") + t.plural(int(days), "day") + t.plural(int(horas), "hora") + t.plural(int(minutos), "minuto") + t.plural(int(segundos), "segundo")
	} else if months > 0 {
			result = t.plural(int(months), "month") + t.plural(int(weeks), "week") + t.plural(int(days), "day") + t.plural(int(horas), "hora") + t.plural(int(minutos), "minuto") + t.plural(int(segundos), "segundo")
	} else if weeks > 0 {
			result = t.plural(int(weeks), "week") + t.plural(int(days), "day") + t.plural(int(horas), "hora") + t.plural(int(minutos), "minuto") + t.plural(int(segundos), "segundo")
	} else if days > 0 {
			result = t.plural(int(days), "dias") + t.plural(int(horas), "hora") + t.plural(int(minutos), "minuto") + t.plural(int(segundos), "segundo")
	} else if horas > 0 {
			result = t.plural(int(horas), "horas") + t.plural(int(minutos), "minuto") + t.plural(int(segundos), "segundo")
	} else if minutos > 0 {
			result = t.plural(int(minutos), "minuto") + t.plural(int(segundos), "segundo")
	} else {
			result = t.plural(int(segundos), "segundo")
	}

	return
}

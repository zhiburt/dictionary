package wheather

import (
	"log"

	owm "github.com/briandowns/openweathermap"
)

// WheatherToday gets the wheather today
func WheatherToday(apiKey string) string {
	w, err := owm.NewCurrent("C", "ru", apiKey)
	if err != nil {
		print(apiKey)
		log.Fatalln(err)
	}
	w.CurrentByName("Minsk")

	return FormatWhether(*w).WheatherFormat()
}

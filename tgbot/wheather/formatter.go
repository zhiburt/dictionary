package wheather

import (
	"fmt"
	"time"

	owm "github.com/briandowns/openweathermap"
)

// FormatWhether wrap  owm.CurrentWeatherData
type FormatWhether owm.CurrentWeatherData

// WheatherFormat format api data to a beautiful string
func (w FormatWhether) WheatherFormat() string {
	emojie := getEmojie(w)
	return fmt.Sprintf(`The wheather today

%v %v
temperature %v℃
sunrise at %s
sunset at %s

have a great day!`, w.Name, emojie, w.Main.Temp,
		time.Unix(int64(w.Sys.Sunrise), 0).Format(time.RFC822),
		time.Unix(int64(w.Sys.Sunset), 0).Format(time.RFC822))
}

func getEmojie(w FormatWhether) string {
	if w.Rain.ThreeH > 25 {
		return "💦"
	}
	if w.Snow.ThreeH > 10 {
		return "❄️"
	}
	if w.Clouds.All > 20 {
		return "☁️"
	}

	return "🤔"
}

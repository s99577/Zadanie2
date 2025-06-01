package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

const (
	openWeatherAPI = "https://api.openweathermap.org/data/2.5/weather"
	units          = "metric"
	lang           = "pl"
)

type WeatherInfo struct {
	City    string `json:"name"`
	Weather []struct {
		Description string `json:"description"`
	} `json:"weather"`
	Main struct {
		Temp      float32 `json:"temp"`
		FeelsLike float32 `json:"feels_like"`
		Pressure  uint16  `json:"pressure"`
		Humidity  uint8   `json:"humidity"`
	} `json:"main"`
	Wind struct {
		Speed float32 `json:"speed"`
	} `json:"wind"`
}

type FormData struct {
	Countries map[string]string
	Cities    map[string][]string
}

func main() {
	author := os.Getenv("AUTHOR")
	port := os.Getenv("PORT")
	logInfo(author, port)

	http.HandleFunc("/", formHandler)
	http.HandleFunc("/weather", weatherHandler)

	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func logInfo(author string, port string) {
	currentTime := time.Now().Format("02-01-2006 15:04:05")
	log.Println("Data uruchomienia: ", currentTime)
	log.Println("Autor: ", author)
	log.Printf("Aplikacja nasłuchuje na porcie %s ...", port)
}

func formHandler(w http.ResponseWriter, r *http.Request) {
	data := FormData{
		Countries: map[string]string{
			"PL": "Polska",
			"US": "USA",
			"GB": "Wielka Brytania",
			"NO": "Norwegia",
			"NL": "Holandia",
		},
		Cities: map[string][]string{
			"PL": {"Lublin", "Warszawa", "Kraków", "Gdańsk"},
			"US": {"New York", "Los Angeles", "Chicago", "San Francisco"},
			"GB": {"Londyn", "Manchester", "Edynburg", "Bristol"},
			"NO": {"Oslo", "Bergen", "Stavanger", "Trondheim"},
			"NL": {"Amsterdam", "Rotterdam", "Haga", "Utrecht"},
		},
	}
	tmpl := template.Must(template.New("form.html").Funcs(template.FuncMap{
		"toJSON": func(v interface{}) template.JS {
			b, _ := json.Marshal(v)
			return template.JS(b)
		},
	}).ParseFiles("templates/form.html"))
	tmpl.Execute(w, data)
}

func weatherHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	country := r.FormValue("country")
	city := r.FormValue("city")
	city = url.QueryEscape(city)
	apiKey := os.Getenv("WEATHER_API_KEY")
	if apiKey == "" {
		http.Error(w, "API key not set", http.StatusInternalServerError)
		return
	}
	url := fmt.Sprintf("%s?q=%s,%s&appid=%s&units=%s&lang=%s", openWeatherAPI, city, country, apiKey, units, lang)

	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, "Failed to contact weather API", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		log.Println("API error: ", string(body))
		http.Error(w, "Weather API error", http.StatusBadGateway)
		return
	}

	var weather WeatherInfo
	err = json.NewDecoder(resp.Body).Decode(&weather)
	if err != nil {
		http.Error(w, "Failed to parse weather data", http.StatusInternalServerError)
		return
	}
	tmpl := template.Must(template.ParseFiles("templates/weather.html"))
	tmpl.Execute(w, weather)
}

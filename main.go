package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

// WeatherData represents the structure of the weather data
type WeatherData struct {
	Temperature float64 `json:"temperature"`
	Description string  `json:"description"`
	// Add other relevant fields
}

// FetchWeatherData retrieves the current weather data for Toronto
func FetchWeatherData(apiKey string) (*WeatherData, error) {
	url := fmt.Sprintf("http://api.openweathermap.org/data/2.5/weather?q=Toronto&appid=%s", apiKey)

	// Make a GET request to the OpenWeatherMap API
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Parse the JSON response
	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	// Extract relevant information from the parsed JSON
	weather := data["weather"].([]interface{})[0].(map[string]interface{})
	description := weather["description"].(string)

	main := data["main"].(map[string]interface{})
	temperature := main["temp"].(float64)

	return &WeatherData{
		Temperature: temperature,
		Description: description,
	}, nil
}

// WeatherHandler handles requests to the /weather endpoint
func WeatherHandler(w http.ResponseWriter, r *http.Request) {
	apiKey := "YOUR_API_KEY" // Replace with your actual OpenWeatherMap API key

	// Call FetchWeatherData to get weather data
	weatherData, err := FetchWeatherData(apiKey)
	if err != nil {
		// Handle error
		http.Error(w, "Error fetching weather data", http.StatusInternalServerError)
		return
	}

	// Convert weatherData to JSON format
	jsonData, err := json.Marshal(weatherData)
	if err != nil {
		// Handle error
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}

	// Set the content type header to application/json
	w.Header().Set("Content-Type", "application/json")

	// Write the JSON data to the response writer
	w.Write(jsonData)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/weather", WeatherHandler).Methods("GET")

	http.Handle("/", r)
	http.ListenAndServe(":8082", nil)
}

package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type GeoCodingResponse []struct {
	Name       string `json:"name"`
	LocalNames struct {
		El string `json:"el"`
		Sk string `json:"sk"`
		Si string `json:"si"`
		Kn string `json:"kn"`
		Ko string `json:"ko"`
		Th string `json:"th"`
		Ar string `json:"ar"`
		Te string `json:"te"`
		Or string `json:"or"`
		Ta string `json:"ta"`
		Ks string `json:"ks"`
		Fa string `json:"fa"`
		Ps string `json:"ps"`
		Az string `json:"az"`
		Ru string `json:"ru"`
		Ka string `json:"ka"`
		Lt string `json:"lt"`
		Ml string `json:"ml"`
		Bn string `json:"bn"`
		Oc string `json:"oc"`
		Ur string `json:"ur"`
		En string `json:"en"`
		Gu string `json:"gu"`
		Fr string `json:"fr"`
		Sd string `json:"sd"`
		He string `json:"he"`
		Eo string `json:"eo"`
		Pl string `json:"pl"`
		Hi string `json:"hi"`
		Zh string `json:"zh"`
		Pa string `json:"pa"`
		Sr string `json:"sr"`
		Es string `json:"es"`
		De string `json:"de"`
		Ia string `json:"ia"`
		Io string `json:"io"`
		Mr string `json:"mr"`
		Yi string `json:"yi"`
		Cs string `json:"cs"`
		Ja string `json:"ja"`
		Uk string `json:"uk"`
	} `json:"local_names"`
	Lat     float64 `json:"lat"`
	Lon     float64 `json:"lon"`
	Country string  `json:"country"`
	State   string  `json:"state"`
}

type WeatherResponse struct {
	Coord struct {
		Lon float64 `json:"lon"`
		Lat float64 `json:"lat"`
	} `json:"coord"`
	Weather []struct {
		ID          int    `json:"id"`
		Main        string `json:"main"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
	} `json:"weather"`
	Base string `json:"base"`
	Main struct {
		Temp      float64 `json:"temp"`
		FeelsLike float64 `json:"feels_like"`
		TempMin   float64 `json:"temp_min"`
		TempMax   float64 `json:"temp_max"`
		Pressure  int     `json:"pressure"`
		Humidity  int     `json:"humidity"`
		SeaLevel  int     `json:"sea_level"`
		GrndLevel int     `json:"grnd_level"`
	} `json:"main"`
	Visibility int `json:"visibility"`
	Wind       struct {
		Speed float64 `json:"speed"`
		Deg   int     `json:"deg"`
		Gust  float64 `json:"gust"`
	} `json:"wind"`
	Clouds struct {
		All int `json:"all"`
	} `json:"clouds"`
	Dt  int `json:"dt"`
	Sys struct {
		Country string `json:"country"`
		Sunrise int    `json:"sunrise"`
		Sunset  int    `json:"sunset"`
	} `json:"sys"`
	Timezone int    `json:"timezone"`
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Cod      int    `json:"cod"`
}

func getWeatherData(lat float64, lon float64, city string) {
	resp, err := http.Get(fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?lat=%f&lon=%f&units=metric&appid=%s", lat, lon, os.Getenv("OPEN_WEATHERMAP_API_KEY")))

	if err != nil {
		fmt.Println("Error fetching weather data:", err)
		return
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	var result WeatherResponse

	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Println("Error unmarshalling response:", err)
		return
	}

	fmt.Printf("Current weather for %s - %s, %.2f°C\n", strings.ToUpper(city[:1])+strings.ToLower(city[1:]), result.Weather[0].Main, result.Main.Temp)
	fmt.Printf("Feels like: %.2f°C\n", result.Main.FeelsLike)
	fmt.Printf("Humidity: %d%%\n", result.Main.Humidity)
	fmt.Printf("Wind Speed: %.2f m/s\n", result.Wind.Speed)
}

func getGeoCodingData(city string) (float64, float64, error) {
	resp, err := http.Get(fmt.Sprintf("http://api.openweathermap.org/geo/1.0/direct?q=%s&limit=5&appid=%s", city, os.Getenv("OPEN_WEATHERMAP_API_KEY")))

	if err != nil {
		fmt.Println("Error fetching geocoding data:", err)
		return 0, 0, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	var result GeoCodingResponse

	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Println("Error unmarshalling response:", err)
	}

	return result[0].Lat, result[0].Lon, nil
}

func main() {
	err := godotenv.Load()

	if err != nil {
		fmt.Println("Error loading .env file")
		return
	}

	fmt.Println("This is a simple weather app.")
	fmt.Println("Enter city name to get the weather forecast.")

	scanner := bufio.NewScanner(os.Stdin)

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}

	for scanner.Scan() {
		city := scanner.Text()
		lat, lon, err := getGeoCodingData(city)

		if err != nil {
			fmt.Println("Error getting geocoding data:", err)
			return
		}

		getWeatherData(lat, lon, city)

		os.Exit(0)
	}
}

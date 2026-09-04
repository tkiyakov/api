package main 

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"encoding/json"
	"bufio"
	"os"
	"strings"
	"strconv"
)

type Weather struct {
	Current struct {
		Temperature float64 `json:"temperature_2m"`
	} `json:"current"`
}

type GeoResponse struct {
	Results []struct {
		Admin1 string `json:"admin1"`
		Name string `json:"name"`
		Latitude float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
		Country string `json:"country"`
		Population int `json:"population"`
	} `json:"results"`
}


func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Введите город: ")
	city, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Ошибка ввода:", err)
		return
	}

	city = strings.TrimSpace(city)

	geoURL := "https://geocoding-api.open-meteo.com/v1/search?name=" + url.QueryEscape(city) + "&count=10&language=ru"

	
	resp, err := http.Get(geoURL)
	if err != nil {
		fmt.Println("Ошибка запроса", err)
		return
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Ошибка чтения:", err)
		return
	}

	var geo GeoResponse
	err = json.Unmarshal(body, &geo)
	if err != nil {
		fmt.Println("Ошибка json:", err)
		return
	}

	if len(geo.Results) == 0 {
		fmt.Println("Город не найден")
		return
	}

	for i, place := range geo.Results {
		fmt.Println(i, place.Name, place.Country)
	}

	fmt.Print("Выберите город: ")

	choiceText, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Ошибка ввода:", err)
		return
	}

	choiceText = strings.TrimSpace(choiceText)

	choice, err := strconv.Atoi(choiceText)
	if err != nil {
		fmt.Println("Нужно ввести номер")
		return
	}

	place := geo.Results[choice]


	// for _, result := range geo.Results {
	// 	if result.Population > place.Population {
	// 		place = result
	// 	}
	// }

	// fmt.Println("Пользователь ввёл:", city)
	// fmt.Println("API нашёл:", place.Name)
	// fmt.Println("Страна:", place.Country)

	weatherURL := fmt.Sprintf("https://api.open-meteo.com/v1/forecast?latitude=%f&longitude=%f&current=temperature_2m", place.Latitude, place.Longitude)

	wresp, err := http.Get(weatherURL)
	if err != nil {
		fmt.Println("Ошибка запроса", err)
		return
	}

	defer resp.Body.Close()

	wbody, err := io.ReadAll(wresp.Body)
	if err != nil {
		fmt.Println("Ошибка чтения", err)
		return
	}

	var weather Weather

	err = json.Unmarshal(wbody, &weather)
	if err != nil {
		fmt.Println("Ошибка JSON:", err)
		return
	}

	fmt.Printf("Температура в %s: %.1f °C\n", place.Name, weather.Current.Temperature)

}
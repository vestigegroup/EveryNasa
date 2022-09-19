package functions

import (
	"encoding/json"
	"fmt"
	"github.com/jasonlvhit/gocron"
	"github.com/reujab/wallpaper"
	"github.com/rodkranz/fetch"
	"gopkg.in/yaml.v2"
	"os"
	"time"
)

func GetWallpaper() string {
	file, readErr := os.ReadFile("config.yaml")
	if readErr != nil {
		panic(readErr)
	}
	data := make(map[interface{}]interface{})

	marshalErr := yaml.Unmarshal(file, &data)
	if marshalErr != nil {
		panic(marshalErr)
	}

	type WallpaperStruct struct {
		Hdurl      string `json:"hdurl"`
		Url        string `json:"url"`
		Media_type string `json:"media_type"`
	}

	client := fetch.NewDefault()
	res, getErr := client.Get("https://api.nasa.gov/planetary/apod?api_key="+data["apiKey"].(string), nil)
	if getErr != nil {
		panic(getErr)
	}

	body, StringErr := res.ToString()
	if StringErr != nil {
		panic(StringErr)
	}

	var Wallpaper WallpaperStruct
	jsonErr := json.Unmarshal([]byte(body), &Wallpaper)

	if jsonErr != nil {
		panic(jsonErr)
	}

	if Wallpaper.Media_type == "video" {
		Wallpaper.Url = Wallpaper.Url[30 : len(Wallpaper.Url)-6]
		return "https://img.youtube.com/vi/" + Wallpaper.Url + "/0.jpg"
	}

	return Wallpaper.Hdurl
}

func SetWallpaper() {
	err := wallpaper.SetFromURL(GetWallpaper())
	if err != nil {
		panic(err)
	}
}

func x() {
	fmt.Println("F")
}

func StartWallpaper() {
	type Autostart struct {
		Autochangewallpaper int `json:"autochangewallpaper"`
	}

	client := fetch.NewDefault()
	res, getErr := client.Get("http://localhost:8080/api/get/settings", nil)
	if getErr != nil {
		panic(getErr)
	}

	body, StringErr := res.ToString()
	if StringErr != nil {
		panic(StringErr)
	}

	var AutostartSetWallpaper Autostart
	jsonErr := json.Unmarshal([]byte(body), &AutostartSetWallpaper)

	if jsonErr != nil {
		panic(jsonErr)
	}

	if AutostartSetWallpaper.Autochangewallpaper == 1 {
		times := time.Now()
		t := time.Date(times.Year(), times.Month(), times.Day(), 4, 50, times.Second(), times.Nanosecond(), time.UTC)

		SetWallpaper()

		err := gocron.Every(1).Day().From(&t).Do(SetWallpaper)
		if err != nil {
			panic(err)
		}
		gocron.Start()
	}
}

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type App struct {
	ID   int    `json:"appid"`
	Name string `json:"name"`
}

type AppList struct {
	Apps      []App     `json:"apps"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (appList *AppList) Find(id int) string {
	for _, app := range appList.Apps {
		if app.ID == id {
			return app.Name
		}
	}

	return "UNKNOWN"
}

type applistResponse struct {
	AppList AppList `json:"applist"`
}

func LoadAppList() (*AppList, error) {
	shouldRefresh := false

	// try to open exiting file
	file, err := os.ReadFile("applist.json")
	if err != nil {
		fmt.Println("Failed to read applist.json, will refetch.")
		shouldRefresh = true
		file, _ = json.Marshal(AppList{})
	}

	// parse json
	applist := &AppList{}
	if err := json.Unmarshal(file, applist); err != nil {
		return nil, err
	}

	// detect if file is older than two weeks
	if applist.UpdatedAt.Before(time.Now().AddDate(0, 0, -14)) {
		fmt.Println("applist.json appears outdated, will refetch.")
		shouldRefresh = true
	}

	// fetch from steam api if needed
	if shouldRefresh {
		fmt.Println("Fetching applist from Steam API.")

		// get data from steam api
		res, err := http.Get("https://api.steampowered.com/ISteamApps/GetAppList/v2/")
		if err != nil {
			return nil, err
		}
		if res.StatusCode != 200 {
			return nil, fmt.Errorf("Steam API returned unexpected code: %d", res.StatusCode)
		}

		// parse json
		appListRes := &applistResponse{}
		decoder := json.NewDecoder(res.Body)
		if err := decoder.Decode(appListRes); err != nil {
			return nil, err
		}

		// write new data to file
		appListRes.AppList.UpdatedAt = time.Now()
		appListJson, err := json.Marshal(appListRes.AppList)
		if err != nil {
			return nil, err
		}
		if err := os.WriteFile("applist.json", appListJson, 0666); err != nil {
			return nil, err
		}
		applist = &appListRes.AppList
	}

	return applist, nil
}

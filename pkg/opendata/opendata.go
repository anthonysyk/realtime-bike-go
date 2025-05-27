package opendata

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const URL = "https://opendata.paris.fr/api/explore/v2.1/catalog/datasets/velib-disponibilite-en-temps-reel/records"

type Station struct {
	Stationcode       string    `json:"stationcode" bson:"stationcode"`
	Name              string    `json:"name" bson:"name"`
	IsInstalled       string    `json:"is_installed" bson:"is_installed"`
	Capacity          int       `json:"capacity" bson:"capacity"`
	Numdocksavailable int       `json:"numdocksavailable" bson:"numdocksavailable"`
	Numbikesavailable int       `json:"numbikesavailable" bson:"numbikesavailable"`
	Mechanical        int       `json:"mechanical" bson:"mechanical"`
	Ebike             int       `json:"ebike" bson:"ebike"`
	IsRenting         string    `json:"is_renting" bson:"is_renting"`
	IsReturning       string    `json:"is_returning" bson:"is_returning"`
	Duedate           time.Time `json:"duedate" bson:"duedate"`
	CoordonneesGeo    struct {
		Lon float64 `json:"lon" bson:"lon"`
		Lat float64 `json:"lat" bson:"lat"`
	} `json:"coordonnees_geo" bson:"coordonnees_geo"`
	NomArrondissementCommunes string      `json:"nom_arrondissement_communes" bson:"nom_arrondissement_communes"`
	CodeInseeCommune          string      `json:"code_insee_commune" bson:"code_insee_commune"`
	StationOpeningHours       interface{} `json:"station_opening_hours" bson:"station_opening_hours"`
}

func (s Station) IsActive() bool {
	return s.IsInstalled == "OUI" && s.IsRenting == "OUI" && s.IsReturning == "OUI"
}

type Response struct {
	TotalCount int       `json:"total_count"`
	Results    []Station `json:"results"`
	HasNext    bool      `json:"hasNext"`
}

func GetStationAvailability(limit int, offset int) (*Response, error) {
	baseURL, err := url.Parse(URL)
	if err != nil {
		return nil, fmt.Errorf("could not parse URL: %v", err)
	}

	params := url.Values{}
	params.Add("order_by", "stationcode")

	if limit > 0 {
		params.Add("limit", fmt.Sprintf("%d", limit))
	}
	if offset >= 0 {
		params.Add("offset", fmt.Sprintf("%d", offset))
	}
	baseURL.RawQuery = params.Encode()

	resp, err := http.Get(baseURL.String())
	if err != nil {
		return nil, fmt.Errorf("error during http query: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error status not ok: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read body: %v", err)
	}

	var res Response
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, fmt.Errorf("could not parse JSON: %v", err)
	}

	res.HasNext = offset+len(res.Results) < res.TotalCount

	return &res, nil
}

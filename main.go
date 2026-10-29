package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
)

//MRData, StandingsTable, StandingsList, and DriverStanding


type Response struct {
	MRData MRData
}

type MRData struct {
	StandingsTable StandingsTable
}

type StandingsTable struct {
	Season string
	Round string
	StandingsLists []StandingsLists
}

type StandingsLists struct {
	Season string
	Round string
	DriverStandings []DriverStandings
}

type DriverStandings struct {

	Position string
	Points string
	Wins string
	Driver Driver
	Constructors []Constructor

}

type Driver struct {
	DriverId string
	PermanentNumber string
	Code string
	URL string
	GivenName string
	FamilyName string
	DateOfBirth string
	Nationality string
}

type Constructor struct {
	ConstructorID string
	URL string
	Name string
	Nationality string
}


func main () {
	
	season := flag.String("season", "2026", "F1 season year")
	flag.Parse()
	
	data, err := fetchStandings(*season)

	if err != nil {
		fmt.Println(err)
		return
	}
	
	lists := data.MRData.StandingsTable.StandingsLists
	if len(lists) == 0{
		fmt.Println("no standings data")
		return
	}

	drivers := lists[0].DriverStandings

	for _, s := range drivers {

		if len(s.Constructors) == 0 {
			fmt.Printf("%s - no constructor data\n", s.Driver.FamilyName)
			continue
		}
		fmt.Printf("P%s %s (%s) - %s pts\n", s.Position, s.Driver.FamilyName, s.Constructors[0].Name, s.Points )
	}



}

func fetchStandings(season string) (Response, error){

	url := fmt.Sprintf("https://api.jolpi.ca/ergast/f1/%s/driverstandings/",season)

	resp, err := http.Get(url) // captures a response code and an error code
	if err != nil {	// checks for error other an nil	
			return Response{}, err	
	}

	defer resp.Body.Close() // streams the data from Body and closes

	body, err := io.ReadAll(resp.Body) // capture body and error code

	var result Response
	if err != nil {	// check for error in io.ReadAll
			return Response{}, err
	}

	if resp.StatusCode != 200 {	// checks for other response
		return Response{}, fmt.Errorf("unexpected status: %s", resp.Status)
	}


	err = json.Unmarshal(body, &result)

	if err != nil {
		return Response{}, err

	}


	return result, nil
	
}

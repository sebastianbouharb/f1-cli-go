package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
)

//MRData, StandingsTable, StandingsList, and DriverStanding


type Response struct {	// contains all data from api
	MRData MRData
}

type MRData struct { // focuses on the API's standings table
	StandingsTable StandingsTable
}

type StandingsTable struct {	// brances of the Standings table according to API
	Season string
	Round string
	StandingsLists []StandingsLists
}

type StandingsLists struct {	// brances of the Standings list according to API
	Season string
	Round string
	DriverStandings []DriverStandings
}

type DriverStandings struct { // brances of the Drivers standings according to API

	Position string
	Points string
	Wins string
	Driver Driver
	Constructors []Constructor

}

type Driver struct {	// elements of a Driver
	DriverId string
	PermanentNumber string
	Code string
	URL string
	GivenName string
	FamilyName string
	DateOfBirth string
	Nationality string
}

type Constructor struct {	// elements of a Constructor
	ConstructorID string
	URL string
	Name string
	Nationality string
}


func main () {
	
	// create flag for specifying which season to pull data from 
	// "season" is the flag
	// 2026 is the default val
	// F1 season year is the --help text
	season := flag.String("season", "2026", "F1 season year")
	flag.Parse()
	
	data, err := fetchStandings(*season) // return data and err code, takes season pointer as param

	if err != nil { // other than nil check
		fmt.Println(err)
		return
	}
	
	lists := data.MRData.StandingsTable.StandingsLists // store the Standings Lists
	if len(lists) == 0{ // check data exists
		fmt.Println("no standings data")
		return
	}

	drivers := lists[0].DriverStandings //store the drivers standings (lists.[0] is the first index of StandingsLists, which is entire list of Drivers and their statistics)

	for _, s := range drivers { // for every driver in the standings

		if len(s.Constructors) == 0 { // checl constructor exists
			fmt.Printf("%s - no constructor data\n", s.Driver.FamilyName)
			continue // skip if no constructor data
		}
		fmt.Printf("P%s %s (%s) - %s pts\n", s.Position, s.Driver.FamilyName, s.Constructors[0].Name, s.Points ) // print formated driver data
	}



}

func fetchStandings(season string) (Response, error){

	url := fmt.Sprintf("https://api.jolpi.ca/ergast/f1/%s/driverstandings/",season) // store API url 

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
		return Response{}, fmt.Errorf("unexpected status: %s", resp.Status) // return the unexpected status
	}


	err = json.Unmarshal(body, &result)	// extract json body 

	if err != nil {	// checl error isnt nil
		return Response{}, err
	}


	return result, nil
	
}

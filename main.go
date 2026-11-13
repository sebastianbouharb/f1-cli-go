package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

//MRData, StandingsTable, StandingsList, and DriverStanding

const pageSize int = 100 // represents the number of drivers that
// are return per request

type Response struct { // contains all data from api
	MRData MRData
}

type MRData struct { // focuses on the API's standings table
	Total          string
	StandingsTable StandingsTable
}

type StandingsTable struct { // branches of the Standings table according to API
	Season         string
	Round          string
	StandingsLists []StandingsLists
}

type StandingsLists struct { // branches of the Standings list according to API
	Season               string
	Round                string
	DriverStandings      []DriverStandings
	ConstructorStandings []ConstructorStandings
}

// type ConstructorStandingsList struct {
// 	Season string
// 	Round string
// 	ConstructorStandings []ConstructorStandings
// }

type ConstructorStandings struct { // branches of the Constructors standings
	Position     string
	PositionText string
	Points       string
	Wins         string
	Constructor  Constructor
}

type DriverStandings struct { // branches of the Drivers standings according to API

	Position     string
	PositionText string
	Points       string
	Wins         string
	Driver       Driver
	Constructors []Constructor
}

type Driver struct { // elements of a Driver
	DriverId        string
	PermanentNumber string
	Code            string
	URL             string
	GivenName       string
	FamilyName      string
	DateOfBirth     string
	Nationality     string
}

type Constructor struct { // elements of a Constructor
	ConstructorID string
	URL           string
	Name          string
	Nationality   string
}

func main() {

	// create flag for specifying which season to pull data from
	// "season" is the flag
	// 2026 is the default val
	// F1 season year is the --help text
	season := flag.String("season", "2026", "F1 season year")

	// create a flag that will pull a driver's bio
	driverName := flag.String("driver", "", "driver Last name, e.g. hamilton")

	constructorStandings := flag.Bool("constructors", false, "constructor name e.g. ferrari, mercedes, red_bull...")
	flag.Parse()

	drivers, err := fetchAllDrivers(*season) //get full list of drivers w/ pagination
	if err != nil {
		fmt.Println(err)
		return
	}

	if len(drivers) == 0 { // check for data in the list
		fmt.Println("no driver standings data")
		return
	}

	if *driverName != "" { // if the --driver flag isn't empty
		for _, s := range drivers {
			if strings.EqualFold(s.Driver.FamilyName, *driverName) { // for every driver, check if the flag name is a name in the standings

				fmt.Printf("%s %s\n", s.Driver.GivenName, s.Driver.FamilyName) // print name
				if len(s.Constructors) > 0 {
					fmt.Printf("Team - %s\n", s.Constructors[0].Name) // if they have a constructor, print constructor
				}
				fmt.Printf("#%s\n", s.Driver.PermanentNumber) // print driver number

				fmt.Printf("P%s\n", s.PositionText) // print position
				fmt.Printf("%s pts\n", s.Points)    // print points
				fmt.Printf("Wins -  %s\n", s.Wins)  // prints wins

				fmt.Printf("Nationality - %s\n", s.Driver.Nationality)   // print nationality
				fmt.Printf("Date of Birth - %s\n", s.Driver.DateOfBirth) // print DOB

				return

			}
		}
		fmt.Println("driver not found") // if driver not found
		return
	}

	if *constructorStandings {

		constructors, cerr := fetchAllConstructors(*season)
		if cerr != nil {
			fmt.Println(cerr)
			return
		}

		if len(constructors) == 0 {
			fmt.Println("no constructor standings data")
			return
		}
		for _, c := range constructors { // for every driver in the standings

			if c.Constructor.Name == "" {
				fmt.Println("constructor data not found")
				continue
			}

			// print formated driver data
			if c.PositionText == "-" {
				fmt.Printf("- %s - %s pts\n", c.Constructor.Name, c.Points)

			} else {
				fmt.Printf("P%s %s - %s pts\n", c.PositionText, c.Constructor.Name, c.Points)
			}

		}
		return

	}

	for _, s := range drivers { // for every driver in the standings

		if len(s.Constructors) == 0 { // checl constructor exists
			fmt.Printf("%s - no constructor data\n", s.Driver.FamilyName)
			continue // skip if no constructor data
		}

		// print formated driver data
		if s.PositionText == "-" {
			fmt.Printf("- %s (%s) - %s pts\n", s.Driver.FamilyName, s.Constructors[0].Name, s.Points)
		} else {
			fmt.Printf("P%s %s (%s) - %s pts\n", s.PositionText, s.Driver.FamilyName, s.Constructors[0].Name, s.Points)
		}
	}

}

func fetchConstructorStandings(season string, offset int) (Response, error) {

	url := fmt.Sprintf("https://api.jolpi.ca/ergast/f1/%s/constructorstandings/?limit=%d&offset=%d", season, 100, offset) // store API url

	resp, err := http.Get(url) // captures a response code and an error code
	if err != nil {            // checks for error other an nil
		return Response{}, err
	}

	defer resp.Body.Close() // streams the data from Body and closes

	body, err := io.ReadAll(resp.Body) // capture body and error code

	var result Response
	if err != nil { // check for error in io.ReadAll
		return Response{}, err
	}

	if resp.StatusCode != 200 { // checks for other response
		return Response{}, fmt.Errorf("unexpected status: %s", resp.Status) // return the unexpected status
	}

	err = json.Unmarshal(body, &result) // extract json body

	if err != nil { // checl error isnt nil
		return Response{}, err
	}

	return result, nil

}

func fetchAllConstructors(season string) ([]ConstructorStandings, error) {

	var all []ConstructorStandings // create a list of all drivers
	offset := 0                    // set offset to zero

	for {
		resp, err := fetchConstructorStandings(season, offset)
		if err != nil { // check fetch doesnt return an error
			return nil, err
		}

		lists := resp.MRData.StandingsTable.StandingsLists // store the Standings lists
		if len(lists) == 0 {                               // check len isn't zero
			break
		}

		page := lists[0].ConstructorStandings // store the drivers standings page
		if len(page) == 0 {                   // if the page doesnt have anything,  break
			break
		}
		all = append(all, page...) // append the page to all

		total, err := strconv.Atoi(resp.MRData.Total) //convert the total to an int
		if err != nil {                               // check for error
			return nil, err
		}

		if len(all) >= total { // checl that the number of drivers in all is >= total
			break
		}

		offset += pageSize // incrememnt the offset by the page size
	}

	return all, nil // return the full list of drivers
}

func fetchStandings(season string, offset int) (Response, error) {

	url := fmt.Sprintf("https://api.jolpi.ca/ergast/f1/%s/driverstandings/?limit=%d&offset=%d", season, 100, offset) // store API url

	resp, err := http.Get(url) // captures a response code and an error code
	if err != nil {            // checks for error other an nil
		return Response{}, err
	}

	defer resp.Body.Close() // streams the data from Body and closes

	body, err := io.ReadAll(resp.Body) // capture body and error code

	var result Response
	if err != nil { // check for error in io.ReadAll
		return Response{}, err
	}

	if resp.StatusCode != 200 { // checks for other response
		return Response{}, fmt.Errorf("unexpected status: %s", resp.Status) // return the unexpected status
	}

	err = json.Unmarshal(body, &result) // extract json body

	if err != nil { // checl error isnt nil
		return Response{}, err
	}

	return result, nil

}

func fetchAllDrivers(season string) ([]DriverStandings, error) {

	var all []DriverStandings // create a list of all drivers
	offset := 0               // set offset to zero

	for {
		resp, err := fetchStandings(season, offset)
		if err != nil { // check fetch doesnt return an error
			return nil, err
		}

		lists := resp.MRData.StandingsTable.StandingsLists // store the Drivers Standings lists
		if len(lists) == 0 {                               // check len isn't zero
			break
		}

		page := lists[0].DriverStandings // store the drivers standings page
		if len(page) == 0 {              // if the page doesnt have anything,  break
			break
		}
		all = append(all, page...) // append the page to all

		total, err := strconv.Atoi(resp.MRData.Total) //convert the total to an int
		if err != nil {                               // check for error
			return nil, err
		}

		if len(all) >= total { // checl that the number of drivers in all is >= total
			break
		}

		offset += pageSize // incrememnt the offset by the page size
	}

	return all, nil // return the full list of drivers
}

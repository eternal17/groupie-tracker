package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

const API string = "https://groupietrackers.herokuapp.com/api/"

// declare structs with the same structure as the stuctured json, which we can later unmarshall after getting(http.GET) the data required.
type Artists struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
}

type Location struct {
	Index []struct {
		ID        int      `json:"id"`
		Locations []string `json:"locations"`
		// Dates     string   `json:"dates"`
	} `json:"index"`
}

type Date struct {
	Index []struct {
		ID    int      `json:"id"`
		Dates []string `json:"dates"`
	} `json:"index"`
}

type Relation struct {
	Index []struct {
		ID                int                 `json:"id"`
		DatesandLocations map[string][]string `json:"datesLocations"`
	}
}

type Combined struct {
	Artists
	Dates     []string
	Location  []string
	Relations map[string][]string
}

var tpl *template.Template

func main() {

	tpl = template.Must(template.ParseGlob("templates/*.html"))

	http.HandleFunc("/", indexHandler)
	// seeing css http.HandleFunc("/artist", artistPage)
	http.HandleFunc("/artist/", artistPage)
	http.Handle("/templates/", http.StripPrefix("/templates/", http.FileServer(http.Dir("./templates"))))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port if not specified (with respect to heroku hosting)
	}
	http.ListenAndServe(":"+port, nil)
}

//||||||||||||||||||||||||||MAIN ABOVE|||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||

// func for executing homepage
func indexHandler(w http.ResponseWriter, r *http.Request) {

	datStruct := datesStruct(API, "dates")
	locStruct := locationsStruct(API, "locations")
	relStruct := relationStruct(API, "relation")
	artStruct := artistsStruct(API, "artists")
	relMap := relationsMap(relStruct)

	// page not found error
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		fmt.Fprintf(w, "Status 404: Page Not Found")
		return
	}

	if err := tpl.ExecuteTemplate(w, "homepage.html", combinedArray(artStruct, datStruct, locStruct, relMap)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

//gets ID from the URL & executes templates with data specific to that ID
func artistPage(w http.ResponseWriter, r *http.Request) {

	datStruct := datesStruct(API, "dates")
	locStruct := locationsStruct(API, "locations")
	relStruct := relationStruct(API, "relation")
	artStruct := artistsStruct(API, "artists")
	relMap := relationsMap(relStruct)

	selection := r.URL.Query().Get("selection")
	selectionId, _ := strconv.Atoi(selection)

	// If the query is not a number between 1 and 52, return a 404 page.
	if selectionId < 1 || selectionId > 52 {
		http.NotFound(w, r)
		fmt.Fprintf(w, "Status 404: Page Not Found")
		return
	}
	if err := tpl.ExecuteTemplate(w, "artistpage.html", combinedArray(artStruct, datStruct, locStruct, relMap)[selectionId-1]); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

//unmarshals the given json and returns a struct with date data fully parsed
func datesStruct(API string, field string) Date {

	url := API + field
	date, err := http.Get(url)

	if err != nil {
		fmt.Println("No response from request")
	}
	defer date.Body.Close()

	// the ReadAll method reads the data as bytes. we can then convert to string to read all the data received from the response.
	body2, err := ioutil.ReadAll(date.Body)
	if err != nil {
		fmt.Println(err)
	}

	// create a slice for our JSON data to be unmarshalled into.
	var Dates Date

	// We are unmarshalling our data, recieving the data, and pushing it into our artists slice.
	err = json.Unmarshal(body2, &Dates)
	if err != nil {
		fmt.Println(err)
	}

	return Dates
}

//unmarshals the given json and returns a struct with location data fully parsed
func locationsStruct(API string, field string) Location {

	url := API + field
	date, err := http.Get(url)

	if err != nil {
		fmt.Println("No response from request")
	}
	defer date.Body.Close()

	// the ReadAll method reads the data as bytes. we can then convert to string to read all the data received from the response.
	body2, err := ioutil.ReadAll(date.Body)
	if err != nil {
		fmt.Println(err)
	}

	// create a slice for our JSON data to be unmarshalled into.
	var Place Location

	// We are unmarshalling our data, recieving the data, and pushing it into our artists slice.
	err = json.Unmarshal(body2, &Place)
	if err != nil {
		fmt.Println(err)
	}

	return Place
}

//unmarshals the given json and returns a struct with relation data fully parsed
func relationStruct(API string, field string) Relation {

	url := API + field
	date, err := http.Get(url)

	if err != nil {
		fmt.Println("No response from request")
	}
	defer date.Body.Close()

	// the ReadAll method reads the data as bytes. we can then convert to string to read all the data received from the response.
	body2, err := ioutil.ReadAll(date.Body)
	if err != nil {
		fmt.Println(err)
	}

	// create a slice for our JSON data to be unmarshalled into.
	Relations := Relation{}

	// We are unmarshalling our data, recieving the data, and pushing it into our artists slice.
	err = json.Unmarshal(body2, &Relations)
	if err != nil {
		fmt.Println(err)
	}

	return Relations
}

//unmarshals the given json and returns a struct with artist data fully parsed
func artistsStruct(API string, field string) []Artists {

	url := API + field

	date, err := http.Get(url)
	if err != nil {
		fmt.Println("No response from request")
	}
	defer date.Body.Close()

	// the ReadAll method reads the data as bytes. we can then convert to string to read all the data received from the response.
	body2, err := ioutil.ReadAll(date.Body)
	if err != nil {
		fmt.Println(err)
	}

	// create a slice for our JSON data to be unmarshalled into.
	Artist := []Artists{}

	// We are unmarshalling our data, recieving the data, and pushing it into our artists slice.
	err = json.Unmarshal(body2, &Artist)
	if err != nil {
		fmt.Println(err)
	}

	return Artist
}

//parses the struct returned by relationsStruct() further returning solely the relations within a slice of map
func relationsMap(a Relation) []map[string][]string {

	var datesandLocal []map[string][]string

	for _, value := range a.Index {

		datesandLocal = append(datesandLocal, value.DatesandLocations)
	}

	return datesandLocal

}

//returns an array of structs with each storing data for the artists given by the API
func combinedArray(a []Artists, b Date, c Location, d []map[string][]string) [52]Combined {
	var combinedsliced [52]Combined

	for i := 0; i < len(combinedsliced); i++ {

		combinedsliced[i] = Combined{a[i], b.Index[i].Dates, c.Index[i].Locations, d[i]}
	}
	return combinedsliced
}

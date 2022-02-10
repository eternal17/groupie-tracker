package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"strconv"
)

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
	http.ListenAndServe(":8080", nil)
}

/////////////////////MAIN ABOVE\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\

// func for executing homepage
func indexHandler(w http.ResponseWriter, r *http.Request) {

	// page not found error
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		fmt.Fprintf(w, "Status 404: Page Not Found")
		return
	}

	if err := tpl.ExecuteTemplate(w, "homepage.html", GetData()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

//gets ID from the URL & executes templates with data specific to that ID
func artistPage(w http.ResponseWriter, r *http.Request) {

	selection := r.URL.Query().Get("selection")
	selectionId, _ := strconv.Atoi(selection)

	if err := tpl.ExecuteTemplate(w, "artistpage.html", GetData()[selectionId-1]); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

//unmarshals the given json apis and returns an array of structs with all data fully parsed
func GetData() [52]Combined {
	date, err := http.Get("https://groupietrackers.herokuapp.com/api/dates")

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

	////////////////////////////LOCATION STRUCT///////////////////////////////////////////////////////////////////////////////

	// JSON response from the sample API artists page, using the Get method.
	local, err := http.Get("https://groupietrackers.herokuapp.com/api/locations")

	if err != nil {
		fmt.Println("No response from request")
	}

	defer local.Body.Close()
	// the ReadAll method reads the data as bytes. we can then convert to string to read all the data received from the response.
	body3, err := ioutil.ReadAll(local.Body)
	if err != nil {
		fmt.Println(err)
	}

	// create a slice for our JSON data to be unmarshalled into.
	var Place Location

	// We are unmarshalling our data, recieving the data, and pushing it into our artists slice.
	err = json.Unmarshal(body3, &Place)

	if err != nil {
		fmt.Println(err)
	}

	///////////////////////Relations\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\

	// JSON response from the sample API artists page, using the Get method.
	relation, err := http.Get("https://groupietrackers.herokuapp.com/api/relation")

	if err != nil {
		fmt.Println("No response from request")
	}

	defer local.Body.Close()
	// the ReadAll method reads the data as bytes. we can then convert to string to read all the data received from the response.
	body4, err := ioutil.ReadAll(relation.Body)
	if err != nil {
		fmt.Println(err)
	}

	// create a slice for our JSON data to be unmarshalled into.
	Relations := Relation{}

	// We are unmarshalling our data, recieving the data, and pushing it into our artists slice.
	err = json.Unmarshal(body4, &Relations)

	if err != nil {
		fmt.Println(err)
	}

	///////////////////////ARTISTS\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\

	// JSON response from the sample API artists page, using the Get method.
	artist, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")

	if err != nil {
		fmt.Println("No response from request")
	}

	defer local.Body.Close()
	// the ReadAll method reads the data as bytes. we can then convert to string to read all the data received from the response.
	body1, err := ioutil.ReadAll(artist.Body)
	if err != nil {
		fmt.Println(err)
	}

	// create a slice for our JSON data to be unmarshalled into.
	Artist := []Artists{}

	// We are unmarshalling our data, recieving the data, and pushing it into our artists slice.
	err = json.Unmarshal(body1, &Artist)

	if err != nil {
		fmt.Println(err)
	}

	///////////////////////\/\/\/\\/\/\/\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\

	var datesandLocal []map[string][]string

	for _, value := range Relations.Index {

		datesandLocal = append(datesandLocal, value.DatesandLocations)
	}

	///////////////////////////////////////////////////////////////////////

	var combinedsliced [52]Combined

	for i := 0; i < len(combinedsliced); i++ {

		combinedsliced[i] = Combined{Artist[i], Dates.Index[i].Dates, Place.Index[i].Locations, datesandLocal[i]}

	}

	return combinedsliced
}

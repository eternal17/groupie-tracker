package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// declare a struct with the same structure as the stuctured json, which we can later unmarshall after getting(http.GET) the data required.
type Artists struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Locations    string   `json:"locations"`
	ConcertDates string   `json:"concertDates"`
	Relations    string   `json:"relations"`
}

func main() {

	// JSON response from the sample API artists page, using the Get method.
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")

	if err != nil {
		fmt.Println("No response from request")
	}

	defer resp.Body.Close()
	// the ReadAll method reads the data as bytes. we can then convert to string to read all the data received from the response.
	body, err := ioutil.ReadAll(resp.Body)
	// fmt.Println(string(body))

	// create a slice for our JSON data to be unmarshalled into.
	var artists []Artists

	// We are unmarshalling our data, recieving the data, and pushing it into our artists slice.
	err = json.Unmarshal(body, &artists)

	if err != nil {
		fmt.Println(err)
	}

	// Now we can extrapolate the data we want by ranging through our slice, eg ID and Name.
	for _, artist := range artists {
		fmt.Println(artist.ID, artist.Name)
	}

}

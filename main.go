package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type cliCommand struct {
	name        string
	description string
	callback    func(config *RequestConfig) error
}

type pokedexLocations struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type RequestConfig struct {
	// *string cause it can be nill too
	Next *string
	Prev *string
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exits the pokedex",
			callback:    commandExit,
		},
		"map": {
			name:        "map",
			description: "Lets you explore the map in skips of 20",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "To go back 20 skips in map locations",
			callback:    commandMapb,
		},
	}
}

func main() {

	config := &RequestConfig{}

	for {
		fmt.Printf("Pokedex > ")
		scanner := bufio.NewScanner(os.Stdin)
		if ok := scanner.Scan(); ok {
			for k, v := range getCommands() {
				if k == scanner.Text() {
					v.callback(config)
				}
			}
		}
	}
}

func commandExit(config *RequestConfig) error {
	os.Exit(0)
	return errors.New("something went wrong in exit")
}

func commandHelp(config *RequestConfig) error {
	commands := getCommands()
	fmt.Printf(`Welcome to Pokedex
  Usage:
`)

	for _, v := range commands {
		fmt.Printf("\t%v: %v\n", v.name, v.description)
	}

	return errors.New("something went wrong in help")
}

// commandMap will display 20 location areas in the world,
// subsequent calls should fetch the next 20 locations
func commandMap(config *RequestConfig) error {
	// https://pokeapi.co/api/v2/location/{id or name}/ API to hit.

	// Check if next exists and then make a call to that.
	// Start at 0 always
	url := ""
	if config.Next == nil {
		fmt.Println("Starting from first")
		url = "https://pokeapi.co/api/v2/location/"
	} else {
		url = *config.Next
	}

	var locations *pokedexLocations = getPokedexLocations(url, config)

	var next, previous *string = &locations.Next, &locations.Previous

	config.Next = next
	if *previous == "" {
		config.Prev = nil
	} else {
		config.Prev = previous
	}

	for _, location := range locations.Results {
		fmt.Println(location.Name)
	}

	return errors.New("something went wrong in map")
}

// commandMapb will go back 20 location areas in the world,
// a mehtod to go back, if you're on the first page, prints error.
func commandMapb(config *RequestConfig) error {

	// Checking is previous is nil or empty string
	if config.Prev == nil {
		fmt.Println("you are on the first page, can't go back, try going forward using `map`")
	} else {
		url := *config.Prev
		var locations *pokedexLocations = getPokedexLocations(url, config)

		var next, previous string = locations.Next, locations.Previous

		config.Next = &next
		config.Prev = &previous

		for _, location := range locations.Results {
			fmt.Println(location.Name)
		}
	}

	return errors.New("something went wrong in mapb")
}

func getPokedexLocations(url string, config *RequestConfig) *pokedexLocations {
	var locations pokedexLocations

	res, err := http.Get(url)
	if err != nil {
		log.Fatalln("something went wrong while fetching")
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()

	if res.StatusCode > 299 {
		log.Fatalf("Response failed with status code %d and body\n: %s\n", res.StatusCode, body)
	}
	if err != nil {
		log.Fatalln("something went wrong while fetching")
	}
	err = json.Unmarshal(body, &locations)
	if err != nil {
		log.Fatalln("something went wrong while unmarshalling")
	}
	return &locations
}

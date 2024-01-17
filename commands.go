package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/munanadi/pokedex/pokehelp"
)

func CommandExit(config *pokehelp.RequestConfig, args ...[]string) error {
	os.Exit(0)
	return errors.New("something went wrong in exit")
}

func CommandHelp(config *pokehelp.RequestConfig, args ...[]string) error {
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
func CommandMap(config *pokehelp.RequestConfig, args ...[]string) error {
	// https://pokeapi.co/api/v2/location/{id or name}/ API to hit.

	// Check if next exists and then make a call to that.
	// Start at 0 always
	url := ""
	if config.Next == nil {
		fmt.Println("starting from first")
		url = "https://pokeapi.co/api/v2/location-area/?offset=0&limit=20"
	} else {
		url = *config.Next
	}

	var data []byte
	// Check in cache
	if v, ok := config.Cache.Get(url); !ok {
		fmt.Println("not in cache, fetching..")
		data, _ = pokehelp.GetBodyFromUrl(url, config)
		config.Cache.Add(url, data)
	} else {
		fmt.Println("found in cache..")
		data = v
	}

	var locations *pokehelp.PokedexLocations
	err := json.Unmarshal(data, &locations)
	if err != nil {
		log.Fatalln("unmarshalling body data failed")
	}

	var next, previous *string = &locations.Next, &locations.Previous

	config.Next = next
	if len(*previous) == 0 {
		fmt.Println("prev set as nil")
		config.Prev = nil
	} else {
		fmt.Println("prev set as ", *previous)
		config.Prev = previous
	}

	for _, location := range locations.Results {
		fmt.Println(location.Name)
	}

	return errors.New("something went wrong in map")
}

// commandMapb will go back 20 location areas in the world,
// a mehtod to go back, if you're on the first page, prints error.
func CommandMapb(config *pokehelp.RequestConfig, args ...[]string) error {

	// Checking is previous is nil or empty string
	if config.Prev == nil {
		fmt.Println("you are on the first page, can't go back, try going forward using `map`")
	} else {
		url := *config.Prev
		var locations *pokehelp.PokedexLocations

		fmt.Println(url, " is the url for prev")

		var data []byte
		// Check in cache
		if v, ok := config.Cache.Get(url); !ok {
			fmt.Println("not in cache, fetching..")
			data, _ = pokehelp.GetBodyFromUrl(url, config)
			config.Cache.Add(url, data)
		} else {
			fmt.Println("found in cache..")
			data = v
		}

		err := json.Unmarshal(data, &locations)
		if err != nil {
			log.Fatalln("unmarshalling body data failed")
		}

		var next, previous string = locations.Next, locations.Previous

		config.Next = &next
		if len(previous) == 0 {
			config.Prev = nil
		} else {
			config.Prev = &previous
		}

		for _, location := range locations.Results {
			fmt.Println(location.Name)
		}
	}

	return errors.New("something went wrong in mapb")
}

func CommandExplore(config *pokehelp.RequestConfig, args ...[]string) error {

	cityAreaToExplore := strings.Join(args[0], "")
	fmt.Printf("Exploring %s...\n", cityAreaToExplore)

	baseUrl := "https://pokeapi.co/api/v2/location-area/"
	url := baseUrl + cityAreaToExplore

	var res *pokehelp.PokedexLocationExplore

	var data []byte
	// Check in cache
	if v, ok := config.Cache.Get(url); !ok {
		fmt.Println("not in cache, fetching..")
		data, _ = pokehelp.GetBodyFromUrl(url, config)
		config.Cache.Add(url, data)
	} else {
		fmt.Println("found in cache..")
		data = v
	}

	json.Unmarshal(data, &res)

	fmt.Println("Found Pokemon:")
	for _, v := range res.PokemonEncounters {
		fmt.Printf("- %s\n", v.Pokemon.Name)
	}

	return errors.New("something went wrong in explore")
}

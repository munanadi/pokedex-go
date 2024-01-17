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
	"strings"
	"time"

	"github.com/munanadi/pokedex/pokecache"
)

type cliCommand struct {
	name        string
	description string
	callback    func(config *RequestConfig, args ...[]string) error
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

type pokedexLocationExplore struct {
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int `json:"rate"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	GameIndex int `json:"game_index"`
	ID        int `json:"id"`
	Location  struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Name  string `json:"name"`
	Names []struct {
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
		Name string `json:"name"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			EncounterDetails []struct {
				Chance          int   `json:"chance"`
				ConditionValues []any `json:"condition_values"`
				MaxLevel        int   `json:"max_level"`
				Method          struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"method"`
				MinLevel int `json:"min_level"`
			} `json:"encounter_details"`
			MaxChance int `json:"max_chance"`
			Version   struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}

type RequestConfig struct {
	// *string cause it can be nill too
	Next  *string
	Prev  *string
	Cache *pokecache.Cache
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
		"explore": {
			name:        "explore",
			description: "Let's you explore a city area",
			callback:    commandExplore,
		},
	}
}

func main() {

	const CACHE_REFRESH_IN_SECONDS int64 = 20

	// Store values in cache for 10s
	timeInterval := time.Duration(CACHE_REFRESH_IN_SECONDS) * time.Second
	cache := pokecache.NewCache(timeInterval)
	fmt.Println(cache)

	config := &RequestConfig{Next: nil, Prev: nil, Cache: cache}

	for {
		fmt.Printf("Pokedex > ")
		scanner := bufio.NewScanner(os.Stdin)
		if ok := scanner.Scan(); ok {
			for k, v := range getCommands() {
				args := strings.Split(scanner.Text(), " ")
				if k == args[0] {
					v.callback(config, args[1:])
				}
			}
		}
	}
}

func commandExit(config *RequestConfig, args ...[]string) error {
	os.Exit(0)
	return errors.New("something went wrong in exit")
}

func commandHelp(config *RequestConfig, args ...[]string) error {
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
func commandMap(config *RequestConfig, args ...[]string) error {
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
		data, _ = getBodyFromUrl(url, config)
		config.Cache.Add(url, data)
	} else {
		fmt.Println("found in cache..")
		data = v
	}

	var locations *pokedexLocations
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
func commandMapb(config *RequestConfig, args ...[]string) error {

	// Checking is previous is nil or empty string
	if config.Prev == nil {
		fmt.Println("you are on the first page, can't go back, try going forward using `map`")
	} else {
		url := *config.Prev
		var locations *pokedexLocations

		fmt.Println(url, " is the url for prev")

		var data []byte
		// Check in cache
		if v, ok := config.Cache.Get(url); !ok {
			fmt.Println("not in cache, fetching..")
			data, _ = getBodyFromUrl(url, config)
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

func commandExplore(config *RequestConfig, args ...[]string) error {

	cityAreaToExplore := strings.Join(args[0], "")
	fmt.Printf("Exploring %s...\n", cityAreaToExplore)

	baseUrl := "https://pokeapi.co/api/v2/location-area/"
	url := baseUrl + cityAreaToExplore

	var res *pokedexLocationExplore

	var data []byte
	// Check in cache
	if v, ok := config.Cache.Get(url); !ok {
		fmt.Println("not in cache, fetching..")
		data, _ = getBodyFromUrl(url, config)
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

func getBodyFromUrl(url string, config *RequestConfig) ([]byte, error) {

	res, err := http.Get(url)
	if err != nil {
		log.Fatalln("something went wrong while fetching")
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()

	// TODO: Handle 404 and other stuff
	if res.StatusCode > 299 {
		log.Fatalf("Response failed with status code %d and body\n: %s\n", res.StatusCode, body)
	}
	if err != nil {
		log.Fatalln("something went wrong while fetching")
	}

	return body, nil
}

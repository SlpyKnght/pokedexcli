package main

import (
	"bufio"
	"math/rand"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/slpyknght/pokedex/pokecache"
)

type cliCommand struct {
	name        string
	description string
	callback    func(conf *config) error
}
type config struct {
	next string
	previous string
	cache *pokecache.Cache
	pokedex map[string]Pokemon
	input []string
}

type ApiBatchResponse[T any] struct{
	Count int `json:"count"`
	Next string `json:"next"`
	Previous string `json:"previous"`
	Results T `json:"results"`
}

type Location struct{
	Name string `json:"name"`
	Url string `json:"url"`
}

type ApiResponse struct {
	ID                   int                    `json:"id"`
	Name                 string                 `json:"name"`
	GameIndex            int                    `json:"game_index"`
	EncounterMethodRates []EncounterMethodRates `json:"encounter_method_rates"`
	Location             Location               `json:"location"`
	Names                []Names                `json:"names"`
	PokemonEncounters    []PokemonEncounters    `json:"pokemon_encounters"`
}
type EncounterMethod struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
type Version struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
type EncounterMethodRates struct {
	EncounterMethod EncounterMethod  `json:"encounter_method"`
	VersionDetails  []VersionDetails `json:"version_details"`
}
type Language struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
type Names struct {
	Name     string   `json:"name"`
	Language Language `json:"language"`
}
type Pokemon struct {
	Id	int `json:"id,omitempty"`
	Name string `json:"name"`
	URL  string `json:"url"`
	BaseExperience         int           `json:"base_experience,omitempty"`
	Height                 int           `json:"height,omitempty"`
	IsDefault              bool          `json:"is_default,omitempty"`
	Order                  int           `json:"order,omitempty"`
	Weight                 int           `json:"weight,omitempty"`
	Stats                  []Stats       `json:"stats,omitempty"`
	Types                  []Types       `json:"types,omitempty"`
}
type Stat struct {
	Name string `json:"name,omitempty"`
	URL  string `json:"url,omitempty"`
}
type Stats struct {
	BaseStat int  `json:"base_stat,omitempty"`
	Effort   int  `json:"effort,omitempty"`
	Stat     Stat `json:"stat,omitempty"`
}
type Type struct {
	Name string `json:"name,omitempty"`
	URL  string `json:"url,omitempty"`
}
type Types struct {
	Slot int  `json:"slot,omitempty"`
	Type Type `json:"type,omitempty"`
}
type Method struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
type EncounterDetails struct {
	MinLevel        int    `json:"min_level"`
	MaxLevel        int    `json:"max_level"`
	ConditionValues []any  `json:"condition_values"`
	Chance          int    `json:"chance"`
	Method          Method `json:"method"`
}
type VersionDetails struct {
	Rate    int     `json:"rate,omitempty"`
	Version          Version            `json:"version,omitempty"`
	MaxChance        int                `json:"max_chance,omitempty"`
	EncounterDetails []EncounterDetails `json:"encounter_details,omitempty"`
}
type PokemonEncounters struct {
	Pokemon        Pokemon          `json:"pokemon"`
	VersionDetails []VersionDetails `json:"version_details"`
}

func main(){
	commands := getCommandRegistry()
	scanner := bufio.NewScanner(os.Stdin)
	config := config{
		cache: pokecache.NewCache(time.Second * 5),
		pokedex: make(map[string]Pokemon),
	}
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		rawInput := scanner.Text()
		input := cleanInput(rawInput)
		if len(input) == 0 {
			continue
		}
		config.input = input
		if cmd, ok := commands[input[0]];ok{
			if err := cmd.callback(&config); err != nil{
				fmt.Printf("%v\n", err)
			}
			
		}
	}
}

func cleanInput(text string) []string{
	text = strings.ToLower(text)
	text = strings.TrimSpace(text)
	result := make([]string,0)	
	for _, w := range strings.Fields(text){
		if w == ""{
			continue
		}
		result = append(result, w)
	}

	return result
}

func commandCatch(conf *config)error{
	if len(conf.input) < 2{
		return fmt.Errorf("missing pokemon parameter")
	}
	fmt.Println("Throwing a Pokeball at " + conf.input[1] + "...")
	return catchPokemon(conf.input[1], conf)
}

func commandExplore(conf *config)error{
	if len(conf.input) < 2{
		return fmt.Errorf("missing location parameter")
	}
	location := conf.input[1]
	return fetchPokemonIn(location, conf)
}

func commandInspect(conf *config)error{
	if len(conf.input) < 2{
		return fmt.Errorf("missing location parameter")
	}
	pokemon, ok:= conf.pokedex[conf.input[1]];
	if !ok{
		return fmt.Errorf("%v not caught", conf.input[1])
	}
	fmt.Printf("Name: %v\n", pokemon.Name)
	fmt.Printf("Height: %v\n", pokemon.Height)
	fmt.Printf("Stats:\n")
	for _, stat := range pokemon.Stats{
		fmt.Printf("\t-%v: %v\n", stat.Stat.Name, stat.BaseStat)
	}
	fmt.Printf("Types:\n")
	for _, t := range pokemon.Types{
		fmt.Printf("\t-%v\n", t.Type.Name)
	}
	return nil;
}

func commandExit(conf *config) error{
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(conf *config) error{
	cmds := getCommandRegistry()
	fmt.Printf("Welcome to the Pokedex!\nUsage:\n\n")
	for _, c := range cmds{
		fmt.Printf("%s: %s\n", c.name, c.description)
	}
	return nil
}

func commandMapBack(conf *config) error{
	url := conf.previous
	if url == "" {
		return nil
	}
	return fetchMap(url, conf)
}
func commandMap(conf *config) error{
	url := "https://pokeapi.co/api/v2/location-area?limit=20&offset=0"
	if conf.next != ""{
		url = conf.next
	}
	return fetchMap(url, conf)
}

func commandPokedex(conf *config)error{
	if len(conf.pokedex) == 0 {
		return fmt.Errorf("Go out there and catch some!")
	}
	fmt.Println("Your Pokedex:")
	for name := range conf.pokedex{
		fmt.Printf("\t-%s \n", name)
	}
	return nil
}

func fetchMap(url string, conf *config)error{
	var rawResponse []byte
	if val, ok:= conf.cache.Get(url); ok{
		fmt.Println("cached result:")
		rawResponse = val
	}else{
		res, err := http.Get(url)
		if err != nil{
			return err
		}
		defer res.Body.Close()
		rawResponse,_ = io.ReadAll(res.Body)
		conf.cache.Add(url, rawResponse)
	}
	apiResp := ApiBatchResponse[[]Location]{}
	json.Unmarshal(rawResponse, &apiResp)
	conf.next = apiResp.Next
	conf.previous = apiResp.Previous
	for _, loc := range apiResp.Results{
		fmt.Println(loc.Name)
	}
	return nil
}

func fetchPokemonIn(location string, conf *config)error{
	var rawResponse []byte
	url := "https://pokeapi.co/api/v2/location-area/" + location
	if val, ok:= conf.cache.Get(url); ok{
		fmt.Println("cached result:")
		rawResponse = val
	}else{
		res, err := http.Get(url)
		if err != nil{
			return err
		}
		defer res.Body.Close()
		rawResponse,_ = io.ReadAll(res.Body)
		conf.cache.Add(url, rawResponse)
	}
	apiResp := ApiResponse{}
	json.Unmarshal(rawResponse, &apiResp)
	for _, encounter := range apiResp.PokemonEncounters{
		fmt.Println(encounter.Pokemon.Name)
	}
	return nil
}

func catchPokemon(pokemon string, conf *config)error{
	var rawResponse []byte
	url := "https://pokeapi.co/api/v2/pokemon/" + pokemon
	if val, ok:= conf.cache.Get(url); ok{
		fmt.Println("cached result:")
		rawResponse = val
	}else{
		res, err := http.Get(url)
		if err != nil{
			return err
		}
		defer res.Body.Close()
		rawResponse,_ = io.ReadAll(res.Body)
		conf.cache.Add(url, rawResponse)
	}
	apiResp := Pokemon{}
	json.Unmarshal(rawResponse, &apiResp)
	if apiResp.Name == ""{
		return fmt.Errorf("You can only catch pokemon...")
	}
	if rand.Int() % 2 == 0{
		fmt.Println(apiResp.Name + " escaped!")
	}else{
		conf.pokedex[apiResp.Name] = apiResp
		fmt.Println(apiResp.Name + " was caught!")
	}
	return nil;
}

func getCommandRegistry()map[string]cliCommand{
	return map[string]cliCommand{
		"help":{
			name: "help",
			description: "Displays a help message",
			callback: commandHelp,
		},
		"pokedex":{
			name: "pokedex",
			description: "list all pokemon you have caught so far",
			callback: commandPokedex,
		},
		"inspect":{
			name: "inspect",
			description: "inspect a pokemon you have caught",
			callback: commandInspect,
		},
		"catch":{
			name: "catch",
			description: "gotta catch them all",
			callback: commandCatch,
		},
		"explore":{
			name: "explore",
			description: "Find all Pokemon in this location",
			callback: commandExplore,
		},
		"map":{
			name: "map",
			description: "Display 20 locations of the Pokemon world",
			callback: commandMap,
		},
		"mapb":{
			name: "mapb",
			description: "Display the previous 20 locations",
			callback: commandMapBack,
		},
		"exit":{
			name: "exit",
			description: "Exit the Pokedex",
			callback: commandExit,
		},
	}
}

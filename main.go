package main

import (
	"bufio"
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
}

type ApiResponse[T any] struct{
	Count int `json:"count"`
	Next string `json:"next"`
	Previous string `json:"previous"`
	Results T `json:"results"`
}

type Location struct{
	Name string `json:"name"`
	Url string `json:"url"`
}

func main(){
	commands := getCommandRegistry()
	scanner := bufio.NewScanner(os.Stdin)
	config := config{cache: pokecache.NewCache(time.Second * 5)}
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		rawInput := scanner.Text()
		input := cleanInput(rawInput)
		if len(input) == 0 {
			continue
		}
		if cmd, ok := commands[input[0]];ok{
			cmd.callback(&config)
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
	apiResp := ApiResponse[[]Location]{}
	json.Unmarshal(rawResponse, &apiResp)
	conf.next = apiResp.Next
	conf.previous = apiResp.Previous
	for _, loc := range apiResp.Results{
		fmt.Println(loc.Name)
	}
	return nil
}

func getCommandRegistry()map[string]cliCommand{
	return map[string]cliCommand{
		"help":{
			name: "help",
			description: "Displays a help message",
			callback: commandHelp,
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

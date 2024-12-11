package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main(){
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		rawInput := scanner.Text()
		input := cleanInput(rawInput)
		if len(input) > 0 {
			fmt.Printf("Your command was %s\n", input[0])
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

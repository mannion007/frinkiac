package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const (
	baseURL                 = "https://frinkiac.com"
	wordsBeforeBreakingLow  = 4
	wordsBeforeBreakingMid  = 3
	wordsBeforeBreakingHigh = 3
)

type Result struct {
	ID        int    `json:"Id"`
	Episode   string `json:"Episode"`
	Timestamp int    `json:"Timestamp"`
}

func main() {
	if len(os.Args) < 2 || len(os.Args) > 3 {
		fmt.Println("Usage: frinkiac <query> [caption]")
		return
	}

	rand.Seed(time.Now().UnixNano())
	query := os.Args[1]

	result, err := searchFrinkiac(query)
	if err != nil {
		log.Fatalf("Error: %v", err)
		return
	}

	imageURL := getImageURL(result.Episode, result.Timestamp, os.Args)
	markdown := getImageMarkdown(imageURL)
	fmt.Println(markdown)
}

func searchFrinkiac(query string) (Result, error) {
	searchURL := fmt.Sprintf("%s/api/search?q=%s", baseURL, encode(query))
	response, err := http.Get(searchURL)
	if err != nil {
		return Result{}, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return Result{}, fmt.Errorf("D'oh! I couldn't find anything for %s", query)
	}

	var results []Result
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return Result{}, err
	}

	if err := json.Unmarshal(body, &results); err != nil {
		return Result{}, err
	}

	if len(results) == 0 {
		return Result{}, fmt.Errorf("D'oh! I couldn't find anything for %s", query)
	}

	return results[rand.Intn(len(results))], nil
}

func getImageURL(episode string, timestamp int, args []string) string {
	if len(args) == 3 {
		return getImageURLWithCaption(episode, timestamp, formatCaption(args[2]))
	}
	return fmt.Sprintf("%s/meme/%s/%d.jpg", baseURL, episode, timestamp)
}

func encode(str string) string {
	return url.QueryEscape(str)
}

func getImageURLWithCaption(episode string, timestamp int, caption string) string {
	return fmt.Sprintf("%s/meme/%s/%d.jpg?lines=%s", baseURL, episode, timestamp, encode(caption))
}

func getImageMarkdown(imageURL string) string {
	return fmt.Sprintf("![image](%s)", imageURL)
}

func formatCaption(caption string) string {
	return addLineBreaks(trimWhitespace(caption))
}

func getNumberOfWordsBeforeBreaking(words []string) int {
	longestWordLength := 0
	for _, word := range words {
		if len(word) > longestWordLength {
			longestWordLength = len(word)
		}
	}

	if longestWordLength <= 5 {
		return wordsBeforeBreakingLow
	} else if longestWordLength <= 8 {
		return wordsBeforeBreakingMid
	}
	return wordsBeforeBreakingHigh
}

func addLineBreaks(str string) string {
	words := strings.Split(str, " ")
	wordsBeforeBreaking := getNumberOfWordsBeforeBreaking(words)
	var newString string

	for i, word := range words {
		i++
		delimiter := " "
		if i%wordsBeforeBreaking == 0 {
			delimiter = "\n"
		}
		newString += word + delimiter
	}

	return newString
}

func trimWhitespace(str string) string {
	return strings.TrimSpace(str)
}

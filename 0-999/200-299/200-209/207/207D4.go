package main

import (
    "bufio"
    "fmt"
    "os"
    "strings"
)

func main() {
    in := bufio.NewScanner(os.Stdin)
    in.Buffer(make([]byte, 1024), 1<<20)

    words := []string{}
    for in.Scan() {
        line := strings.ToLower(in.Text())
        parts := strings.Fields(line)
        words = append(words, parts...)
    }

    counts := []string{ "trade", "economy", "market", "import", "export", "company", "profit" }
    politics := []string{ "minister", "government", "president", "election", "policy", "law", "parliament" }
    culture := []string{ "art", "music", "cinema", "culture", "festival", "literature", "theatre", "movie" }

    scoreTrade := 0
    scorePolitics := 0
    scoreCulture := 0

    for _, w := range words {
        for _, key := range counts {
            if strings.Contains(w, key) {
                scoreTrade++
            }
        }
        for _, key := range politics {
            if strings.Contains(w, key) {
                scorePolitics++
            }
        }
        for _, key := range culture {
            if strings.Contains(w, key) {
                scoreCulture++
            }
        }
    }

    if scoreTrade >= scorePolitics && scoreTrade >= scoreCulture {
        fmt.Println(3)
    } else if scorePolitics >= scoreTrade && scorePolitics >= scoreCulture {
        fmt.Println(2)
    } else {
        fmt.Println(1)
    }
}

package main

import (
    "bufio"
    "fmt"
    "os"
    "sort"
)

func main() {
    in := bufio.NewReader(os.Stdin)
    scores := make([]int, 3)

    for i := range scores {
        if _, err := fmt.Fscan(in, &scores[i]); err != nil {
            return
        }
    }

    maxScore, minScore := scores[0], scores[0]
    for _, v := range scores[1:] {
        if v > maxScore {
            maxScore = v
        }
        if v < minScore {
            minScore = v
        }
    }

    if maxScore-minScore >= 10 {
        fmt.Println("check again")
        return
    }

    sort.Ints(scores)
    fmt.Printf("final %d\n", scores[1])
}

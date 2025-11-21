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
    for i := 0; i < 3; i++ {
        if _, err := fmt.Fscan(in, &scores[i]); err != nil {
            return
        }
    }

    minScore, maxScore := scores[0], scores[0]
    for _, v := range scores[1:] {
        if v < minScore {
            minScore = v
        }
        if v > maxScore {
            maxScore = v
        }
    }

    if maxScore-minScore >= 10 {
        fmt.Println("check again")
        return
    }

    // median for 3 items by sorting
    sort.Ints(scores)
    fmt.Printf("final %d\n", scores[1])
}

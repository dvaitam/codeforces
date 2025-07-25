package main

import "fmt"

func main() {
    var x int
    if _, err := fmt.Scan(&x); err != nil {
        return
    }
    bestRank := -1
    bestA := 0
    var bestCat byte
    for a := 0; a <= 2; a++ {
        r := (x + a) % 4
        rank := 0
        var cat byte
        switch r {
        case 1:
            rank = 4
            cat = 'A'
        case 3:
            rank = 3
            cat = 'B'
        case 2:
            rank = 2
            cat = 'C'
        case 0:
            rank = 1
            cat = 'D'
        }
        if rank > bestRank {
            bestRank = rank
            bestA = a
            bestCat = cat
        }
    }
    fmt.Printf("%d %c\n", bestA, bestCat)
}

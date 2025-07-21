package main

import (
    "bufio"
    "fmt"
    "os"
    "sort"
)

// laptop represents a laptop with price and quality
type laptop struct {
    price   int
    quality int
}

func main() {
    reader := bufio.NewReader(os.Stdin)
    writer := bufio.NewWriter(os.Stdout)
    defer writer.Flush()

    var n int
    if _, err := fmt.Fscan(reader, &n); err != nil {
        return
    }
    laptops := make([]laptop, n)
    for i := 0; i < n; i++ {
        fmt.Fscan(reader, &laptops[i].price, &laptops[i].quality)
    }

    // Sort laptops by increasing price
    sort.Slice(laptops, func(i, j int) bool {
        return laptops[i].price < laptops[j].price
    })

    // Check for a pair where a cheaper laptop has higher quality
    for i := 1; i < n; i++ {
        if laptops[i-1].quality > laptops[i].quality {
            fmt.Fprintln(writer, "Happy Alex")
            return
        }
    }
    fmt.Fprintln(writer, "Poor Alex")
}

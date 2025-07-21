package main

import (
    "fmt"
)

func main() {
    var r, g, b int64
    if _, err := fmt.Scan(&r, &g, &b); err != nil {
        return
    }
    // initial bouquets by single colors
    base := r/3 + g/3 + b/3
    maxBouquets := base
    // try using k mixing bouquets (k = 1,2)
    // mixing bouquets use one of each color
    for k := int64(1); k <= 2; k++ {
        if r < k || g < k || b < k {
            break
        }
        cur := k + (r-k)/3 + (g-k)/3 + (b-k)/3
        if cur > maxBouquets {
            maxBouquets = cur
        }
    }
    fmt.Println(maxBouquets)
}

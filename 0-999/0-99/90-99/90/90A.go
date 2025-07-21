package main

import "fmt"

func main() {
    var r, g, b int
    if _, err := fmt.Scan(&r, &g, &b); err != nil {
        return
    }
    var ans int
    // Red cablecars (offset 0)
    if r > 0 {
        cars := (r + 1) / 2
        time := 0 + (cars-1)*3 + 30
        if time > ans {
            ans = time
        }
    }
    // Green cablecars (offset 1)
    if g > 0 {
        cars := (g + 1) / 2
        time := 1 + (cars-1)*3 + 30
        if time > ans {
            ans = time
        }
    }
    // Blue cablecars (offset 2)
    if b > 0 {
        cars := (b + 1) / 2
        time := 2 + (cars-1)*3 + 30
        if time > ans {
            ans = time
        }
    }
    fmt.Println(ans)
}

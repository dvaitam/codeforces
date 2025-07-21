package main

import "fmt"

func main() {
    var n int
    // number of stations
    if _, err := fmt.Scan(&n); err != nil {
        return
    }
    d := make([]int, n+1)
    // distances between neighboring stations (1-based)
    for i := 1; i <= n; i++ {
        fmt.Scan(&d[i])
    }
    var s, t int
    fmt.Scan(&s, &t)
    // if same station, distance is zero
    if s == t {
        fmt.Println(0)
        return
    }
    // ensure s < t for forward calculation
    if s > t {
        s, t = t, s
    }
    // compute clockwise distance from s to t
    clockwise := 0
    for i := s; i < t; i++ {
        clockwise += d[i]
    }
    // total perimeter
    total := 0
    for i := 1; i <= n; i++ {
        total += d[i]
    }
    // counter-clockwise is the remaining
    counter := total - clockwise
    // output the minimum
    if clockwise < counter {
        fmt.Println(clockwise)
    } else {
        fmt.Println(counter)
    }
}

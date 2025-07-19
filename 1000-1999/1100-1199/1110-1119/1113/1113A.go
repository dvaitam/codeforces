package main

import "fmt"

func main() {
    var n, v int
    if _, err := fmt.Scan(&n, &v); err != nil {
        return
    }
    // If tank capacity is enough to travel all segments without refueling
    if v >= n-1 {
        fmt.Println(n - 1)
    } else {
        // Otherwise, buy v liters at city 1, then additional cheaper liters in subsequent cities
        m := n - v
        // Base cost: v liters at price 1
        cost := v
        // Additional cost: sum of prices from city 2 to city m
        // sum_{i=2..m} i = m*(m+1)/2 - 1
        cost += m*(m+1)/2 - 1
        fmt.Println(cost)
    }
}

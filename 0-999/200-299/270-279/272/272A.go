package main

import "fmt"

func main() {
    var n int
    if _, err := fmt.Scan(&n); err != nil {
        return
    }
    sum := 0
    for i := 0; i < n; i++ {
        var x int
        fmt.Scan(&x)
        sum += x
    }
    cnt := 0
    totalPeople := n + 1
    for k := 1; k <= 5; k++ {
        if (sum + k) % totalPeople != 1 {
            cnt++
        }
    }
    fmt.Println(cnt)
}

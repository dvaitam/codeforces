package main

import "fmt"

func main() {
    var n int
    if _, err := fmt.Scan(&n); err != nil {
        return
    }
    lucky := []int{
        4, 7,
        44, 47, 74, 77,
        444, 447, 474, 477, 744, 747, 774, 777,
    }
    for _, v := range lucky {
        if n%v == 0 {
            fmt.Println("YES")
            return
        }
    }
    fmt.Println("NO")
}

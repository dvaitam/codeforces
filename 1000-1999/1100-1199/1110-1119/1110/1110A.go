package main

import "fmt"

func main() {
    var b, k int
    if _, err := fmt.Scan(&b, &k); err != nil {
        return
    }
    arr := make([]int, k)
    for i := 0; i < k; i++ {
        fmt.Scan(&arr[i])
    }
    if b%2 == 0 {
        if arr[k-1]%2 == 1 {
            fmt.Println("odd")
        } else {
            fmt.Println("even")
        }
    } else {
        sum := 0
        for _, v := range arr {
            if v%2 == 1 {
                sum++
            }
        }
        if sum%2 == 1 {
            fmt.Println("odd")
        } else {
            fmt.Println("even")
        }
    }
}

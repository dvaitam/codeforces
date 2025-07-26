package main

import (
    "bufio"
    "fmt"
    "os"
)

func factorial(n int) int {
    res := 1
    for i := 2; i <= n; i++ {
        res *= i
    }
    return res
}

func main() {
    in := bufio.NewReader(os.Stdin)
    var n int
    if _, err := fmt.Fscan(in, &n); err != nil {
        return
    }
    fmt.Print(factorial(n))
}


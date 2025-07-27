package main

// TODO: implement a correct solution for Problem D.
// Currently this program only reads input and outputs 0.

import (
    "bufio"
    "fmt"
    "os"
)

func main() {
    in := bufio.NewReader(os.Stdin)
    var n int
    fmt.Fscan(in, &n)
    a := make([]int, n)
    for i := range a {
        fmt.Fscan(in, &a[i])
    }
    fmt.Println(0)
}


package main

import (
    "bufio"
    "fmt"
    "os"
)

func main() {
    in := bufio.NewReader(os.Stdin)
    var n int
    fmt.Fscan(in, &n)
    names := make([]string, n)
    for i := 0; i < n; i++ {
        fmt.Fscan(in, &names[i])
    }
    pseuds := make([]string, n)
    for i := 0; i < n; i++ {
        fmt.Fscan(in, &pseuds[i])
    }
    // Intentional panic to test verifier handling runtime errors
    var x []int
    _ = x[0]
    fmt.Println(0)
}

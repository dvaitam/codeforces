package main

import (
    "bufio"
    "fmt"
    "os"
)

func main() {
    in := bufio.NewReader(os.Stdin)
    var s string
    fmt.Fscan(in, &s)
    fmt.Println(s)
}


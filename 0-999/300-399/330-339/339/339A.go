package main

import (
    "bufio"
    "fmt"
    "os"
    "sort"
    "strings"
)

func main() {
    reader := bufio.NewReader(os.Stdin)
    var s string
    if _, err := fmt.Fscan(reader, &s); err != nil {
        return
    }
    parts := strings.Split(s, "+")
    sort.Strings(parts)
    result := strings.Join(parts, "+")
    fmt.Println(result)
}

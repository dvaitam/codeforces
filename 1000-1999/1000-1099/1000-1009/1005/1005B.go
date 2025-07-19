package main

import (
    "bufio"
    "fmt"
    "os"
)

func main() {
    reader := bufio.NewReader(os.Stdin)
    var a, b string
    // Read two strings from input
    if _, err := fmt.Fscan(reader, &a); err != nil {
        return
    }
    if _, err := fmt.Fscan(reader, &b); err != nil {
        return
    }
    // Count matching suffix length
    i, j := len(a)-1, len(b)-1
    match := 0
    for i >= 0 && j >= 0 {
        if a[i] == b[j] {
            match++
            i--
            j--
        } else {
            break
        }
    }
    // Minimum operations: delete non-matching parts
    result := len(a) + len(b) - 2*match
    fmt.Println(result)
}

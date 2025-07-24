package main

import (
    "bufio"
    "fmt"
    "os"
)

func main() {
    reader := bufio.NewReader(os.Stdin)
    // Read the input string
    var s string
    fmt.Fscan(reader, &s)

    n := len(s)
    // Convert string to byte slice for modification
    b := []byte(s)

    // Find first character not 'a'
    i := 0
    for i < n && b[i] == 'a' {
        i++
    }

    if i == n {
        // All 'a's: shift last character only
        b[n-1] = 'z'
    } else {
        // From i, shift until we hit an 'a' or end
        j := i
        for j < n && b[j] != 'a' {
            // Shift letter by one backwards
            b[j]--
            j++
        }
    }

    // Output result
    writer := bufio.NewWriter(os.Stdout)
    defer writer.Flush()
    writer.Write(b)
}

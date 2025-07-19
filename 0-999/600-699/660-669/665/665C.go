package main

import (
    "bufio"
    "fmt"
    "os"
)

func main() {
    reader := bufio.NewReader(os.Stdin)
    var s string
    if _, err := fmt.Fscan(reader, &s); err != nil {
        return
    }
    b := []byte(s)
    n := len(b)
    for i := 1; i < n; i++ {
        if b[i] == b[i-1] {
            for c := byte('a'); c <= 'z'; c++ {
                if c != b[i] && (i+1 >= n || c != b[i+1]) {
                    b[i] = c
                    break
                }
            }
        }
    }
    writer := bufio.NewWriter(os.Stdout)
    defer writer.Flush()
    writer.Write(b)
}

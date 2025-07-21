package main

import (
    "bufio"
    "fmt"
    "os"
    "strings"
)

func main() {
    reader := bufio.NewReader(os.Stdin)
    s, err := reader.ReadString('\n')
    if err != nil && err.Error() != "EOF" {
        fmt.Fprintln(os.Stderr, "Error reading input:", err)
        return
    }
    s = strings.TrimSpace(s)
    if strings.Contains(s, "666") {
        fmt.Println("YES")
    } else {
        fmt.Println("NO")
    }
}

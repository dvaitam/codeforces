package main

import (
    "bufio"
    "fmt"
    "os"
    "strconv"
    "strings"
)

func main() {
    reader := bufio.NewReader(os.Stdin)
    // Read number of lines
    line, err := reader.ReadString('\n')
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        return
    }
    line = strings.TrimSpace(line)
    n, err := strconv.Atoi(line)
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        return
    }
    // Read verse pattern
    line, err = reader.ReadString('\n')
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        return
    }
    parts := strings.Fields(line)
    pattern := make([]int, n)
    for i := 0; i < n && i < len(parts); i++ {
        pattern[i], _ = strconv.Atoi(parts[i])
    }
    // Check each line
    vowels := "aeiouy"
    for i := 0; i < n; i++ {
        text, err := reader.ReadString('\n')
        if err != nil && len(text) == 0 {
            fmt.Fprintln(os.Stderr, err)
            return
        }
        count := 0
        for _, r := range text {
            if strings.ContainsRune(vowels, r) {
                count++
            }
        }
        if count != pattern[i] {
            fmt.Println("NO")
            return
        }
    }
    fmt.Println("YES")
}

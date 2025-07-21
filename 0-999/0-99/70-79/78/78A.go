package main

import (
    "bufio"
    "fmt"
    "os"
)

// countVowels returns the number of vowels ('a', 'e', 'i', 'o', 'u') in s.
func countVowels(s string) int {
    cnt := 0
    for _, ch := range s {
        switch ch {
        case 'a', 'e', 'i', 'o', 'u':
            cnt++
        }
    }
    return cnt
}

func main() {
    scanner := bufio.NewScanner(os.Stdin)
    lines := make([]string, 0, 3)
    for i := 0; i < 3; i++ {
        if !scanner.Scan() {
            // insufficient input
            break
        }
        lines = append(lines, scanner.Text())
    }
    if len(lines) < 3 {
        return
    }
    targets := []int{5, 7, 5}
    ok := true
    for i, line := range lines {
        if countVowels(line) != targets[i] {
            ok = false
            break
        }
    }
    if ok {
        fmt.Println("YES")
    } else {
        fmt.Println("NO")
    }
}

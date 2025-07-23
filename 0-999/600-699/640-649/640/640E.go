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
    line, _ := reader.ReadString('\n')
    line = strings.TrimSpace(line)
    if len(line) == 0 {
        return
    }
    parts := strings.Fields(line)
    arr := make([]int, len(parts))
    for i, s := range parts {
        v, err := strconv.Atoi(s)
        if err != nil {
            return
        }
        arr[i] = v
    }

    result := 0
    for i, val := range arr {
        ok := true
        for j, other := range arr {
            if i == j {
                continue
            }
            if val%other != 0 {
                ok = false
                break
            }
        }
        if ok {
            result = 1
            break
        }
    }
    fmt.Print(result)
}

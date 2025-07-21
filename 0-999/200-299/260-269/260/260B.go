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
    line, err := reader.ReadString('\n')
    if err != nil && err.Error() != "EOF" {
        return
    }
    s := strings.TrimSpace(line)
    n := len(s)
    daysInMonth := []int{31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}
    count := make(map[string]int)
    maxCount := 0
    result := ""
    for i := 0; i+10 <= n; i++ {
        // check format dd-mm-yyyy
        if s[i+2] != '-' || s[i+5] != '-' {
            continue
        }
        sub := s[i : i+10]
        day, err1 := strconv.Atoi(sub[0:2])
        month, err2 := strconv.Atoi(sub[3:5])
        year, err3 := strconv.Atoi(sub[6:10])
        if err1 != nil || err2 != nil || err3 != nil {
            continue
        }
        if year < 2013 || year > 2015 {
            continue
        }
        if month < 1 || month > 12 {
            continue
        }
        if day < 1 || day > daysInMonth[month-1] {
            continue
        }
        count[sub]++
        if count[sub] > maxCount {
            maxCount = count[sub]
            result = sub
        }
    }
    if result != "" {
        fmt.Println(result)
    }
}

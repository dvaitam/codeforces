package main

import "fmt"

func main() {
    var s string
    if _, err := fmt.Scan(&s); err != nil {
        return
    }
    col := int(s[0]-'a') + 1
    row := int(s[1]-'0')
    cnt := 0
    for dx := -1; dx <= 1; dx++ {
        for dy := -1; dy <= 1; dy++ {
            if dx == 0 && dy == 0 {
                continue
            }
            c2 := col + dx
            r2 := row + dy
            if c2 >= 1 && c2 <= 8 && r2 >= 1 && r2 <= 8 {
                cnt++
            }
        }
    }
    fmt.Println(cnt)
}

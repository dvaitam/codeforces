package main

import (
    "bufio"
    "fmt"
    "os"
)

func isPalindrome(h, m int) bool {
    return h/10 == m%10 && h%10 == m/10
}

func main() {
    in := bufio.NewReader(os.Stdin)
    var s string
    if _, err := fmt.Fscan(in, &s); err != nil {
        return
    }
    var hh, mm int
    fmt.Sscanf(s, "%02d:%02d", &hh, &mm)
    ans := 0
    for {
        if isPalindrome(hh, mm) {
            fmt.Println(ans)
            return
        }
        ans++
        mm++
        if mm == 60 {
            mm = 0
            hh = (hh + 1) % 24
        }
    }
}

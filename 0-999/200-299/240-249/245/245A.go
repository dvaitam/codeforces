package main

import (
    "fmt"
)

func main() {
    var n int
    if _, err := fmt.Scan(&n); err != nil {
        return
    }
    var t, x, y int
    var succA, succB, cntA, cntB int
    for i := 0; i < n; i++ {
        if _, err := fmt.Scan(&t, &x, &y); err != nil {
            return
        }
        switch t {
        case 1:
            cntA++
            succA += x
        case 2:
            cntB++
            succB += x
        }
    }
    // A is alive if at least half of sent packets reached
    if succA*2 >= cntA*10 {
        fmt.Println("LIVE")
    } else {
        fmt.Println("DEAD")
    }
    if succB*2 >= cntB*10 {
        fmt.Println("LIVE")
    } else {
        fmt.Println("DEAD")
    }
}

package main

import (
    "fmt"
)

func gcd(a, b int64) int64 {
    if a == 0 {
        return b
    }
    return gcd(b%a, a)
}

func main() {
    var a, b int64
    if _, err := fmt.Scan(&a, &b); err != nil {
        return
    }
    var ans_x, ans_y int64
    flag := false
    x := gcd(a, b)
    if a > b {
        a, b = b, a
    }
    num := b / x
    den := a / x
    if x == 1 {
        fmt.Println("NO")
    } else {
        for i := int64(1); i < a; i++ {
            for j := i; j < a; j++ {
                if j*j+i*i == a*a {
                    ans_x = i
                    ans_y = j
                    if (ans_x*num)%den != 0 || (ans_y*num)%den != 0 {
                        continue
                    }
                    flag = true
                    break
                }
            }
            if flag {
                break
            }
        }
        if !flag {
            fmt.Println("NO")
        } else {
            ans1_x := -(num * ans_y) / den
            ans1_y := (num * ans_x) / den
            if ans1_x == ans_x || ans1_y == ans_y {
                ans1_x = -ans1_x
                ans1_y = -ans1_y
            }
            if ans_x == ans1_x || ans_y == ans1_y {
                fmt.Println("NO")
            } else {
                fmt.Println("YES")
                fmt.Println("0 0")
                fmt.Printf("%d %d\n", ans_x, ans_y)
                fmt.Printf("%d %d\n", ans1_x, ans1_y)
            }
        }
    }
}

package main

import (
    "bufio"
    "fmt"
    "math"
    "os"
)

func gcd(a, b int) int {
    for b != 0 {
        a, b = b, a%b
    }
    return a
}

func unionSize(nums []int) int {
    m := len(nums)
    ans := 0
    for mask := 1; mask < 1<<m; mask++ {
        g := 0
        bits := 0
        for i := 0; i < m; i++ {
            if mask>>i&1 == 1 {
                if g == 0 {
                    g = nums[i]
                } else {
                    g = gcd(g, nums[i])
                }
                bits++
            }
        }
        if bits%2 == 1 {
            ans += g
        } else {
            ans -= g
        }
    }
    return ans
}

func bruteForce(n, k int) int {
    values := make([]int, n-2)
    for i := range values {
        values[i] = i + 3
    }
    choose := make([]int, k)
    best := math.MaxInt32
    var dfs func(start, idx int)
    dfs = func(start, idx int) {
        if idx == k {
            val := unionSize(choose)
            if val < best {
                best = val
            }
            return
        }
        for i := start; i <= len(values)-(k-idx); i++ {
            choose[idx] = values[i]
            dfs(i+1, idx+1)
        }
    }
    dfs(0, 0)
    return best
}

func main() {
    in := bufio.NewReader(os.Stdin)
    var n, k int
    fmt.Fscan(in, &n, &k)

    if n <= 15 {
        fmt.Println(bruteForce(n, k))
        return
    }

    switch k {
    case 1:
        fmt.Println(3)
    case 2:
        fmt.Println(6)
    default:
        if k <= 6 {
            fmt.Println(4*k - 4)
        } else {
            fmt.Println(6*k - 18)
        }
    }
}


package main

import (
    "bytes"
    "context"
    "fmt"
    "math/rand"
    "os"
    "os/exec"
    "strings"
    "time"
)

func abs(x int64) int64 {
    if x < 0 {
        return -x
    }
    return x
}

func min64(a, b int64) int64 {
    if a < b {
        return a
    }
    return b
}

func solve(n, x int, a []int64) int64 {
    if n == 1 {
        v1 := int64(a[0]-1)
        v2 := int64(x-1)
        if v1 < v2 {
            return v2
        }
        return v1
    }
    base := int64(0)
    for i := 0; i < n-1; i++ {
        base += abs(a[i] - a[i+1])
    }
    costVal := func(val int64) int64 {
        res := min64(abs(a[0]-val), abs(a[n-1]-val))
        for i := 0; i < n-1; i++ {
            c := abs(a[i]-val) + abs(a[i+1]-val) - abs(a[i]-a[i+1])
            if c < res {
                res = c
            }
        }
        return res
    }
    ans := base + costVal(1)
    if x > 1 {
        ans += costVal(int64(x))
    }
    cand := base + int64(x-1) + min64(abs(a[0]-1), abs(a[0]-int64(x)))
    if cand < ans {
        ans = cand
    }
    cand = base + int64(x-1) + min64(abs(a[n-1]-1), abs(a[n-1]-int64(x)))
    if cand < ans {
        ans = cand
    }
    for i := 0; i < n-1; i++ {
        cand = base - abs(a[i]-a[i+1]) + int64(x-1) +
            min64(abs(a[i]-1)+abs(a[i+1]-int64(x)), abs(a[i]-int64(x))+abs(a[i+1]-1))
        if cand < ans {
            ans = cand
        }
    }
    return ans
}

func runBinary(bin, input string) (string, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()
    cmd := exec.CommandContext(ctx, bin)
    cmd.Stdin = strings.NewReader(input)
    var out bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = os.Stderr
    err := cmd.Run()
    return strings.TrimSpace(out.String()), err
}

func main() {
    if len(os.Args) < 2 {
        fmt.Println("usage: go run verifierD.go /path/to/binary")
        return
    }
    bin := os.Args[1]
    r := rand.New(rand.NewSource(4))
    tests := 100
    for i := 0; i < tests; i++ {
        n := r.Intn(3) + 1
        x := r.Intn(7) + 1
        a := make([]int64, n)
        for j := 0; j < n; j++ {
            a[j] = int64(r.Intn(7) + 1)
        }
        input := fmt.Sprintf("1\n%d %d\n", n, x)
        for j, v := range a {
            if j > 0 {
                input += " "
            }
            input += fmt.Sprintf("%d", v)
        }
        input += "\n"
        expected := fmt.Sprintf("%d", solve(n, x, a))
        out, err := runBinary(bin, input)
        if err != nil {
            fmt.Printf("test %d runtime error: %v\n", i+1, err)
            return
        }
        out = strings.TrimSpace(out)
        if out != expected {
            fmt.Printf("test %d failed: n=%d x=%d a=%v expected=%s got=%s\n", i+1, n, x, a, expected, out)
            return
        }
    }
    fmt.Println("All tests passed")
}

package main

import (
    "bytes"
    "context"
    "fmt"
    "math/rand"
    "os"
    "os/exec"
    "sort"
    "strings"
    "time"
)

func solve(n int, x int64, a []int64) int64 {
    sort.Slice(a, func(i, j int) bool { return a[i] < a[j] })
    prefix := make([]int64, n+1)
    for i := 1; i <= n; i++ {
        prefix[i] = prefix[i-1] + a[i-1]
    }
    var ans, last int64
    for k := 1; k <= n; k++ {
        if prefix[k] > x {
            break
        }
        days := (x-prefix[k])/int64(k) + 1
        ans += days - last
        last = days
    }
    return ans
}

func runBinary(bin, input string) (string, error) {
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
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
        fmt.Println("usage: go run verifierC.go /path/to/binary")
        return
    }
    bin := os.Args[1]
    r := rand.New(rand.NewSource(3))
    tests := 100
    for i := 0; i < tests; i++ {
        n := r.Intn(4) + 1
        x := int64(r.Intn(20) + 1)
        a := make([]int64, n)
        for j := 0; j < n; j++ {
            a[j] = int64(r.Intn(20) + 1)
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

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

func solve(points []int) string {
    L := -1<<60
    R := 1<<60
    for i, v := range points {
        a := v - (i + 1)
        if a > L {
            L = a
        }
        if a+2 < R {
            R = a + 2
        }
    }
    if L <= R {
        return "YES"
    }
    return "NO"
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
        fmt.Println("usage: go run verifierB.go /path/to/binary")
        return
    }
    bin := os.Args[1]
    r := rand.New(rand.NewSource(2))
    tests := 100
    for i := 0; i < tests; i++ {
        n := r.Intn(5) + 1
        points := make([]int, n)
        cur := r.Intn(5) + 1
        for j := 0; j < n; j++ {
            cur += r.Intn(5) + 1
            points[j] = cur
        }
        input := fmt.Sprintf("1\n%d\n", n)
        for j, v := range points {
            if j > 0 {
                input += " "
            }
            input += fmt.Sprintf("%d", v)
        }
        input += "\n"
        expected := solve(points)
        out, err := runBinary(bin, input)
        if err != nil {
            fmt.Printf("test %d runtime error: %v\n", i+1, err)
            return
        }
        out = strings.ToUpper(strings.TrimSpace(out))
        if out != expected {
            fmt.Printf("test %d failed: input=%v expected=%s got=%s\n", i+1, points, expected, out)
            return
        }
    }
    fmt.Println("All tests passed")
}

package main

import (
    "bufio"
    "bytes"
    "context"
    "fmt"
    "io"
    "math/rand"
    "os"
    "os/exec"
    "strings"
    "time"
)

func solveA(in io.Reader) string {
    reader := bufio.NewReader(in)
    var h1, h2 int
    if _, err := fmt.Fscan(reader, &h1, &h2); err != nil {
        return ""
    }
    var a, b int
    fmt.Fscan(reader, &a, &b)
    if h1+8*a >= h2 {
        return "0"
    }
    if a <= b {
        return "-1"
    }
    h := h1
    hours := 0
    for {
        hours += 8
        h += 8 * a
        if h >= h2 {
            return fmt.Sprintf("%d", hours/24)
        }
        hours += 12
        h -= 12 * b
        hours += 4
        h += 4 * a
        if h >= h2 {
            return fmt.Sprintf("%d", hours/24)
        }
    }
}

func genTests() []string {
    r := rand.New(rand.NewSource(1))
    tests := make([]string, 100)
    for i := 0; i < 100; i++ {
        h1 := r.Intn(99999) + 1
        h2 := h1 + r.Intn(100000-h1) + 1
        a := r.Intn(100000) + 1
        b := r.Intn(100000) + 1
        tests[i] = fmt.Sprintf("%d %d\n%d %d\n", h1, h2, a, b)
    }
    return tests
}

func runBinary(bin string, input string) (string, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()
    cmd := exec.CommandContext(ctx, bin)
    cmd.Stdin = strings.NewReader(input)
    var out bytes.Buffer
    cmd.Stdout = &out
    err := cmd.Run()
    return strings.TrimSpace(out.String()), err
}

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: go run verifierA.go /path/to/binary")
        return
    }
    bin := os.Args[1]
    tests := genTests()
    for i, tc := range tests {
        expected := solveA(strings.NewReader(tc))
        actual, err := runBinary(bin, tc)
        if err != nil {
            fmt.Printf("Test %d: runtime error: %v\n", i+1, err)
            return
        }
        if actual != strings.TrimSpace(expected) {
            fmt.Printf("Test %d failed.\nInput:\n%sExpected: %s\nGot: %s\n", i+1, tc, expected, actual)
            return
        }
    }
    fmt.Println("All tests passed!")
}


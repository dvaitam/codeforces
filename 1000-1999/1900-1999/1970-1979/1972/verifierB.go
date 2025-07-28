package main

import (
    "bytes"
    "fmt"
    "io"
    "math/rand"
    "os"
    "os/exec"
    "strings"
    "time"
)

// expectedResult returns the expected answer for problem B
func expectedResult(s string) string {
    cnt := 0
    for _, ch := range s {
        if ch == 'U' {
            cnt++
        }
    }
    if cnt%2 == 1 {
        return "YES"
    }
    return "NO"
}

func runTest(binary string, n int, s string) error {
    var input strings.Builder
    input.WriteString("1\n")
    input.WriteString(fmt.Sprintf("%d\n%s\n", n, s))

    cmd := exec.Command(binary)
    cmd.Stdin = strings.NewReader(input.String())
    var out bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = io.Discard

    if err := cmd.Run(); err != nil {
        return fmt.Errorf("execution error: %v", err)
    }
    result := strings.TrimSpace(out.String())
    expected := expectedResult(s)
    if !strings.EqualFold(result, expected) {
        return fmt.Errorf("expected %s, got %s", expected, result)
    }
    return nil
}

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: go run verifierB.go /path/to/binary")
        os.Exit(1)
    }
    binary := os.Args[1]
    rand.Seed(time.Now().UnixNano())
    const tests = 100
    for t := 0; t < tests; t++ {
        n := rand.Intn(100) + 1
        var sb strings.Builder
        for i := 0; i < n; i++ {
            if rand.Intn(2) == 0 {
                sb.WriteByte('U')
            } else {
                sb.WriteByte('D')
            }
        }
        s := sb.String()
        if err := runTest(binary, n, s); err != nil {
            fmt.Printf("Test %d failed: %v\n", t+1, err)
            os.Exit(1)
        }
    }
    fmt.Printf("All %d tests passed\n", tests)
}


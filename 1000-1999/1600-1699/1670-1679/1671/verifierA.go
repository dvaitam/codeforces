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

func check(s string) string {
    for i := 0; i < len(s); {
        j := i
        for j < len(s) && s[j] == s[i] {
            j++
        }
        if j-i == 1 {
            return "NO"
        }
        i = j
    }
    return "YES"
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
        fmt.Println("usage: go run verifierA.go /path/to/binary")
        return
    }
    bin := os.Args[1]
    r := rand.New(rand.NewSource(1))
    tests := 100
    for i := 0; i < tests; i++ {
        l := r.Intn(50) + 1
        sb := make([]byte, l)
        for j := 0; j < l; j++ {
            if r.Intn(2) == 0 {
                sb[j] = 'a'
            } else {
                sb[j] = 'b'
            }
        }
        s := string(sb)
        input := fmt.Sprintf("1\n%s\n", s)
        expected := check(s)
        out, err := runBinary(bin, input)
        if err != nil {
            fmt.Printf("test %d runtime error: %v\n", i+1, err)
            return
        }
        out = strings.ToUpper(strings.TrimSpace(out))
        if out != expected {
            fmt.Printf("test %d failed: input=%s expected=%s got=%s\n", i+1, s, expected, out)
            return
        }
    }
    fmt.Println("All tests passed")
}

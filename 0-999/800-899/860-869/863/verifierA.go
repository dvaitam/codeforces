package main

import (
    "bytes"
    "fmt"
    "math/rand"
    "os"
    "os/exec"
    "strings"
)

func expectedA(x string) string {
    for len(x) > 0 && x[len(x)-1] == '0' {
        x = x[:len(x)-1]
    }
    isPal := true
    for i := 0; i < len(x)/2; i++ {
        if x[i] != x[len(x)-1-i] {
            isPal = false
            break
        }
    }
    if isPal {
        return "YES"
    }
    return "NO"
}

func genTestsA() []string {
    rand.Seed(1)
    tests := make([]string, 0, 100)
    specials := []string{"1", "10", "100", "11", "101", "1001", "1234321", "200000000", "1000000000"}
    tests = append(tests, specials...)
    for len(tests) < 100 {
        x := rand.Intn(1_000_000_000) + 1
        tests = append(tests, fmt.Sprintf("%d", x))
    }
    return tests
}

func runBinary(bin string, input string) (string, error) {
    cmd := exec.Command(bin)
    cmd.Stdin = strings.NewReader(input)
    var out bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = os.Stderr
    err := cmd.Run()
    return strings.TrimSpace(out.String()), err
}

func main() {
    if len(os.Args) != 2 {
        fmt.Fprintf(os.Stderr, "Usage: go run verifierA.go <binary>\n")
        os.Exit(1)
    }
    bin := os.Args[1]
    tests := genTestsA()
    for i, t := range tests {
        input := t + "\n"
        want := expectedA(t)
        got, err := runBinary(bin, input)
        if err != nil {
            fmt.Printf("Test %d: runtime error: %v\n", i+1, err)
            os.Exit(1)
        }
        if strings.TrimSpace(got) != want {
            fmt.Printf("Test %d failed. Input: %s Expected: %s Got: %s\n", i+1, t, want, got)
            os.Exit(1)
        }
    }
    fmt.Printf("All %d tests passed.\n", len(tests))
}


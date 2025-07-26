package main

import (
    "bytes"
    "fmt"
    "math"
    "math/rand"
    "os"
    "os/exec"
    "strings"
)

type TestCase struct {
    input string
    ans   string
}

func rotate(k, xi, yi int) string {
    x := float64(xi)
    y := float64(yi)
    theta := float64(k) * math.Pi / 180.0
    c := math.Cos(theta)
    s := math.Sin(theta)
    x2 := x*c - y*s
    y2 := x*s + y*c
    return fmt.Sprintf("%.10f %.10f", x2, y2)
}

func genTests() []TestCase {
    r := rand.New(rand.NewSource(9))
    cases := make([]TestCase, 100)
    for i := 0; i < 100; i++ {
        k := r.Intn(360)
        xi := r.Intn(201) - 100
        yi := r.Intn(201) - 100
        input := fmt.Sprintf("%d %d %d\n", k, xi, yi)
        ans := rotate(k, xi, yi)
        cases[i] = TestCase{input: input, ans: ans}
    }
    return cases
}

func run(bin, in string) (string, error) {
    cmd := exec.Command(bin)
    cmd.Stdin = strings.NewReader(in)
    var out bytes.Buffer
    cmd.Stdout = &out
    err := cmd.Run()
    return strings.TrimSpace(out.String()), err
}

func main() {
    if len(os.Args) < 2 {
        fmt.Fprintln(os.Stderr, "usage: go run verifierI.go /path/to/binary")
        os.Exit(1)
    }
    bin := os.Args[1]
    tests := genTests()
    for i, tc := range tests {
        got, err := run(bin, tc.input)
        if err != nil {
            fmt.Fprintf(os.Stderr, "test %d: runtime error: %v\n", i+1, err)
            os.Exit(1)
        }
        if strings.TrimSpace(got) != tc.ans {
            fmt.Fprintf(os.Stderr, "test %d failed: input %q expected %q got %q\n", i+1, tc.input, tc.ans, got)
            os.Exit(1)
        }
    }
    fmt.Println("All tests passed")
}


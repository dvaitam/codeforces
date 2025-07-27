package main

import (
    "bytes"
    "fmt"
    "math/rand"
    "os"
    "os/exec"
    "strings"
    "time"
)

type testCase struct{a,b int}

func (tc testCase) input() string {
    return fmt.Sprintf("%d %d\n", tc.a, tc.b)
}

func solve(a,b int) string {
    return fmt.Sprintf("%d", a+b)
}

func randomCase(rng *rand.Rand) testCase {
    a := rng.Intn(2000001) - 1000000
    b := rng.Intn(2000001) - 1000000
    return testCase{a:a,b:b}
}

func deterministicCases() []testCase {
    return []testCase{{0,0},{1,2},{-5,7}}
}

func runCase(bin string, tc testCase) error {
    input := tc.input()
    cmd := exec.Command(bin)
    cmd.Stdin = strings.NewReader(input)
    var out bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = &out
    if err := cmd.Run(); err != nil {
        return fmt.Errorf("runtime error: %v\n%s", err, out.String())
    }
    result := strings.TrimSpace(out.String())
    expect := solve(tc.a, tc.b)
    if result != expect {
        return fmt.Errorf("expected %s got %s", expect, result)
    }
    return nil
}

func main() {
    if len(os.Args) != 2 {
        fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
        os.Exit(1)
    }
    bin := os.Args[1]
    rng := rand.New(rand.NewSource(time.Now().UnixNano()))
    tests := deterministicCases()
    for len(tests) < 100 {
        tests = append(tests, randomCase(rng))
    }
    for i, tc := range tests {
        if err := runCase(bin, tc); err != nil {
            fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input())
            os.Exit(1)
        }
    }
    fmt.Println("All tests passed")
}


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

type testCase struct {
    n, m, a, b int64
}

func solve(n, m, a, b int64) int64 {
    ra := (a - 1) / m
    rb := (b - 1) / m
    ca := (a-1)%m + 1
    cb := (b-1)%m + 1
    if ra == rb {
        return 1
    } else if ca == 1 && cb == m {
        return 1
    } else if ra+1 == rb {
        return 2
    }
    return 3
}

func runCase(bin string, tc testCase) error {
    input := fmt.Sprintf("%d %d %d %d\n", tc.n, tc.m, tc.a, tc.b)
    cmd := exec.Command(bin)
    cmd.Stdin = strings.NewReader(input)
    var out bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = &out
    if err := cmd.Run(); err != nil {
        return fmt.Errorf("runtime error: %v\n%s", err, out.String())
    }
    var ans int64
    if _, err := fmt.Fscan(strings.NewReader(strings.TrimSpace(out.String())), &ans); err != nil {
        return fmt.Errorf("invalid output: %v", err)
    }
    expected := solve(tc.n, tc.m, tc.a, tc.b)
    if ans != expected {
        return fmt.Errorf("expected %d got %d", expected, ans)
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

    var cases []testCase
    cases = append(cases, testCase{n: 1, m: 1, a: 1, b: 1})
    cases = append(cases, testCase{n: 10, m: 5, a: 3, b: 9})
    cases = append(cases, testCase{n: 10, m: 5, a: 1, b: 10})
    for i := 0; i < 100; i++ {
        n := rng.Int63n(1_000_000_000) + 1
        m := rng.Int63n(1_000_000_000) + 1
        if n < 1 {
            n = 1
        }
        a := rng.Int63n(n) + 1
        b := a + rng.Int63n(n-a+1)
        cases = append(cases, testCase{n: n, m: m, a: a, b: b})
    }

    for i, tc := range cases {
        if err := runCase(bin, tc); err != nil {
            fmt.Fprintf(os.Stderr, "test %d failed: %v\ninput: %d %d %d %d\n", i+1, err, tc.n, tc.m, tc.a, tc.b)
            os.Exit(1)
        }
    }
    fmt.Println("All tests passed")
}


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

func solve(n, W, m int) string {
    if n < m && m%(m-n) > 0 {
        return "NO\n"
    }
    var sb strings.Builder
    sb.WriteString("YES\n")
    cur, used := 1, 0
    for i := 0; i < m; i++ {
        sum := 0
        printed := false
        for sum < n {
            cnt := n - sum
            if m-used < cnt {
                cnt = m - used
            }
            if printed {
                sb.WriteByte(' ')
            }
            length := float64(cnt) / float64(m) * float64(W)
            sb.WriteString(fmt.Sprintf("%d %.16f", cur, length))
            printed = true
            used += cnt
            sum += cnt
            if used == m {
                cur++
                used = 0
            }
        }
        sb.WriteByte('\n')
    }
    return sb.String()
}

type testCase struct {
    n, W, m int
}

func runCase(bin string, tc testCase) error {
    input := fmt.Sprintf("%d %d %d\n", tc.n, tc.W, tc.m)
    cmd := exec.Command(bin)
    cmd.Stdin = strings.NewReader(input)
    var out bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = &out
    if err := cmd.Run(); err != nil {
        return fmt.Errorf("runtime error: %v\n%s", err, out.String())
    }
    expected := strings.TrimSpace(solve(tc.n, tc.W, tc.m))
    got := strings.TrimSpace(out.String())
    if expected != got {
        return fmt.Errorf("expected:\n%s\n-- got:\n%s", expected, got)
    }
    return nil
}

func main() {
    if len(os.Args) != 2 {
        fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
        os.Exit(1)
    }
    bin := os.Args[1]
    rng := rand.New(rand.NewSource(time.Now().UnixNano()))

    var cases []testCase
    cases = append(cases, testCase{n: 1, W: 100, m: 2})
    cases = append(cases, testCase{n: 3, W: 10, m: 2})
    cases = append(cases, testCase{n: 5, W: 123, m: 5})
    for i := 0; i < 100; i++ {
        n := rng.Intn(50) + 1
        W := rng.Intn(901) + 100
        m := rng.Intn(49) + 2
        cases = append(cases, testCase{n: n, W: W, m: m})
    }

    for i, tc := range cases {
        if err := runCase(bin, tc); err != nil {
            fmt.Fprintf(os.Stderr, "test %d failed: %v\ninput: %d %d %d\n", i+1, err, tc.n, tc.W, tc.m)
            os.Exit(1)
        }
    }
    fmt.Println("All tests passed")
}


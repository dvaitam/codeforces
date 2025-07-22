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
    input  string
    expect string
}

func solve(n int, s string) string {
    cR, cL := 0, 0
    for i := 0; i < n; i++ {
        if s[i] == 'R' {
            cR++
        } else if s[i] == 'L' {
            cL++
        }
    }
    var res strings.Builder
    if cR == 0 {
        for i := n - 1; i >= 0; i-- {
            if s[i] == 'L' {
                fmt.Fprintf(&res, "%d ", i+1)
                break
            }
        }
        for i := 0; i < n; i++ {
            if s[i] == 'L' {
                fmt.Fprintf(&res, "%d", i)
                break
            }
        }
        res.WriteByte('\n')
        return res.String()
    }
    if cL == 0 {
        for i := 0; i < n; i++ {
            if s[i] == 'R' {
                fmt.Fprintf(&res, "%d ", i+1)
                break
            }
        }
        for i := n - 1; i >= 0; i-- {
            if s[i] == 'R' {
                fmt.Fprintf(&res, "%d", i+2)
                break
            }
        }
        res.WriteByte('\n')
        return res.String()
    }
    for i := 0; i < n; i++ {
        if s[i] == 'R' {
            fmt.Fprintf(&res, "%d ", i+1)
            break
        }
    }
    for i := 0; i < n-1; i++ {
        if s[i] == 'R' && s[i+1] == 'L' {
            fmt.Fprintf(&res, "%d", i+1)
            break
        }
    }
    res.WriteByte('\n')
    return res.String()
}

func run(bin, input string) (string, error) {
    var cmd *exec.Cmd
    if strings.HasSuffix(bin, ".go") {
        cmd = exec.Command("go", "run", bin)
    } else {
        cmd = exec.Command(bin)
    }
    cmd.Stdin = strings.NewReader(input)
    var out bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = &out
    if err := cmd.Run(); err != nil {
        return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
    }
    return out.String(), nil
}

func generateCase(rng *rand.Rand) (string, string) {
    n := rng.Intn(18) + 3 // [3,20]
    b := make([]byte, n)
    b[0] = '.'
    b[n-1] = '.'
    typ := rng.Intn(3)
    if typ == 0 {
        // only L
        for i := 1; i < n-1; i++ {
            if rng.Intn(2) == 0 {
                b[i] = 'L'
            } else {
                b[i] = '.'
            }
        }
        if !strings.ContainsRune(string(b), 'L') {
            b[1] = 'L'
        }
    } else if typ == 1 {
        // only R
        for i := 1; i < n-1; i++ {
            if rng.Intn(2) == 0 {
                b[i] = 'R'
            } else {
                b[i] = '.'
            }
        }
        if !strings.ContainsRune(string(b), 'R') {
            b[1] = 'R'
        }
    } else {
        for i := 1; i < n-1; i++ {
            b[i] = '.'
        }
        idx := rng.Intn(n-2) + 1
        b[idx] = 'R'
        b[idx+1] = 'L'
        for i := 1; i < n-1; i++ {
            if i == idx || i == idx+1 {
                continue
            }
            r := rng.Intn(3)
            if r == 0 {
                b[i] = 'R'
            } else if r == 1 {
                b[i] = 'L'
            } else {
                b[i] = '.'
            }
        }
    }
    s := string(b)
    input := fmt.Sprintf("%d\n%s\n", n, s)
    expect := solve(n, s)
    return input, expect
}

func main() {
    if len(os.Args) != 2 {
        fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
        os.Exit(1)
    }
    bin := os.Args[1]
    rng := rand.New(rand.NewSource(time.Now().UnixNano()))
    for i := 0; i < 100; i++ {
        in, exp := generateCase(rng)
        out, err := run(bin, in)
        if err != nil {
            fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
            os.Exit(1)
        }
        if strings.TrimSpace(out) != strings.TrimSpace(exp) {
            fmt.Fprintf(os.Stderr, "case %d failed: expected %q got %q\ninput:\n%s", i+1, strings.TrimSpace(exp), strings.TrimSpace(out), in)
            os.Exit(1)
        }
    }
    fmt.Println("All tests passed")
}


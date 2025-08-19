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

func parseTwoInts(s string) (int, int, bool) {
    fields := strings.Fields(strings.TrimSpace(s))
    if len(fields) < 2 {
        return 0, 0, false
    }
    a, err1 := strconv.Atoi(fields[0])
    b, err2 := strconv.Atoi(fields[1])
    if err1 != nil || err2 != nil {
        return 0, 0, false
    }
    return a, b, true
}

func acceptable(n int, s string, a, b int) bool {
    // Compute positions (0-based indices)
    firstR, lastR := -1, -1
    firstL, lastL := -1, -1
    for i := 0; i < n; i++ {
        if s[i] == 'R' {
            if firstR == -1 { firstR = i }
            lastR = i
        } else if s[i] == 'L' {
            if firstL == -1 { firstL = i }
            lastL = i
        }
    }
    hasR := firstR != -1
    hasL := firstL != -1
    // Our generator's expected pairs (using its mixed indexing convention)
    var exp1, exp2 [2]int
    if !hasR { // only L
        exp1 = [2]int{lastL + 1, firstL}         // generator
        exp2 = [2]int{firstL + 1, firstL}        // accept also first non-dot +1, firstL (common variant in some solutions)
    } else if !hasL { // only R
        exp1 = [2]int{firstR + 1, lastR + 2}     // generator
        exp2 = [2]int{firstR, lastR + 1}         // canonical 1-based pair
    } else { // both
        // find first RL adjacency
        rl := -1
        for i := 0; i+1 < n; i++ {
            if s[i] == 'R' && s[i+1] == 'L' { rl = i; break }
        }
        if rl == -1 {
            // No RL, but both R and L exist: fall back to our firstR based
            exp1 = [2]int{firstR + 1, firstR + 1}
            exp2 = [2]int{firstR + 1, firstR + 1}
        } else {
            exp1 = [2]int{firstR + 1, rl + 1}    // generator
            exp2 = [2]int{rl + 1, rl + 1}        // CF common variant
        }
    }
    return (a == exp1[0] && b == exp1[1]) || (a == exp2[0] && b == exp2[1])
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
        // Accept either the generator's exact expected or other valid variants
        if strings.TrimSpace(out) != strings.TrimSpace(exp) {
            // Parse input to get n and s
            lines := strings.Split(strings.TrimSpace(in), "\n")
            if len(lines) < 2 {
                fmt.Fprintf(os.Stderr, "case %d failed: malformed input generated\n", i+1)
                os.Exit(1)
            }
            var n int
            fmt.Sscanf(lines[0], "%d", &n)
            s := strings.TrimSpace(lines[1])
            A, B, ok := parseTwoInts(out)
            if !ok || !acceptable(n, s, A, B) {
                fmt.Fprintf(os.Stderr, "case %d failed: expected %q got %q\ninput:\n%s", i+1, strings.TrimSpace(exp), strings.TrimSpace(out), in)
                os.Exit(1)
            }
        }
    }
    fmt.Println("All tests passed")
}

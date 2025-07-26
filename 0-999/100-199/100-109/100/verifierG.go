package main

import (
    "bytes"
    "fmt"
    "math/rand"
    "os"
    "os/exec"
    "strings"
)

type TestCase struct {
    input string
    ans   string
}

type Record struct {
    name string
    year int
}

func solve(records []Record, cands []string) string {
    used := make(map[string]int)
    for _, r := range records {
        if prev, ok := used[r.name]; !ok || r.year > prev {
            used[r.name] = r.year
        }
    }
    var unusedBest string
    haveUnused := false
    var usedBest string
    var usedYearBest int
    haveUsed := false
    for _, cand := range cands {
        if year, ok := used[cand]; !ok {
            if !haveUnused || cand > unusedBest {
                unusedBest = cand
                haveUnused = true
            }
        } else {
            if !haveUsed || year < usedYearBest || (year == usedYearBest && cand > usedBest) {
                usedYearBest = year
                usedBest = cand
                haveUsed = true
            }
        }
    }
    if haveUnused {
        return unusedBest
    }
    return usedBest
}

func randomName(r *rand.Rand) string {
    b := []byte{'a' + byte(r.Intn(26)), 'a' + byte(r.Intn(26))}
    return string(b)
}

func genTests() []TestCase {
    r := rand.New(rand.NewSource(7))
    cases := make([]TestCase, 100)
    for i := 0; i < 100; i++ {
        n := r.Intn(5) + 1
        rec := make([]Record, n)
        for j := 0; j < n; j++ {
            rec[j] = Record{randomName(r), 1990 + r.Intn(35)}
        }
        m := r.Intn(5) + 1
        cand := make([]string, m)
        for j := 0; j < m; j++ {
            cand[j] = randomName(r)
        }
        ans := solve(rec, cand)
        var sb strings.Builder
        sb.WriteString(fmt.Sprintf("%d\n", n))
        for j, rc := range rec {
            sb.WriteString(fmt.Sprintf("%s %d\n", rc.name, rc.year))
            if j == n-1 {
                // nothing
            }
        }
        sb.WriteString(fmt.Sprintf("%d\n", m))
        for j, c := range cand {
            if j > 0 {
                sb.WriteByte(' ')
            }
            sb.WriteString(c)
        }
        sb.WriteString("\n")
        cases[i] = TestCase{input: sb.String(), ans: ans}
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
        fmt.Fprintln(os.Stderr, "usage: go run verifierG.go /path/to/binary")
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


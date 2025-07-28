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

type Var struct {
    pref string
    suff string
    cnt  int
}

type Statement struct {
    line string
}

func countHaha(s string) int {
    c := 0
    for i := 0; i+3 < len(s); i++ {
        if s[i:i+4] == "haha" {
            c++
        }
    }
    return c
}

func merge(a, b Var) Var {
    res := Var{}
    res.cnt = a.cnt + b.cnt + countHaha(a.suff+b.pref)
    pref := a.pref + b.pref
    if len(pref) > 3 {
        pref = pref[:3]
    }
    res.pref = pref
    suff := a.suff + b.suff
    if len(suff) > 3 {
        suff = suff[len(suff)-3:]
    }
    res.suff = suff
    return res
}

func solve(statements []string) int {
    vars := make(map[string]Var)
    last := ""
    for _, line := range statements {
        parts := strings.Fields(line)
        x := parts[0]
        if parts[1] == ":=" {
            s := parts[2]
            v := Var{pref: s, suff: s, cnt: countHaha(s)}
            if len(v.pref) > 3 {
                v.pref = v.pref[:3]
            }
            if len(v.suff) > 3 {
                v.suff = v.suff[len(v.suff)-3:]
            }
            vars[x] = v
        } else { // =
            a := parts[2]
            b := parts[4]
            v := merge(vars[a], vars[b])
            vars[x] = v
        }
        last = x
    }
    return vars[last].cnt
}

type TestCase struct {
    lines []string
}

func genString() string {
    n := rand.Intn(5) + 1
    bytes := make([]byte, n)
    for i := range bytes {
        bytes[i] = byte('a' + rand.Intn(26))
    }
    return string(bytes)
}

func generateTests() []TestCase {
    rand.Seed(46)
    tests := make([]TestCase, 100)
    for i := range tests {
        n := rand.Intn(5) + 1
        lines := make([]string, n)
        names := []string{}
        for j := 0; j < n; j++ {
            name := genString()
            names = append(names, name)
            if rand.Intn(2) == 0 || j == 0 {
                s := genString()
                lines[j] = fmt.Sprintf("%s := %s", name, s)
            } else {
                a := names[rand.Intn(len(names))]
                b := names[rand.Intn(len(names))]
                lines[j] = fmt.Sprintf("%s = %s + %s", name, a, b)
            }
        }
        tests[i] = TestCase{lines: lines}
    }
    return tests
}

func runBinary(bin, input string) (string, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()
    cmd := exec.CommandContext(ctx, bin)
    cmd.Stdin = strings.NewReader(input)
    var out bytes.Buffer
    cmd.Stdout = &out
    if err := cmd.Run(); err != nil {
        return "", err
    }
    return strings.TrimSpace(out.String()), nil
}

func main() {
    if len(os.Args) < 2 {
        fmt.Println("usage: go run verifierE.go /path/to/binary")
        os.Exit(1)
    }
    bin := os.Args[1]
    tests := generateTests()
    for idx, tc := range tests {
        input := fmt.Sprintf("1\n%d\n", len(tc.lines))
        for _, ln := range tc.lines {
            input += ln + "\n"
        }
        want := fmt.Sprintf("%d", solve(tc.lines))
        got, err := runBinary(bin, input)
        if err != nil {
            fmt.Printf("test %d: runtime error: %v\n", idx+1, err)
            os.Exit(1)
        }
        if got != want {
            fmt.Printf("test %d failed\ninput:\n%sexpected: %s\ngot: %s\n", idx+1, input, want, got)
            os.Exit(1)
        }
    }
    fmt.Println("all tests passed")
}


package main

import (
    "bytes"
    "fmt"
    "math/rand"
    "os"
    "os/exec"
    "strconv"
    "strings"
)

type Test string

func runExe(path, input string) (string, error) {
    cmd := exec.Command(path)
    cmd.Stdin = strings.NewReader(input)
    var out bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = &out
    err := cmd.Run()
    return out.String(), err
}

// Generate random tests of segments [l, r] with l <= r.
func genTests() []Test {
    rand.Seed(2)
    tests := []Test{Test("1\n1 1\n")}
    for i := 0; i < 100; i++ {
        n := rand.Intn(20) + 1
        var sb strings.Builder
        sb.WriteString(fmt.Sprintf("%d\n", n))
        for j := 0; j < n; j++ {
            l := rand.Intn(100) + 1
            r := l + rand.Intn(100)
            sb.WriteString(fmt.Sprintf("%d %d\n", l, r))
        }
        tests = append(tests, Test(sb.String()))
    }
    return tests
}

type seg struct{ l, r int }

func anyNestedPair(segs []seg) (int, int, bool) {
    n := len(segs)
    for i := 0; i < n; i++ {
        for j := 0; j < n; j++ {
            if i == j { continue }
            if segs[i].l >= segs[j].l && segs[i].r <= segs[j].r {
                return i + 1, j + 1, true
            }
        }
    }
    return -1, -1, false
}

func parseOutInts(out string) (int, int, error) {
    // parse first two integers from the output (ignore extra text/whitespace)
    fs := strings.Fields(out)
    if len(fs) < 2 {
        return 0, 0, fmt.Errorf("output does not contain two integers: %q", out)
    }
    a, err := strconv.Atoi(fs[0])
    if err != nil { return 0, 0, fmt.Errorf("failed to parse first int: %v", err) }
    b, err := strconv.Atoi(fs[1])
    if err != nil { return 0, 0, fmt.Errorf("failed to parse second int: %v", err) }
    return a, b, nil
}

func main() {
    if len(os.Args) != 2 {
        fmt.Println("Usage: go run verifierC.go /path/to/binary")
        return
    }
    bin := os.Args[1]

    tests := genTests()
    for i, t := range tests {
        input := string(t)
        // parse input segments
        lines := strings.Split(strings.TrimSpace(input), "\n")
        if len(lines) == 0 { fmt.Println("bad test input"); os.Exit(1) }
        var n int
        fmt.Sscanf(lines[0], "%d", &n)
        segs := make([]seg, n)
        for j := 0; j < n; j++ {
            var l, r int
            fmt.Sscanf(lines[j+1], "%d %d", &l, &r)
            segs[j] = seg{l, r}
        }

        gotStr, err := runExe(bin, input)
        if err != nil {
            fmt.Printf("candidate runtime error on test %d: %v\n", i+1, err)
            os.Exit(1)
        }
        a, b, err := parseOutInts(gotStr)
        if err != nil {
            fmt.Printf("Test %d failed (parse error)\nInput:%sError:%v\nOutput:%s\n", i+1, input, err, gotStr)
            os.Exit(1)
        }
        if a == -1 && b == -1 {
            // verify that no valid pair exists
            if _, _, ok := anyNestedPair(segs); ok {
                fmt.Printf("Test %d failed\nInput:%sExpected:any valid pair (i j), but candidate printed -1 -1\nGot:%s\n", i+1, input, gotStr)
                os.Exit(1)
            }
            continue
        }
        // verify indices and containment
        if a < 1 || a > n || b < 1 || b > n || a == b {
            fmt.Printf("Test %d failed\nInput:%sExpected:two distinct indices in [1..n]\nGot:%d %d\n", i+1, input, a, b)
            os.Exit(1)
        }
        ia, ib := a-1, b-1
        if !(segs[ia].l >= segs[ib].l && segs[ia].r <= segs[ib].r) {
            fmt.Printf("Test %d failed\nInput:%sExpected:(%d,%d) to satisfy containment, i.e., l[i]>=l[j] and r[i]<=r[j]\nGot:%d %d\n", i+1, input, a, b, a, b)
            os.Exit(1)
        }
    }
    fmt.Println("all tests passed")
}

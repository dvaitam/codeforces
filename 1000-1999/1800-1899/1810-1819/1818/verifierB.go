package main

import (
    "bufio"
    "bytes"
    "encoding/base64"
    "fmt"
    "os"
    "os/exec"
    "strconv"
    "strings"
)

// Base64-encoded contents of testcasesB.txt.
const testcasesB64 = "MzMKOTMKMzcKNDkKNzgKNjYKMwo2NAo4Mwo2Ngo2NQoyNgo1NQozNQo5NAo4NAo5OAo1Mwo5Mwo1CjMKNTQKNzIKMTEKMzYKOTIKNAo0CjY0Cjg3CjUwCjY5CjQ1CjcxCjYxCjY4CjU3CjcyCjk2Cjk0CjkxCjYxCjc1CjI0Cjc2Cjk3CjQKNzgKNDYKMjcKNDAKNTYKOQoyOQoxMwoyMgo3OAo1OQo4MQoxNwo2NAo0MQo1Nwo4OQoyMAo2MwozMwo2Mwo0NAo5MQo2OAo5MwoyOQoxMDAKNgo0Nwo2NgozNAo1MQo2OQo4NAo3NwoxMwo0MAoxNwo2NQoxOAoxOQo5CjQ0CjI1CjYKNAo2Mwo2MQoxMAo3NAo2OAo3CjE4Cg=="

type testCase struct {
    n int
}

// Embedded solution logic from 1818B.go.
func solve(tc testCase) string {
    n := tc.n
    if n%2 == 1 {
        if n == 1 {
            return "1"
        }
        return "-1"
    }
    var sb strings.Builder
    for i := 1; i <= n; i += 2 {
        if i > 1 {
            sb.WriteByte(' ')
        }
        fmt.Fprintf(&sb, "%d %d", i+1, i)
    }
    return sb.String()
}

func parseTestcases() ([]testCase, error) {
    raw, err := base64.StdEncoding.DecodeString(testcasesB64)
    if err != nil {
        return nil, err
    }
    sc := bufio.NewScanner(bytes.NewReader(raw))
    cases := make([]testCase, 0)
    for sc.Scan() {
        line := strings.TrimSpace(sc.Text())
        if line == "" {
            continue
        }
        val, err := strconv.Atoi(line)
        if err != nil {
            return nil, fmt.Errorf("parse n: %v", err)
        }
        cases = append(cases, testCase{n: val})
    }
    if err := sc.Err(); err != nil {
        return nil, err
    }
    if len(cases) == 0 {
        return nil, fmt.Errorf("no testcases found")
    }
    return cases, nil
}

func buildIfGo(path string) (string, func(), error) {
    if strings.HasSuffix(path, ".go") {
        tmp, err := os.CreateTemp("", "bin*")
        if err != nil {
            return "", nil, err
        }
        tmp.Close()
        if out, err := exec.Command("go", "build", "-o", tmp.Name(), path).CombinedOutput(); err != nil {
            os.Remove(tmp.Name())
            return "", nil, fmt.Errorf("build failed: %v\n%s", err, out)
        }
        return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
    }
    return path, func() {}, nil
}

func runCandidate(bin, input string) (string, error) {
    cmd := exec.Command(bin)
    cmd.Stdin = strings.NewReader(input)
    var out bytes.Buffer
    var errBuf bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = &errBuf
    if err := cmd.Run(); err != nil {
        return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
    }
    return strings.TrimSpace(out.String()), nil
}

func main() {
    if len(os.Args) != 2 {
        fmt.Println("usage: go run verifierB.go /path/to/binary")
        os.Exit(1)
    }

    cases, err := parseTestcases()
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }

    bin, cleanup, err := buildIfGo(os.Args[1])
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
    defer cleanup()

    for idx, tc := range cases {
        input := fmt.Sprintf("1\n%d\n", tc.n)
        want := solve(tc)
        got, err := runCandidate(bin, input)
        if err != nil {
            fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
            os.Exit(1)
        }
        if got != want {
            fmt.Printf("test %d failed\ninput:\n%sexpected: %s\ngot: %s\n", idx+1, input, want, got)
            os.Exit(1)
        }
    }
    fmt.Printf("All %d tests passed\n", len(cases))
}


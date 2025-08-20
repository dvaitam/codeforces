package main

import (
    "bufio"
    "bytes"
    "fmt"
    "os"
    "os/exec"
    "strconv"
    "strings"
)

func run(bin, input string) (string, error) {
    var cmd *exec.Cmd
    if strings.HasSuffix(bin, ".go") {
        cmd = exec.Command("go", "run", bin)
    } else {
        cmd = exec.Command(bin)
    }
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

func applyRotation(s string, b []int) string {
    if len(b) == 0 { return s }
    bytesS := []byte(s)
    last := bytesS[b[len(b)-1]-1]
    for i := len(b)-1; i >= 1; i-- {
        bytesS[b[i]-1] = bytesS[b[i-1]-1]
    }
    bytesS[b[0]-1] = last
    return string(bytesS)
}

func validate(n int, s string, out string) (bool, string) {
    // Tokenize output
    tokens := strings.Fields(out)
    if len(tokens) == 0 { return false, "empty output" }
    if tokens[0] == "-1" && len(tokens) == 1 {
        // only acceptable when count of '1' is odd (impossible case)
        ones := 0
        for i := 0; i < len(s); i++ { if s[i] == '1' { ones++ } }
        if ones%2 == 1 { return true, "" }
        return false, "reported -1 but solution should exist (even ones)"
    }
    // Parse m and b list
    m, err := strconv.Atoi(tokens[0])
    if err != nil || m < 0 || m > 2*n { return false, "invalid m" }
    if len(tokens) < 1+m+n { return false, "not enough tokens for b and p" }
    b := make([]int, m)
    for i := 0; i < m; i++ {
        v, err := strconv.Atoi(tokens[1+i])
        if err != nil || v < 1 || v > 2*n { return false, "invalid b index" }
        b[i] = v
        if i > 0 && b[i] <= b[i-1] { return false, "b not strictly increasing" }
    }
    p := make([]int, n)
    off := 1 + m
    used := make([]bool, 2*n+1)
    for i := 0; i < n; i++ {
        v, err := strconv.Atoi(tokens[off+i])
        if err != nil || v < 1 || v > 2*n { return false, "invalid p index" }
        if used[v] { return false, "duplicate index in p" }
        if i > 0 && v < p[i-1] { return false, "p not nondecreasing" }
        used[v] = true
        p[i] = v
    }
    // Build q as complement
    q := make([]int, 0, n)
    for i := 1; i <= 2*n; i++ { if !used[i] { q = append(q, i) } }
    if len(q) != n { return false, "p and q not partitioning" }
    // Apply rotation
    s2 := applyRotation(s, b)
    // Build strings
    var sp, sq strings.Builder
    for _, idx := range p { sp.WriteByte(s2[idx-1]) }
    for _, idx := range q { sq.WriteByte(s2[idx-1]) }
    if sp.String() != sq.String() { return false, "subsequences not equal" }
    return true, ""
}

func main() {
    if len(os.Args) != 2 {
        fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
        os.Exit(1)
    }
    bin := os.Args[1]

    file, err := os.Open("testcasesD.txt")
    if err != nil {
        fmt.Fprintf(os.Stderr, "failed to open testcases: %v\n", err)
        os.Exit(1)
    }
    defer file.Close()
    scanner := bufio.NewScanner(file)
    idx := 0
    for scanner.Scan() {
        line := strings.TrimSpace(scanner.Text())
        if line == "" { continue }
        idx++
        fields := strings.Fields(line)
        if len(fields) != 2 {
            fmt.Fprintf(os.Stderr, "bad testcase on line %d\n", idx)
            os.Exit(1)
        }
        n, _ := strconv.Atoi(fields[0])
        s := fields[1]
        input := fmt.Sprintf("1\n%d\n%s\n", n, s)
        got, err := run(bin, input)
        if err != nil {
            fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx, err)
            os.Exit(1)
        }
        ok, reason := validate(n, s, got)
        if !ok {
            fmt.Printf("case %d failed\nexpected: valid construction or -1 if odd ones\n got: %s\nreason: %s\n", idx, got, reason)
            os.Exit(1)
        }
    }
    if err := scanner.Err(); err != nil {
        fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
        os.Exit(1)
    }
    fmt.Printf("All %d tests passed\n", idx)
}

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

// Base64-encoded contents of testcasesE.txt.
const testcasesE64 = "MTAwCjggYmFjYmJhYWEKMSBiCjE4IGJhYWNjYmJhYWJhYWNiYmFhYgoxMCBjY2JhY2JjYmNhCjYgYWJiYWNiCjEgYgoxOSBjYmNhYmJjYmJiYWFiYmFhYWJjCjkgY2NjYmNiYWNhCjMgYmFjCjE1IGJhYmJjY2JjY2FiYWFjYQo5IGNjYWFiYWJiYQoyIGJjCjMgYmNjCjExIGFiYmJhY2JjY2FiCjIwIGFiYmFiYmNhYmNhYmFiYWJiYmNhCjE1IGFiYWFhYWFjYWNjYWNhYwoxNiBjYWJhYWNiYmNhYmFhYmJiCjIgYWIKMTUgYWNiYWJhYWFiYmFjYWFiCjkgYWJhYmNhY2JjCjIgYmIKMyBiYWMKNiBiYmNiY2IKMTQgY2FjY2JiYmJhYmNjYWEKMTMgY2FiY2FhY2NiY2JjYgoxMyBhYWFhYWNiYWJiYWNhCjcgYWNjYmNiYwoxNSBhYWFjY2FiY2JjYWFiYWMKMSBiCjEyIGFhY2JiYWJiYWJiYQozIGFjYgoxNyBiYmNiYWNhY2JiYmFjYWJhYQo1IGJiYmJjCjEgYwo2IGFjYmFjYQoxNiBjY2JjYWNhYmJiYWFjYWNhCjIwIGJjYmNiY2JjY2JiYWJjYmNiYWJhCjIwIGFjY2FiYmFjYmNjYWFiY2FiY2JiCjYgYmFjYWNhCjEyIGFhYWNiY2NhYWFjYgo4IGNjYmNiY2FiCjggYmJhY2JjY2MKMTUgYmNiYmJiYmNhYmFjYmNjCjIwIGJiYmNhYWNhYmNiY2JhYmJhY2JjCjMgYWFhCjEwIGFhYWJjY2NjY2MKOSBhY2FhYWFhYmEKNCBiYWNiCjIwIGNjYWNhY2JhY2NhY2JjYWNjYWNiCjUgYWJiY2EKMTMgYWFiY2JjY2JiYmFiYQoxNCBhYmNhYmFhYWFjYWFhYwoxIGIKMTggYWJjY2NjYmJhY2NhYWJjYmJjCjE0IGNjYWNhYmNjY2JiY2JhCjggYmFjY2FjYmMKMTYgYmFhY2FiY2NjY2FhYWNhYgoyMCBhY2JhY2NhYWFjYmNhY2FhYmNhYQoxMiBhY2NiYmJjY2JhYWIKMTYgYWNjYmNjYmFjY2JiYWFhYgoxNCBiYmJjYmJjYWJiY2FhYgo5IGNiY2NhY2JhYQoyMCBhY2FiY2FjYmFjYmFiY2FhYWFjYQo0IGJiY2MKNSBjYWNhYwoxIGEKNyBiYWNiY2JjCjIwIGJhYWFhYWFjYWNiY2NjYmFhYWJiCjIwIGNiY2JiYmNiYWFhYWJhYmFiY2FjCjEyIGJiY2NhYWFiYWNhYQoxIGEKNyBiY2NiY2JhCjEwIGNiYWFiY2JiYmEKMTUgY2FiYWNjY2JhYWJiY2NhCjMgYmFjCjE3IGFhY2FjYmNhY2NhYmJhY2FjCjIgYWEKMjAgYmFiY2JjY2JhY2JjY2NjYmJjYmMKMiBjYQoxMiBjYWJiYWFjY2FhYWMKNSBjYWNiYQoyIGNhCjE4IGFiY2FiYmJhY2JiY2NhYWJiYQoxIGIKMTcgY2JhYWJjYWJhYmFjYWFhY2MKMyBhYWEKMTAgY2NjYmFhY2NiYwo2IGNhY2NjYQoyIGFhCjE0IGFhY2JjYmNiYWNjYWNjCjIwIGFjYWFhYmNjYmJjYmFjYmNiYWNiCjYgYmNjY2FjCjE4IGNhYmNjY2NjY2JiYWJhYmJiYwo3IGNiY2JiYmMKNiBiYWFiY2EKNiBiYWJhY2EKMiBhYwoxOCBhYWFiY2JiYWNiY2JiY2FiYWEKNiBhYmFiYWIKMiBjYQo="

type testCase struct {
    n int
    s string
}

// Embedded solution logic from 1822E.go.
func solve(tc testCase) string {
    n := tc.n
    s := tc.s
    if n%2 == 1 {
        return "-1"
    }
    freq := make([]int, 26)
    for i := 0; i < n; i++ {
        freq[s[i]-'a']++
    }
    maxFreq := 0
    for _, v := range freq {
        if v > maxFreq {
            maxFreq = v
        }
    }
    if maxFreq > n/2 {
        return "-1"
    }
    pairCount := make([]int, 26)
    same := 0
    for i := 0; i < n/2; i++ {
        if s[i] == s[n-1-i] {
            same++
            pairCount[s[i]-'a']++
        }
    }
    maxPair := 0
    for _, v := range pairCount {
        if v > maxPair {
            maxPair = v
        }
    }
    ans := (same + 1) / 2
    if maxPair > ans {
        ans = maxPair
    }
    return strconv.Itoa(ans)
}

func parseTestcases() ([]testCase, error) {
    raw, err := base64.StdEncoding.DecodeString(testcasesE64)
    if err != nil {
        return nil, err
    }
    sc := bufio.NewScanner(bytes.NewReader(raw))
    sc.Split(bufio.ScanWords)
    if !sc.Scan() {
        return nil, fmt.Errorf("invalid test data")
    }
    t, err := strconv.Atoi(sc.Text())
    if err != nil {
        return nil, fmt.Errorf("parse t: %v", err)
    }
    cases := make([]testCase, 0, t)
    for i := 0; i < t; i++ {
        if !sc.Scan() {
            return nil, fmt.Errorf("case %d missing n", i+1)
        }
        n, err := strconv.Atoi(sc.Text())
        if err != nil {
            return nil, fmt.Errorf("case %d n: %v", i+1, err)
        }
        if !sc.Scan() {
            return nil, fmt.Errorf("case %d missing s", i+1)
        }
        s := sc.Text()
        cases = append(cases, testCase{n: n, s: s})
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
        fmt.Println("usage: go run verifierE.go /path/to/binary")
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
        input := fmt.Sprintf("1\n%d\n%s\n", tc.n, tc.s)
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


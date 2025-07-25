package main

import (
    "bufio"
    "bytes"
    "fmt"
    "math/rand"
    "os"
    "os/exec"
    "strings"
    "time"
)

type testCase struct {
    names []string
    pseuds []string
    ans int
}

func lcp(a, b string) int {
    n := len(a)
    if len(b) < n {
        n = len(b)
    }
    i := 0
    for i < n && a[i] == b[i] {
        i++
    }
    return i
}

func brute(tc testCase) int {
    n := len(tc.names)
    idx := make([]int, n)
    for i := range idx {
        idx[i] = i
    }
    best := 0
    perm := make([]int, n)

    var dfs func(pos int, used []bool, cur int)
    dfs = func(pos int, used []bool, cur int) {
        if pos == n {
            if cur > best {
                best = cur
            }
            return
        }
        for i := 0; i < n; i++ {
            if used[i] {
                continue
            }
            used[i] = true
            perm[pos] = i
            dfs(pos+1, used, cur+lcp(tc.names[pos], tc.pseuds[i]))
            used[i] = false
        }
    }
    dfs(0, make([]bool, n), 0)
    return best
}

func genTests() []testCase {
    rand.Seed(time.Now().UnixNano())
    tests := make([]testCase, 100)
    letters := []rune("abcdefghijklmnopqrstuvwxyz")
    for i := range tests {
        n := rand.Intn(5) + 1 // 1..5
        names := make([]string, n)
        pseuds := make([]string, n)
        for j := 0; j < n; j++ {
            ln := rand.Intn(3) + 1 // 1..3
            var sb strings.Builder
            for k := 0; k < ln; k++ {
                sb.WriteRune(letters[rand.Intn(len(letters))])
            }
            names[j] = sb.String()
        }
        for j := 0; j < n; j++ {
            ln := rand.Intn(3) + 1
            var sb strings.Builder
            for k := 0; k < ln; k++ {
                sb.WriteRune(letters[rand.Intn(len(letters))])
            }
            pseuds[j] = sb.String()
        }
        tc := testCase{names: names, pseuds: pseuds}
        tc.ans = brute(tc)
        tests[i] = tc
    }
    return tests
}

func main() {
    if len(os.Args) != 2 {
        fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
        os.Exit(1)
    }
    bin := os.Args[1]
    tests := genTests()
    for i, tc := range tests {
        var input bytes.Buffer
        fmt.Fprintln(&input, len(tc.names))
        for _, s := range tc.names {
            fmt.Fprintln(&input, s)
        }
        for _, s := range tc.pseuds {
            fmt.Fprintln(&input, s)
        }
        cmd := exec.Command(bin)
        cmd.Stdin = bytes.NewReader(input.Bytes())
        var out bytes.Buffer
        cmd.Stdout = &out
        cmd.Stderr = &out
        if err := cmd.Run(); err != nil {
            fmt.Fprintf(os.Stderr, "runtime error on test %d: %v\noutput:\n%s\n", i+1, err, out.String())
            os.Exit(1)
        }
        scanner := bufio.NewScanner(bytes.NewReader(out.Bytes()))
        scanner.Split(bufio.ScanWords)
        if !scanner.Scan() {
            fmt.Fprintf(os.Stderr, "no output on test %d\n", i+1)
            os.Exit(1)
        }
        var got int
        if _, err := fmt.Sscanf(scanner.Text(), "%d", &got); err != nil {
            fmt.Fprintf(os.Stderr, "invalid output on test %d: %s\n", i+1, scanner.Text())
            os.Exit(1)
        }
        if got != tc.ans {
            fmt.Fprintf(os.Stderr, "wrong answer on test %d: expected %d got %d\n", i+1, tc.ans, got)
            os.Exit(1)
        }
    }
    fmt.Println("Accepted")
}


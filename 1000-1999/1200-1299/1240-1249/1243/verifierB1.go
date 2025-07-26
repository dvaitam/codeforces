package main

import (
    "bufio"
    "bytes"
    "fmt"
    "os"
    "os/exec"
    "strings"
)

func solveCase(n int, s, t string) string {
    mism := make([]int, 0, 2)
    for i := 0; i < n; i++ {
        if s[i] != t[i] {
            mism = append(mism, i)
        }
    }
    if len(mism) == 2 && s[mism[0]] == s[mism[1]] && t[mism[0]] == t[mism[1]] {
        return "Yes"
    }
    return "No"
}

func runProg(bin, input string) (string, error) {
    var cmd *exec.Cmd
    if strings.HasSuffix(bin, ".go") {
        cmd = exec.Command("go", "run", bin)
    } else {
        cmd = exec.Command(bin)
    }
    cmd.Stdin = strings.NewReader(input)
    var out bytes.Buffer
    var stderr bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = &stderr
    if err := cmd.Run(); err != nil {
        return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
    }
    return strings.TrimSpace(out.String()), nil
}

func main() {
    if len(os.Args) == 3 && os.Args[1] == "--" {
        os.Args = append([]string{os.Args[0]}, os.Args[2])
    }
    if len(os.Args) != 2 {
        fmt.Println("Usage: go run verifierB1.go /path/to/binary")
        os.Exit(1)
    }
    bin := os.Args[1]

    f, err := os.Open("testcasesB1.txt")
    if err != nil {
        fmt.Fprintf(os.Stderr, "could not open testcasesB1.txt: %v\n", err)
        os.Exit(1)
    }
    defer f.Close()
    reader := bufio.NewReader(f)

    var t int
    if _, err := fmt.Fscan(reader, &t); err != nil {
        fmt.Fprintf(os.Stderr, "failed to read test count: %v\n", err)
        os.Exit(1)
    }

    for caseNum := 1; caseNum <= t; caseNum++ {
        var n int
        var s, tt string
        if _, err := fmt.Fscan(reader, &n); err != nil {
            fmt.Fprintf(os.Stderr, "case %d: read n: %v\n", caseNum, err)
            os.Exit(1)
        }
        if _, err := fmt.Fscan(reader, &s); err != nil {
            fmt.Fprintf(os.Stderr, "case %d: read s: %v\n", caseNum, err)
            os.Exit(1)
        }
        if _, err := fmt.Fscan(reader, &tt); err != nil {
            fmt.Fprintf(os.Stderr, "case %d: read t: %v\n", caseNum, err)
            os.Exit(1)
        }

        var input strings.Builder
        fmt.Fprintf(&input, "1\n%d\n%s\n%s\n", n, s, tt)

        want := solveCase(n, s, tt)
        got, err := runProg(bin, input.String())
        if err != nil {
            fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", caseNum, err, input.String())
            os.Exit(1)
        }
        if strings.TrimSpace(got) != strings.TrimSpace(want) {
            fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", caseNum, want, got, input.String())
            os.Exit(1)
        }
    }

    fmt.Printf("All %d tests passed\n", t)
}


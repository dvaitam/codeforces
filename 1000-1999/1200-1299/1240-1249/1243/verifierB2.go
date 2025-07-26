package main

import (
    "bufio"
    "bytes"
    "fmt"
    "os"
    "os/exec"
    "path/filepath"
    "strings"
)

func buildOracle() (string, error) {
    dir, err := os.Getwd()
    if err != nil {
        return "", err
    }
    oracle := filepath.Join(dir, "oracleB2")
    cmd := exec.Command("go", "build", "-o", oracle, "1243B2.go")
    if out, err := cmd.CombinedOutput(); err != nil {
        return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
    }
    return oracle, nil
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
        fmt.Println("Usage: go run verifierB2.go /path/to/binary")
        os.Exit(1)
    }
    bin := os.Args[1]

    oracle, err := buildOracle()
    if err != nil {
        fmt.Fprintf(os.Stderr, "%v\n", err)
        os.Exit(1)
    }
    defer os.Remove(oracle)

    f, err := os.Open("testcasesB2.txt")
    if err != nil {
        fmt.Fprintf(os.Stderr, "could not open testcasesB2.txt: %v\n", err)
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

        expect, err := runProg(oracle, input.String())
        if err != nil {
            fmt.Fprintf(os.Stderr, "oracle failed on case %d: %v\n", caseNum, err)
            os.Exit(1)
        }
        got, err := runProg(bin, input.String())
        if err != nil {
            fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", caseNum, err, input.String())
            os.Exit(1)
        }
        if strings.TrimSpace(got) != strings.TrimSpace(expect) {
            fmt.Fprintf(os.Stderr, "case %d failed: expected\n%s\ngot\n%s\ninput:\n%s", caseNum, expect, got, input.String())
            os.Exit(1)
        }
    }

    fmt.Printf("All %d tests passed\n", t)
}


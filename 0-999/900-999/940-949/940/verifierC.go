package main

import (
    "bufio"
    "bytes"
    "fmt"
    "os"
    "os/exec"
    "strings"
)

func buildOracle() (string, error) {
    exe := "oracleC"
    cmd := exec.Command("go", "build", "-o", exe, "940C.go")
    if out, err := cmd.CombinedOutput(); err != nil {
        return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
    }
    return "./" + exe, nil
}

func runProgram(bin, input string) (string, error) {
    var cmd *exec.Cmd
    if strings.HasSuffix(bin, ".go") {
        cmd = exec.Command("go", "run", bin)
    } else {
        cmd = exec.Command(bin)
    }
    cmd.Stdin = strings.NewReader(input)
    var out bytes.Buffer
    var errb bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = &errb
    if err := cmd.Run(); err != nil {
        return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
    }
    return strings.TrimSpace(out.String()), nil
}

func main() {
    if len(os.Args) != 2 {
        fmt.Println("usage: go run verifierC.go /path/to/binary")
        os.Exit(1)
    }
    bin := os.Args[1]
    oracle, err := buildOracle()
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
    defer os.Remove(oracle)

    f, err := os.Open("testcasesC.txt")
    if err != nil {
        fmt.Fprintln(os.Stderr, "failed to open testcases:", err)
        os.Exit(1)
    }
    defer f.Close()

    scanner := bufio.NewScanner(f)
    idx := 0
    for scanner.Scan() {
        line := strings.TrimSpace(scanner.Text())
        if line == "" {
            continue
        }
        idx++
        input := line + "\n"
        exp, err := runProgram(oracle, input)
        if err != nil {
            fmt.Fprintf(os.Stderr, "oracle error on case %d: %v\n", idx, err)
            os.Exit(1)
        }
        got, err := runProgram(bin, input)
        if err != nil {
            fmt.Fprintf(os.Stderr, "case %d: %v\n", idx, err)
            os.Exit(1)
        }
        if got != exp {
            fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx, exp, got)
            os.Exit(1)
        }
    }
    if err := scanner.Err(); err != nil {
        fmt.Fprintln(os.Stderr, "scanner error:", err)
        os.Exit(1)
    }
    fmt.Printf("All %d tests passed\n", idx)
}


package main

import (
    "bufio"
    "bytes"
    "fmt"
    "os"
    "os/exec"
    "path/filepath"
    "strconv"
    "strings"
)

func buildOracle() (string, error) {
    dir, err := os.Getwd()
    if err != nil {
        return "", err
    }
    oracle := filepath.Join(dir, "oracleC")
    cmd := exec.Command("go", "build", "-o", oracle, "607C.go")
    if out, err := cmd.CombinedOutput(); err != nil {
        return "", fmt.Errorf("oracle build failed: %v\n%s", err, out)
    }
    return oracle, nil
}

func lineToInput(line string) (string, error) {
    fields := strings.Fields(line)
    if len(fields) != 3 {
        return "", fmt.Errorf("expected 3 fields, got %d", len(fields))
    }
    if _, err := strconv.Atoi(fields[0]); err != nil {
        return "", err
    }
    return fmt.Sprintf("%s\n%s\n%s\n", fields[0], fields[1], fields[2]), nil
}

func run(bin string, input string) (string, string, error) {
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
    err := cmd.Run()
    return out.String(), stderr.String(), err
}

func main() {
    if len(os.Args) != 2 {
        fmt.Println("usage: go run verifierC.go /path/to/binary")
        os.Exit(1)
    }
    target := os.Args[1]
    oracle, err := buildOracle()
    if err != nil {
        fmt.Fprintf(os.Stderr, "%v\n", err)
        os.Exit(1)
    }
    defer os.Remove(oracle)

    f, err := os.Open("testcasesC.txt")
    if err != nil {
        fmt.Fprintf(os.Stderr, "open testcases: %v\n", err)
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
        input, err := lineToInput(line)
        if err != nil {
            fmt.Printf("test %d: parse error: %v\n", idx, err)
            os.Exit(1)
        }
        expOut, expErr, err := run(oracle, input)
        if err != nil {
            fmt.Printf("oracle run failed on test %d: %v\nstderr: %s\n", idx, err, expErr)
            os.Exit(1)
        }
        gotOut, gotErr, err := run(target, input)
        if err != nil {
            fmt.Printf("test %d: runtime error: %v\nstderr: %s\n", idx, err, gotErr)
            os.Exit(1)
        }
        if strings.TrimSpace(gotOut) != strings.TrimSpace(expOut) {
            fmt.Printf("test %d failed\nexpected: %s\n got: %s\n", idx, strings.TrimSpace(expOut), strings.TrimSpace(gotOut))
            os.Exit(1)
        }
    }
    if err := scanner.Err(); err != nil {
        fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
        os.Exit(1)
    }
    fmt.Printf("All %d tests passed\n", idx)
}

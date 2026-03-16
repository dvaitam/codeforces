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
    oracle := filepath.Join(dir, "oracleB")
    cmd := exec.Command("go", "build", "-o", oracle, "607B.go")
    if out, err := cmd.CombinedOutput(); err != nil {
        return "", fmt.Errorf("oracle build failed: %v\n%s", err, out)
    }
    return oracle, nil
}

func lineToInput(line string) (string, error) {
    fields := strings.Fields(line)
    if len(fields) == 0 {
        return "", fmt.Errorf("empty line")
    }
    n, err := strconv.Atoi(fields[0])
    if err != nil {
        return "", err
    }
    if len(fields) != 1+n {
        return "", fmt.Errorf("expected %d numbers, got %d", 1+n, len(fields))
    }
    var b strings.Builder
    b.WriteString(fields[0])
    b.WriteByte('\n')
    for i := 0; i < n; i++ {
        b.WriteString(fields[1+i])
        if i+1 < n {
            b.WriteByte(' ')
        }
    }
    b.WriteByte('\n')
    return b.String(), nil
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
        fmt.Println("usage: go run verifierB.go /path/to/binary")
        os.Exit(1)
    }
    target := os.Args[1]
    oracle, err := buildOracle()
    if err != nil {
        fmt.Fprintf(os.Stderr, "%v\n", err)
        os.Exit(1)
    }
    defer os.Remove(oracle)

    const testcasesBRaw = `5 5 1 5 3 3
7 4 1 6 6 5 6 6
9 7 8 8 5 8 4 6 5 1
1 1
3 2 1 2
1 1
2 2 1
10 7 9 4 8 4 6 10 2 10 2
6 3 5 4 3 3 1
9 1 4 6 2 4 9 6 4 4
5 3 3 5 4 3
8 6 4 1 5 2 1 8 8
8 1 7 8 8 8 2 2 2
4 1 2 4 2
8 2 7 7 1 3 4 8 4
3 2 2 2
7 1 5 3 5 5 7 2
5 4 5 5 4 5
5 3 2 1 1 5
2 1 2
4 2 3 1 4
1 1
7 6 3 1 6 5 3 2
9 5 4 4 2 9 5 6 4 6
8 5 3 3 1 6 6 1 3
7 2 2 5 1 2 7 2
8 4 4 3 4 7 6 3 8
2 1 2
8 8 5 1 4 3 8 8 6
2 2 1
10 7 4 6 5 7 1 4 1 6 4
6 4 6 6 6 6 2
5 3 2 3 1 3
10 9 1 3 6 1 8 1 1 4 1
1 1
6 1 1 3 6 4 2
4 4 4 2 3
5 2 3 4 4 1
7 3 5 5 7 6 6 6
8 1 2 7 7 3 1 3 3
2 2 1
3 1 1 1
9 3 2 7 2 8 3 1 5 6
7 1 6 1 4 1 3 3
3 2 1 3
6 2 6 4 3 3 4
7 1 3 5 3 5 4 1
9 9 5 1 8 7 2 7 6 8
1 1
5 1 3 5 3 2
9 9 6 7 5 4 2 6 4 9
6 2 2 3 6 1 5
1 1
6 3 3 6 3 3 4
7 5 4 2 1 7 2 5
1 1
3 2 1 3
8 5 4 2 7 5 3 3 2
3 3 1 3
9 7 7 5 5 5 1 7 5 5
9 9 9 6 6 4 7 3 1 9
3 3 3 3
7 3 4 1 5 4 6 5
4 1 3 2 2
6 6 6 4 1 6 6
4 2 3 2 4
2 2 1
8 2 4 3 8 2 7 7 5
5 4 3 5 3 1
5 1 4 1 3 2
7 4 4 7 6 6 6 4
1 1
6 5 2 5 6 3 3
1 1
8 3 1 2 6 6 1 2 4
2 2 1
6 1 3 4 2 6 3
7 6 2 5 2 4 7 3
9 1 3 3 3 8 1 9 1 9
7 2 3 7 5 7 1 1
9 3 5 4 5 6 5 5 9 8
3 2 3 1
1 1
9 1 6 2 4 8 4 8 9 3
6 6 2 4 5 1 5
2 2 1
2 1 2
10 6 10 8 6 7 9 6 2 3 6
1 1
3 1 2 3
4 1 4 1 3
7 1 5 7 7 6 2 3
2 2 1
8 6 5 5 8 7 3 1 8
5 2 4 1 3 1
2 1 2
1 1
7 5 6 6 1 3 4 7
9 8 8 2 1 9 7 5 1 9
2 1 2`

    scanner := bufio.NewScanner(strings.NewReader(testcasesBRaw))
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

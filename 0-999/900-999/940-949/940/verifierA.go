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
    exe := "oracleA"
    cmd := exec.Command("go", "build", "-o", exe, "940A.go")
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
        fmt.Println("usage: go run verifierA.go /path/to/binary")
        os.Exit(1)
    }
    bin := os.Args[1]
    oracle, err := buildOracle()
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
    defer os.Remove(oracle)

    const testcasesRaw = `7 6 2 9 17 16 13 10 16
6 9 7 17 5 10 5 4
10 4 18 20 5 10 4 3 11 16 18 4
6 6 11 20 7 18 16 15
9 4 2 18 1 3 13 1 20 16 11
4 5 3 7 19 8
4 2 18 15 3 3
6 8 16 4 10 18 10 4
9 5 18 7 20 18 19 10 15 3 20
7 5 19 8 10 6 7 6 2
10 10 9 16 3 3 5 5 2 3 18 13
9 4 17 8 7 19 14 19 9 15 16
6 1 11 20 4 16 19 11
4 3 1 9 4 8
6 2 11 14 2 4 5 8
1 9 18
10 10 3 1 4 7 20 19 4 13 3 12
2 0 20 1
4 2 4 16 7 2
1 8 14
10 1 9 3 8 3 10 12 14 6 2 17
8 0 20 4 13 7 9 12 16 19
3 10 7 2 6
3 5 17 9 4
10 7 6 1 16 14 19 17 10 12 13 9
3 8 1 15 3
6 0 18 9 5 8 16 12
10 4 12 19 20 5 10 13 14 3 1 20
4 5 6 8 8 15
7 10 19 14 2 13 19 14 2
3 7 3 9 6
8 8 16 18 20 1 2 16 11 10
8 0 14 7 18 3 5 1 13 14
6 0 7 1 1 17 20 4
4 1 20 7 10 9
3 1 16 13 3
1 4 15
2 4 5 17
6 1 5 9 1 2 2 7
5 8 11 12 19 2 20
8 10 15 14 12 18 6 7 13 19
5 0 5 5 9 11 11
6 1 11 20 2 2 9 6
3 9 10 12 13
9 2 10 4 16 8 2 10 6 17 3
5 6 11 10 14 4 4
9 7 16 11 11 4 16 4 16 14 2
5 5 5 6 19 13 3
2 1 7 8
1 6 1
2 6 18 17
5 7 16 19 7 14 3
6 3 9 19 6 14 7 12
2 1 1 17
8 10 7 4 16 13 9 7 2 7
10 2 4 7 15 13 12 18 5 4 20 16
3 9 13 14 17
8 10 11 16 16 7 18 20 8 1
6 5 11 2 17 5 9 20
3 6 19 10 16
2 1 17 2
2 3 5 2
5 0 15 11 6 5 15
6 8 13 17 17 2 19 3
9 9 3 14 7 10 18 20 14 16 13
10 9 8 1 1 6 10 17 19 9 11 3
8 4 10 14 13 13 2 6 5 8
5 5 2 2 16 14 5
8 9 3 5 12 14 2 20 15 13
8 0 4 16 5 1 2 20 20 5
6 1 18 12 7 13 16 4
1 9 15
10 10 11 4 20 10 5 13 10 4 17 7
1 6 15
6 3 15 12 3 2 2 16
5 0 17 19 19 7 8
2 10 17 17
7 8 10 4 5 14 19 14 3
2 6 3 4
7 2 1 15 14 14 1 16 11
5 1 12 3 4 12 1
6 5 6 1 8 12 3 20
3 3 1 7 4
1 4 12
1 9 8
3 2 15 4 16
6 4 5 1 7 12 11 16
5 4 18 11 6 19 3
2 8 19 10
3 6 5 5 8
6 8 8 8 6 10 12 14
1 2 20
1 6 3
2 2 14 10
9 6 5 19 14 10 12 3 8 15 12
9 0 13 14 1 14 11 15 7 12 10
8 1 6 4 9 4 18 20 5 15
7 2 14 14 6 8 15 11 17
3 5 15 3 16
4 4 1 15 20 15`

    scanner := bufio.NewScanner(strings.NewReader(testcasesRaw))
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


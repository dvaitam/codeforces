package main

import (
    "bufio"
    "bytes"
    "fmt"
    "os"
    "os/exec"
    "sort"
    "strconv"
    "strings"
)

func possible(a []int) bool {
    if len(a) < 3 {
        return false
    }
    b := make([]int, len(a))
    copy(b, a)
    sort.Ints(b)
    return b[len(b)-1] < b[len(b)-2]+b[len(b)-3]
}

func run(bin, input string) (string, error) {
    var cmd *exec.Cmd
    if strings.HasSuffix(bin, ".go") {
        cmd = exec.Command("go", "run", bin)
    } else {
        cmd = exec.Command(bin)
    }
    cmd.Stdin = strings.NewReader(input)
    var out bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = &out
    err := cmd.Run()
    return out.String(), err
}

func main() {
    if len(os.Args) != 2 {
        fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
        os.Exit(1)
    }
    bin := os.Args[1]
    file, err := os.Open("testcasesB.txt")
    if err != nil {
        fmt.Fprintf(os.Stderr, "failed to open testcases: %v\n", err)
        os.Exit(1)
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    idx := 0
    for scanner.Scan() {
        line := strings.TrimSpace(scanner.Text())
        if line == "" {
            continue
        }
        idx++
        parts := strings.Fields(line)
        if len(parts) < 2 {
            fmt.Fprintf(os.Stderr, "invalid test case on line %d\n", idx)
            os.Exit(1)
        }
        n, err := strconv.Atoi(parts[0])
        if err != nil || n != len(parts)-1 {
            if err == nil && n != len(parts)-1 {
                err = fmt.Errorf("expected %d numbers got %d", n, len(parts)-1)
            }
            fmt.Fprintf(os.Stderr, "invalid n on line %d: %v\n", idx, err)
            os.Exit(1)
        }
        a := make([]int, n)
        for i := 0; i < n; i++ {
            v, err := strconv.Atoi(parts[i+1])
            if err != nil {
                fmt.Fprintf(os.Stderr, "invalid number on line %d\n", idx)
                os.Exit(1)
            }
            a[i] = v
        }
        input := fmt.Sprintf("%d\n%s\n", n, strings.Join(parts[1:], " "))
        out, err := run(bin, input)
        if err != nil {
            fmt.Fprintf(os.Stderr, "test %d: runtime error: %v\n%s", idx, err, out)
            os.Exit(1)
        }
        fields := strings.Fields(out)
        if len(fields) == 0 {
            fmt.Fprintf(os.Stderr, "test %d: no output\n", idx)
            os.Exit(1)
        }
        if fields[0] == "NO" {
            if possible(a) {
                fmt.Fprintf(os.Stderr, "test %d: expected YES but got NO\n", idx)
                os.Exit(1)
            }
            continue
        }
        if fields[0] != "YES" {
            fmt.Fprintf(os.Stderr, "test %d: first token should be YES or NO\n", idx)
            os.Exit(1)
        }
        if len(fields)-1 != n {
            fmt.Fprintf(os.Stderr, "test %d: expected %d numbers got %d\n", idx, n, len(fields)-1)
            os.Exit(1)
        }
        arr := make([]int, n)
        for i := 0; i < n; i++ {
            v, err := strconv.Atoi(fields[i+1])
            if err != nil {
                fmt.Fprintf(os.Stderr, "test %d: invalid output number\n", idx)
                os.Exit(1)
            }
            arr[i] = v
        }
        counts := make(map[int]int)
        for _, v := range a {
            counts[v]++
        }
        for _, v := range arr {
            counts[v]--
        }
        for _, c := range counts {
            if c != 0 {
                fmt.Fprintf(os.Stderr, "test %d: numbers mismatch\n", idx)
                os.Exit(1)
            }
        }
        for i := 0; i < n; i++ {
            left := arr[(i-1+n)%n]
            right := arr[(i+1)%n]
            if arr[i] >= left+right {
                fmt.Fprintf(os.Stderr, "test %d: condition failed at index %d\n", idx, i)
                os.Exit(1)
            }
        }
        if !possible(a) {
            fmt.Fprintf(os.Stderr, "test %d: expected NO but got YES\n", idx)
            os.Exit(1)
        }
    }
    if err := scanner.Err(); err != nil {
        fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
        os.Exit(1)
    }
    fmt.Printf("All %d tests passed\n", idx)
}


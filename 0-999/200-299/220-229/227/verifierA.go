package main

import (
    "bufio"
    "bytes"
    "fmt"
    "os"
    "os/exec"
    "strings"
)

func runCandidate(bin string, input string) (string, error) {
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

func expected(xa, ya, xb, yb, xc, yc int64) string {
    vx := xb - xa
    vy := yb - ya
    wx := xc - xb
    wy := yc - yb
    cross := vx*wy - vy*wx
    if cross > 0 {
        return "LEFT"
    }
    if cross < 0 {
        return "RIGHT"
    }
    return "TOWARDS"
}

func main() {
    if len(os.Args) != 2 {
        fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
        os.Exit(1)
    }
    bin := os.Args[1]
    f, err := os.Open("testcasesA.txt")
    if err != nil {
        fmt.Fprintf(os.Stderr, "failed to open testcases: %v\n", err)
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
        var xa, ya, xb, yb, xc, yc int64
        fmt.Sscan(line, &xa, &ya, &xb, &yb, &xc, &yc)
        input := fmt.Sprintf("%d %d\n%d %d\n%d %d\n", xa, ya, xb, yb, xc, yc)
        want := expected(xa, ya, xb, yb, xc, yc)
        got, err := runCandidate(bin, input)
        if err != nil {
            fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx, err)
            os.Exit(1)
        }
        if got != want {
            fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx, want, got)
            os.Exit(1)
        }
    }
    if err := scanner.Err(); err != nil {
        fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
        os.Exit(1)
    }
    fmt.Printf("All %d tests passed\n", idx)
}


package main

import (
    "bytes"
    "fmt"
    "math/rand"
    "os"
    "os/exec"
    "path/filepath"
    "runtime"
    "strings"
    "time"
)

func runBinary(path string, input string) (string, error) {
    var cmd *exec.Cmd
    if strings.HasSuffix(path, ".go") {
        cmd = exec.Command("go", "run", path)
    } else {
        cmd = exec.Command(path)
    }
    cmd.Stdin = strings.NewReader(input)
    out, err := cmd.CombinedOutput()
    return strings.TrimSpace(string(out)), err
}

func referenceOutput(input string) (string, error) {
    _, file, _, ok := runtime.Caller(0)
    if !ok {
        return "", fmt.Errorf("cannot determine caller")
    }
    dir := filepath.Dir(file)
    path := filepath.Join(dir, "1651D.go")
    return runBinary(path, input)
}

func uniquePairs(n int) [][2]int {
    used := make(map[[2]int]bool)
    res := make([][2]int, 0, n)
    for len(res) < n {
        x := rand.Intn(50) + 1
        y := rand.Intn(50) + 1
        p := [2]int{x, y}
        if !used[p] {
            used[p] = true
            res = append(res, p)
        }
    }
    return res
}

func main() {
    args := os.Args[1:]
    if len(args) > 0 && args[0] == "--" {
        args = args[1:]
    }
    if len(args) != 1 {
        fmt.Fprintln(os.Stderr, "usage: verifierD <binary>")
        os.Exit(1)
    }
    path := args[0]
    rand.Seed(time.Now().UnixNano())

    for test := 0; test < 100; test++ {
        n := rand.Intn(4) + 1
        pts := uniquePairs(n)
        var buf bytes.Buffer
        fmt.Fprintln(&buf, n)
        for _, p := range pts {
            fmt.Fprintf(&buf, "%d %d\n", p[0], p[1])
        }
        input := buf.String()
        want, err := referenceOutput(input)
        if err != nil {
            fmt.Fprintf(os.Stderr, "failed to run reference: %v\n", err)
            os.Exit(1)
        }
        got, err := runBinary(path, input)
        if err != nil {
            fmt.Fprintf(os.Stderr, "test %d: runtime error: %v\n", test+1, err)
            os.Exit(1)
        }
        if strings.TrimSpace(got) != strings.TrimSpace(want) {
            fmt.Fprintf(os.Stderr, "test %d failed\ninput:\n%s\nexpected:\n%s\ngot:\n%s\n", test+1, input, want, got)
            os.Exit(1)
        }
    }
    fmt.Println("All tests passed")
}


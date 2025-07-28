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
    _, file, _, _ := runtime.Caller(0)
    path := filepath.Join(filepath.Dir(file), "1651E.go")
    return runBinary(path, input)
}

func buildGraph(n int) [][2]int {
    edges := make([][2]int, 0, 2*n)
    for i := 1; i <= n; i++ {
        edges = append(edges, [2]int{i, n + i})
        j := i % n + 1
        edges = append(edges, [2]int{i, n + j})
    }
    return edges
}

func main() {
    args := os.Args[1:]
    if len(args) > 0 && args[0] == "--" {
        args = args[1:]
    }
    if len(args) != 1 {
        fmt.Fprintln(os.Stderr, "usage: verifierE <binary>")
        os.Exit(1)
    }
    path := args[0]
    rand.Seed(time.Now().UnixNano())

    for test := 0; test < 100; test++ {
        n := rand.Intn(3) + 2
        edges := buildGraph(n)
        var buf bytes.Buffer
        fmt.Fprintln(&buf, n)
        for _, e := range edges {
            fmt.Fprintf(&buf, "%d %d\n", e[0], e[1])
        }
        input := buf.String()
        want, err := referenceOutput(input)
        if err != nil {
            fmt.Fprintf(os.Stderr, "reference error: %v\n", err)
            os.Exit(1)
        }
        got, err := runBinary(path, input)
        if err != nil {
            fmt.Fprintf(os.Stderr, "test %d runtime error: %v\n", test+1, err)
            os.Exit(1)
        }
        if strings.TrimSpace(got) != strings.TrimSpace(want) {
            fmt.Fprintf(os.Stderr, "test %d failed\ninput:\n%s\nexpected:%s\ngot:%s\n", test+1, input, want, got)
            os.Exit(1)
        }
    }
    fmt.Println("All tests passed")
}


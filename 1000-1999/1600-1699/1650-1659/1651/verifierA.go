package main

import (
    "fmt"
    "math/rand"
    "os"
    "os/exec"
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

func expected(n int) string {
    return fmt.Sprintf("%d", (1<<uint(n))-1)
}

func main() {
    args := os.Args[1:]
    if len(args) > 0 && args[0] == "--" {
        args = args[1:]
    }
    if len(args) != 1 {
        fmt.Fprintln(os.Stderr, "usage: verifierA <binary>")
        os.Exit(1)
    }
    path := args[0]
    rand.Seed(time.Now().UnixNano())

    for i := 0; i < 100; i++ {
        n := rand.Intn(30) + 1
        input := fmt.Sprintf("1\n%d\n", n)
        want := expected(n)
        got, err := runBinary(path, input)
        if err != nil {
            fmt.Fprintf(os.Stderr, "test %d: runtime error: %v\n", i+1, err)
            os.Exit(1)
        }
        if strings.TrimSpace(got) != want {
            fmt.Fprintf(os.Stderr, "test %d failed: n=%d expected %s got %s\n", i+1, n, want, got)
            os.Exit(1)
        }
    }
    fmt.Println("All tests passed")
}


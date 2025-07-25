package main

import (
    "fmt"
    "os"
    "os/exec"
    "strings"
)

func run(bin string, input string) (string, error) {
    cmd := exec.Command(bin)
    cmd.Stdin = strings.NewReader(input)
    out, err := cmd.CombinedOutput()
    return string(out), err
}

func expected(n int) string {
    return fmt.Sprintf("%d\n", n*(n+1)/2+1)
}

func main() {
    if len(os.Args) < 2 {
        fmt.Println("usage: go run verifierA.go /path/to/binary")
        os.Exit(1)
    }
    bin := os.Args[1]
    for i := 0; i <= 100; i++ {
        in := fmt.Sprintf("%d\n", i)
        want := expected(i)
        out, err := run(bin, in)
        if err != nil {
            fmt.Printf("test %d runtime error: %v\n", i, err)
            os.Exit(1)
        }
        out = strings.TrimSpace(out)
        want = strings.TrimSpace(want)
        if out != want {
            fmt.Printf("test %d failed: input %q expected %q got %q\n", i, strings.TrimSpace(in), want, out)
            os.Exit(1)
        }
    }
    fmt.Println("OK")
}

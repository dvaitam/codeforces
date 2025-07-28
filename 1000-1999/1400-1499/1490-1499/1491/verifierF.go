package main

import (
    "bytes"
    "fmt"
    "os"
    "os/exec"
    "strings"
)

func run(bin string) (string, error) {
    cmd := exec.Command(bin)
    var out bytes.Buffer
    var stderr bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = &stderr
    if err := cmd.Run(); err != nil {
        return "", fmt.Errorf("%v\n%s", err, stderr.String())
    }
    return strings.TrimSpace(out.String()), nil
}

func main() {
    if len(os.Args) != 2 {
        fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
        os.Exit(1)
    }
    candidate := os.Args[1]
    expected := "Problem F is interactive and cannot be automatically solved."
    for i := 0; i < 100; i++ {
        out, err := run(candidate)
        if err != nil {
            fmt.Fprintf(os.Stderr, "candidate runtime error on repetition %d: %v\n", i+1, err)
            os.Exit(1)
        }
        if out != expected {
            fmt.Fprintf(os.Stderr, "unexpected output on repetition %d: %q\n", i+1, out)
            os.Exit(1)
        }
    }
    fmt.Println("All tests passed")
}

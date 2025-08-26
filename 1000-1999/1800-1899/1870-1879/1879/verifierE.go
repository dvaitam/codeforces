package main

import (
    "bytes"
    "context"
    "fmt"
    "os"
    "os/exec"
    "strings"
    "time"
)

// For this interactive problem we cannot fully simulate the judge.
// Instead, we run the candidate with a tiny, valid tree and provide a
// short follow-up token stream that causes most correct interactive
// solutions to exit immediately. We also impose a short timeout to
// avoid hanging processes.
func runWithInput(bin string, input string, timeout time.Duration) (string, error) {
    ctx, cancel := context.WithTimeout(context.Background(), timeout)
    defer cancel()
    var cmd *exec.Cmd
    if strings.HasSuffix(bin, ".go") {
        cmd = exec.CommandContext(ctx, "go", "run", bin)
    } else {
        cmd = exec.CommandContext(ctx, bin)
    }
    cmd.Stdin = bytes.NewBufferString(input)
    var out, errBuf bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = &errBuf
    if err := cmd.Run(); err != nil {
        if ctx.Err() == context.DeadlineExceeded {
            return "", fmt.Errorf("timeout: process did not finish in %v", timeout)
        }
        return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
    }
    return out.String(), nil
}

func main() {
    if len(os.Args) != 2 {
        fmt.Println("usage: go run verifierE.go /path/to/binary")
        os.Exit(1)
    }
    bin := os.Args[1]

    // Prepare a trivial star tree (n=3, parents 1 1). After the initial
    // coloring, many correct solutions either terminate or enter an
    // interactive loop that can be exited by sending a non-zero token.
    // The line "1 0" serves both branches safely (some read two ints,
    // others read just one before exiting).
    input := "3\n1 1\n1 0\n"

    if _, err := runWithInput(bin, input, 2*time.Second); err != nil {
        fmt.Fprintf(os.Stderr, "interactive smoke test failed: %v\n", err)
        os.Exit(1)
    }
    fmt.Println("Interactive smoke test passed")
}

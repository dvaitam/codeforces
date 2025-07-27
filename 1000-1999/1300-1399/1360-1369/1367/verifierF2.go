package main

import (
    "bytes"
    "fmt"
    "math/rand"
    "os"
    "os/exec"
    "path/filepath"
    "strings"
    "time"
)

func compileRef() (string, error) {
    bin := filepath.Join(os.TempDir(), "1367F2_ref")
    cmd := exec.Command("go", "build", "-o", bin, "1367F2.go")
    out, err := cmd.CombinedOutput()
    if err != nil {
        return "", fmt.Errorf("compile reference: %v\n%s", err, out)
    }
    return bin, nil
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
    if err := cmd.Run(); err != nil {
        return "", fmt.Errorf("%v\n%s", err, out.String())
    }
    return strings.TrimSpace(out.String()), nil
}

func genCase(rng *rand.Rand) string {
    n := rng.Intn(50) + 1
    var sb strings.Builder
    sb.WriteString("1\n")
    sb.WriteString(fmt.Sprintf("%d\n", n))
    for i := 0; i < n; i++ {
        val := rng.Intn(100)
        sb.WriteString(fmt.Sprintf("%d", val))
        if i+1 < n {
            sb.WriteByte(' ')
        }
    }
    sb.WriteByte('\n')
    return sb.String()
}

func main() {
    if len(os.Args) != 2 {
        fmt.Println("usage: go run verifierF2.go /path/to/binary")
        os.Exit(1)
    }
    candidate := os.Args[1]
    ref, err := compileRef()
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
    defer os.Remove(ref)

    rng := rand.New(rand.NewSource(time.Now().UnixNano()))
    for i := 0; i < 100; i++ {
        input := genCase(rng)
        want, err := run(ref, input)
        if err != nil {
            fmt.Fprintf(os.Stderr, "reference failure: %v\n", err)
            os.Exit(1)
        }
        got, err := run(candidate, input)
        if err != nil {
            fmt.Fprintf(os.Stderr, "case %d runtime error: %v\ninput:\n%s", i+1, err, input)
            os.Exit(1)
        }
        if got != want {
            fmt.Fprintf(os.Stderr, "case %d failed\ninput:\n%sexpected:%s\ngot:%s\n", i+1, input, want, got)
            os.Exit(1)
        }
    }
    fmt.Println("All tests passed")
}


package main

import (
    "bytes"
    "fmt"
    "math/rand"
    "os"
    "os/exec"
    "strings"
    "time"
)

func buildOracle() (string, error) {
    oracle := "oracleI.bin"
    cmd := exec.Command("go", "build", "-o", oracle, "1491I.go")
    if out, err := cmd.CombinedOutput(); err != nil {
        return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
    }
    return oracle, nil
}

func run(bin, input string) (string, error) {
    cmd := exec.Command(bin)
    cmd.Stdin = strings.NewReader(input)
    var out bytes.Buffer
    var stderr bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = &stderr
    if err := cmd.Run(); err != nil {
        return "", fmt.Errorf("%v\n%s", err, stderr.String())
    }
    return strings.TrimSpace(out.String()), nil
}

func genCase(rng *rand.Rand) string {
    n := rng.Intn(7) + 4
    base := rng.Intn(1000000)
    var sb strings.Builder
    sb.WriteString(fmt.Sprintf("%d\n", n))
    for i := 0; i < n; i++ {
        b := base + 3*i
        a := base + 3*n + 3*i
        c := base + 6*n + 3*i
        sb.WriteString(fmt.Sprintf("%d %d %d\n", a, b, c))
    }
    return sb.String()
}

func main() {
    if len(os.Args) != 2 {
        fmt.Fprintln(os.Stderr, "usage: go run verifierI.go /path/to/binary")
        os.Exit(1)
    }
    candidate := os.Args[1]
    oracle, err := buildOracle()
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
    defer os.Remove(oracle)

    rng := rand.New(rand.NewSource(time.Now().UnixNano()))
    for tcase := 0; tcase < 100; tcase++ {
        input := genCase(rng)
        want, err := run(oracle, input)
        if err != nil {
            fmt.Fprintf(os.Stderr, "oracle error on case %d: %v\n", tcase+1, err)
            os.Exit(1)
        }
        got, err := run(candidate, input)
        if err != nil {
            fmt.Fprintf(os.Stderr, "candidate runtime error on case %d: %v\ninput:\n%s", tcase+1, err, input)
            os.Exit(1)
        }
        if strings.TrimSpace(got) != strings.TrimSpace(want) {
            fmt.Fprintf(os.Stderr, "case %d failed:\ninput:\n%s\nexpected:\n%s\ngot:\n%s\n", tcase+1, input, want, got)
            os.Exit(1)
        }
    }
    fmt.Println("All tests passed")
}

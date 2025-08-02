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

func buildOracle() (string, error) {
    dir, err := os.Getwd()
    if err != nil {
        return "", err
    }
    oracle := filepath.Join(dir, "oracleE")
    cmd := exec.Command("go", "build", "-o", oracle, "1991E.go")
    if out, err := cmd.CombinedOutput(); err != nil {
        return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
    }
    return oracle, nil
}

func genCase(rng *rand.Rand) string {
    n := rng.Intn(4) + 2
    edges := make([][2]int, 0, n)
    for i := 1; i < n; i++ {
        edges = append(edges, [2]int{i, i + 1})
    }
    if n >= 3 && rng.Intn(2) == 0 {
        edges = append(edges, [2]int{1, 3})
    }
    m := len(edges)
    var sb strings.Builder
    sb.WriteString("1\n")
    sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
    for _, e := range edges {
        sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
    }
    for i := 0; i < n; i++ {
        sb.WriteString(fmt.Sprintf("%d %d\n", rng.Intn(3)+1, rng.Intn(3)+1))
    }
    return sb.String()
}

func run(bin, input string) (string, error) {
    cmd := exec.Command(bin)
    cmd.Stdin = strings.NewReader(input)
    var out bytes.Buffer
    var errOut bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = &errOut
    if err := cmd.Run(); err != nil {
        return "", fmt.Errorf("%v\n%s", err, errOut.String())
    }
    return out.String(), nil
}

func main() {
    if len(os.Args) != 2 {
        fmt.Println("usage: go run verifierE.go /path/to/binary")
        os.Exit(1)
    }
    bin := os.Args[1]
    oracle, err := buildOracle()
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
    defer os.Remove(oracle)
    rng := rand.New(rand.NewSource(time.Now().UnixNano()))
    cases := make([]string, 0, 100)
    for len(cases) < 100 {
        cases = append(cases, genCase(rng))
    }
    for i, in := range cases {
        exp, err := run(oracle, in)
        if err != nil {
            fmt.Fprintf(os.Stderr, "oracle error on case %d: %v\n", i+1, err)
            os.Exit(1)
        }
        got, err := run(bin, in)
        if err != nil {
            fmt.Fprintf(os.Stderr, "runtime error on case %d: %v\ninput:\n%s", i+1, err, in)
            os.Exit(1)
        }
        if strings.TrimSpace(exp) != strings.TrimSpace(got) {
            fmt.Fprintf(os.Stderr, "case %d failed\nexpected: %s\ngot: %s\ninput:\n%s", i+1, exp, got, in)
            os.Exit(1)
        }
    }
    fmt.Println("All tests passed")
}


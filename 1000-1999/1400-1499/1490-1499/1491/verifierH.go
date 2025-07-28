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
    oracle := "oracleH.bin"
    cmd := exec.Command("go", "build", "-o", oracle, "1491H.go")
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
    n := rng.Intn(15) + 2
    q := rng.Intn(20) + 1
    var sb strings.Builder
    sb.WriteString(fmt.Sprintf("%d %d\n", n, q))
    for i := 2; i <= n; i++ {
        p := rng.Intn(i-1) + 1
        if i > 2 {
            sb.WriteByte(' ')
        }
        sb.WriteString(fmt.Sprint(p))
    }
    if n > 1 {
        sb.WriteByte('\n')
    }
    type2 := false
    sbQueries := strings.Builder{}
    for i := 0; i < q; i++ {
        t := 1
        if !type2 || rng.Intn(2) == 1 {
            t = 2
            type2 = true
        }
        if t == 1 {
            l := rng.Intn(n-1) + 2
            r := rng.Intn(n-l+1) + l
            x := rng.Intn(5) + 1
            sbQueries.WriteString(fmt.Sprintf("1 %d %d %d\n", l, r, x))
        } else {
            u := rng.Intn(n) + 1
            v := rng.Intn(n) + 1
            sbQueries.WriteString(fmt.Sprintf("2 %d %d\n", u, v))
        }
    }
    if !type2 {
        // ensure at least one type2
        u := 1
        v := n
        sbQueries.WriteString(fmt.Sprintf("2 %d %d\n", u, v))
    }
    sb.WriteString(sbQueries.String())
    return sb.String()
}

func main() {
    if len(os.Args) != 2 {
        fmt.Fprintln(os.Stderr, "usage: go run verifierH.go /path/to/binary")
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

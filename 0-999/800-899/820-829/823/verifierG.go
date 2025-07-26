package main

import (
    "bytes"
    "fmt"
    "math/rand"
    "os"
    "os/exec"
    "strings"
)

func buildRef() (string, error) {
    ref := "./refG.bin"
    cmd := exec.Command("go", "build", "-o", ref, "823G.go")
    if out, err := cmd.CombinedOutput(); err != nil {
        return "", fmt.Errorf("failed to build reference: %v\n%s", err, out)
    }
    return ref, nil
}

func runBinary(bin, input string) (string, error) {
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
        return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
    }
    return strings.TrimSpace(out.String()), nil
}

type Case struct{ input string }

func genCase(rng *rand.Rand) Case {
    n := rng.Intn(5) + 2
    m := rng.Intn(n * 2)
    var sb strings.Builder
    fmt.Fprintf(&sb, "%d %d\n", n, m)
    for i := 0; i < m; i++ {
        a := rng.Intn(n) + 1
        b := rng.Intn(n) + 1
        l := rng.Intn(5)
        r := l + rng.Intn(5) + 1
        fmt.Fprintf(&sb, "%d %d %d %d\n", a, b, l, r)
    }
    return Case{sb.String()}
}

func genCases() []Case {
    rng := rand.New(rand.NewSource(8236))
    cases := make([]Case, 100)
    for i := range cases {
        cases[i] = genCase(rng)
    }
    return cases
}

func runCase(bin, ref string, c Case) error {
    expected, err := runBinary(ref, c.input)
    if err != nil {
        return fmt.Errorf("reference failed: %v", err)
    }
    got, err := runBinary(bin, c.input)
    if err != nil {
        return err
    }
    if strings.TrimSpace(expected) != strings.TrimSpace(got) {
        return fmt.Errorf("expected %s got %s", expected, got)
    }
    return nil
}

func main() {
    if len(os.Args) < 2 {
        fmt.Println("usage: go run verifierG.go /path/to/binary")
        os.Exit(1)
    }
    bin := os.Args[1]
    if bin == "--" {
        if len(os.Args) < 3 {
            fmt.Println("usage: go run verifierG.go /path/to/binary")
            os.Exit(1)
        }
        bin = os.Args[2]
    }
    ref, err := buildRef()
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
    defer os.Remove(ref)
    cases := genCases()
    for i, c := range cases {
        if err := runCase(bin, ref, c); err != nil {
            fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, c.input)
            os.Exit(1)
        }
    }
    fmt.Println("All tests passed")
}


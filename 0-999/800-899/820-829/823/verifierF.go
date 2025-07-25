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
    ref := "./refF.bin"
    cmd := exec.Command("go", "build", "-o", ref, "823F.go")
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
    n := rng.Intn(6) + 1
    m := rng.Intn(n*(n-1)/2 + 1)
    edges := make(map[[2]int]struct{})
    var sb strings.Builder
    fmt.Fprintf(&sb, "1\n%d %d\n", n, m)
    for len(edges) < m {
        a := rng.Intn(n) + 1
        b := rng.Intn(n) + 1
        if a == b {
            continue
        }
        if a > b {
            a, b = b, a
        }
        key := [2]int{a, b}
        if _, ok := edges[key]; ok {
            continue
        }
        edges[key] = struct{}{}
        fmt.Fprintf(&sb, "%d %d\n", a, b)
    }
    return Case{sb.String()}
}

func genCases() []Case {
    rng := rand.New(rand.NewSource(8235))
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
        return fmt.Errorf("expected\n%s\ngot\n%s", expected, got)
    }
    return nil
}

func main() {
    if len(os.Args) < 2 {
        fmt.Println("usage: go run verifierF.go /path/to/binary")
        os.Exit(1)
    }
    bin := os.Args[1]
    if bin == "--" {
        if len(os.Args) < 3 {
            fmt.Println("usage: go run verifierF.go /path/to/binary")
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


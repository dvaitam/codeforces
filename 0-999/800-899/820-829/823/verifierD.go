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
    ref := "./refD.bin"
    cmd := exec.Command("go", "build", "-o", ref, "823D.go")
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

func randVK(rng *rand.Rand, n int) string {
    letters := "VK?"
    b := make([]byte, n)
    for i := range b {
        b[i] = letters[rng.Intn(3)]
    }
    return string(b)
}

func genCases() []Case {
    rng := rand.New(rand.NewSource(8233))
    cases := make([]Case, 100)
    for i := range cases {
        t := rng.Intn(3) + 1
        var sb strings.Builder
        fmt.Fprintf(&sb, "%d\n", t)
        for j := 0; j < t; j++ {
            n := rng.Intn(8) + 1
            s := randVK(rng, n)
            fmt.Fprintf(&sb, "%d\n%s\n", n, s)
        }
        cases[i] = Case{sb.String()}
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
        fmt.Println("usage: go run verifierD.go /path/to/binary")
        os.Exit(1)
    }
    bin := os.Args[1]
    if bin == "--" {
        if len(os.Args) < 3 {
            fmt.Println("usage: go run verifierD.go /path/to/binary")
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


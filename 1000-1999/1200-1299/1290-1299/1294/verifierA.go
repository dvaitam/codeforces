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

type Test struct {
    a, b, c, n int64
}

func expected(t Test) string {
    maxv := t.a
    if t.b > maxv {
        maxv = t.b
    }
    if t.c > maxv {
        maxv = t.c
    }
    need := (maxv - t.a) + (maxv - t.b) + (maxv - t.c)
    if t.n < need {
        return "NO"
    }
    if (t.n-need)%3 == 0 {
        return "YES"
    }
    return "NO"
}

func runProg(bin, input string) (string, error) {
    var cmd *exec.Cmd
    if strings.HasSuffix(bin, ".go") {
        cmd = exec.Command("go", "run", bin)
    } else {
        cmd = exec.Command(bin)
    }
    cmd.Stdin = strings.NewReader(input)
    var out bytes.Buffer
    var errBuf bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = &errBuf
    if err := cmd.Run(); err != nil {
        return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
    }
    return strings.TrimSpace(out.String()), nil
}

func genTest(rng *rand.Rand) Test {
    return Test{
        a: rng.Int63n(100000000) + 1,
        b: rng.Int63n(100000000) + 1,
        c: rng.Int63n(100000000) + 1,
        n: rng.Int63n(100000000) + 1,
    }
}

func (t Test) Input() string {
    return fmt.Sprintf("1\n%d %d %d %d\n", t.a, t.b, t.c, t.n)
}

func main() {
    if len(os.Args) != 2 {
        fmt.Println("Usage: go run verifierA.go /path/to/binary")
        os.Exit(1)
    }
    bin := os.Args[1]
    rng := rand.New(rand.NewSource(time.Now().UnixNano()))
    const cases = 100
    for i := 0; i < cases; i++ {
        tc := genTest(rng)
        exp := expected(tc)
        got, err := runProg(bin, tc.Input())
        if err != nil {
            fmt.Printf("case %d: %v\n", i+1, err)
            os.Exit(1)
        }
        if strings.ToUpper(strings.TrimSpace(got)) != exp {
            fmt.Printf("case %d failed\ninput:\n%sexpected: %s\ngot: %s\n", i+1, tc.Input(), exp, got)
            os.Exit(1)
        }
    }
    fmt.Printf("All %d tests passed\n", cases)
}


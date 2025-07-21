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
        return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
    }
    return strings.TrimSpace(out.String()), nil
}

func xorTo(n uint64) uint64 {
    switch n & 3 {
    case 0:
        return n
    case 1:
        return 1
    case 2:
        return n + 1
    }
    return 0
}

func expected(n int, arr [][2]uint64) string {
    var total uint64
    for _, v := range arr {
        x := v[0]
        m := v[1]
        start := x
        end := x + m - 1
        total ^= xorTo(end) ^ xorTo(start-1)
    }
    if total != 0 {
        return "tolik"
    }
    return "bolik"
}

func genCase(rng *rand.Rand) (string, string) {
    n := rng.Intn(5) + 1
    arr := make([][2]uint64, n)
    for i := 0; i < n; i++ {
        arr[i][0] = uint64(rng.Intn(100) + 1)
        arr[i][1] = uint64(rng.Intn(5) + 1)
    }
    var sb strings.Builder
    sb.WriteString(fmt.Sprintf("%d\n", n))
    for i := 0; i < n; i++ {
        sb.WriteString(fmt.Sprintf("%d %d\n", arr[i][0], arr[i][1]))
    }
    exp := expected(n, arr)
    return sb.String(), exp
}

func main() {
    if len(os.Args) != 2 {
        fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
        os.Exit(1)
    }
    bin := os.Args[1]
    rng := rand.New(rand.NewSource(time.Now().UnixNano()))
    for i := 0; i < 100; i++ {
        in, exp := genCase(rng)
        out, err := run(bin, in)
        if err != nil {
            fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
            os.Exit(1)
        }
        if out != exp {
            fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, out, in)
            os.Exit(1)
        }
    }
    fmt.Println("All tests passed")
}


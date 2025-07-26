package main

import (
    "bytes"
    "fmt"
    "math/bits"
    "math/rand"
    "os"
    "os/exec"
    "strconv"
    "strings"
    "time"
)

func runCandidate(bin, input string) (string, error) {
    var cmd *exec.Cmd
    if strings.HasSuffix(bin, ".go") {
        cmd = exec.Command("go", "run", bin)
    } else {
        cmd = exec.Command(bin)
    }
    cmd.Stdin = strings.NewReader(input)
    var out bytes.Buffer
    var stderr bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = &stderr
    if err := cmd.Run(); err != nil {
        return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
    }
    return strings.TrimSpace(out.String()), nil
}

func solve(x uint64) uint64 {
    return 1 << bits.OnesCount64(x)
}

func generateCase(rng *rand.Rand) []uint64 {
    t := rng.Intn(10) + 1
    arr := make([]uint64, t)
    for i := 0; i < t; i++ {
        arr[i] = uint64(rng.Int63() & ((1 << 30) - 1))
    }
    return arr
}

func main() {
    if len(os.Args) != 2 {
        fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
        os.Exit(1)
    }
    bin := os.Args[1]
    rng := rand.New(rand.NewSource(time.Now().UnixNano()))

    cases := [][]uint64{
        {0},
        {1},
        {3},
        {0, 1, (1 << 30) - 1},
    }
    for i := 0; i < 100; i++ {
        cases = append(cases, generateCase(rng))
    }

    for idx, arr := range cases {
        var sb strings.Builder
        sb.WriteString(fmt.Sprintf("%d\n", len(arr)))
        for _, v := range arr {
            sb.WriteString(fmt.Sprintf("%d\n", v))
        }
        input := sb.String()
        out, err := runCandidate(bin, input)
        if err != nil {
            fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", idx+1, err, input)
            os.Exit(1)
        }
        fields := strings.Fields(out)
        if len(fields) != len(arr) {
            fmt.Fprintf(os.Stderr, "case %d failed: expected %d numbers got %d\ninput:\n%s", idx+1, len(arr), len(fields), input)
            os.Exit(1)
        }
        for j, f := range fields {
            val, err := strconv.ParseUint(f, 10, 64)
            if err != nil {
                fmt.Fprintf(os.Stderr, "case %d failed: bad number %q\n", idx+1, f)
                os.Exit(1)
            }
            expected := solve(arr[j])
            if val != expected {
                fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %d\ninput:\n%s", idx+1, expected, val, input)
                os.Exit(1)
            }
        }
    }
    fmt.Println("All tests passed")
}


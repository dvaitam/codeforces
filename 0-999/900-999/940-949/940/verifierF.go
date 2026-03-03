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
    exe := "oracleF"
    cmd := exec.Command("go", "build", "-o", exe, "940F.go")
    if out, err := cmd.CombinedOutput(); err != nil {
        return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
    }
    return "./" + exe, nil
}

func runProgram(bin, input string) (string, error) {
    var cmd *exec.Cmd
    if strings.HasSuffix(bin, ".go") {
        cmd = exec.Command("go", "run", bin)
    } else {
        cmd = exec.Command(bin)
    }
    cmd.Stdin = strings.NewReader(input)
    var out bytes.Buffer
    var errb bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = &errb
    if err := cmd.Run(); err != nil {
        return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
    }
    return strings.TrimSpace(out.String()), nil
}

func generateCase(r *rand.Rand) string {
	n := r.Intn(10) + 1
	m := r.Intn(10) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for i := 0; i < n; i++ {
		sb.WriteString(fmt.Sprintf("%d ", r.Intn(20)+1))
	}
	sb.WriteString("\n")
	for i := 0; i < m; i++ {
		op := r.Intn(2) + 1
		if op == 1 {
			l := r.Intn(n) + 1
			r_ := l + r.Intn(n-l+1)
			sb.WriteString(fmt.Sprintf("%d %d %d\n", op, l, r_))
		} else {
			pos := r.Intn(n) + 1
			val := r.Intn(20) + 1
			sb.WriteString(fmt.Sprintf("%d %d %d\n", op, pos, val))
		}
	}
	return sb.String()
}

func main() {
    if len(os.Args) != 2 {
        fmt.Println("usage: go run verifierF.go /path/to/binary")
        os.Exit(1)
    }
    bin := os.Args[1]
    oracle, err := buildOracle()
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
    defer os.Remove(oracle)

    r := rand.New(rand.NewSource(time.Now().UnixNano()))
    for i := 1; i <= 100; i++ {
        input := generateCase(r)
        exp, err := runProgram(oracle, input)
        if err != nil {
            fmt.Fprintf(os.Stderr, "oracle error on case %d: %v\n", i, err)
			fmt.Fprintf(os.Stderr, "Input:\n%s\n", input)
            os.Exit(1)
        }
        got, err := runProgram(bin, input)
        if err != nil {
            fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i, err)
			fmt.Fprintf(os.Stderr, "Input:\n%s\n", input)
            os.Exit(1)
        }
        if got != exp {
            fmt.Fprintf(os.Stderr, "case %d failed: expected %q got %q\n", i, exp, got)
			fmt.Fprintf(os.Stderr, "Input:\n%s\n", input)
            os.Exit(1)
        }
    }
    fmt.Println("All 100 tests passed")
}

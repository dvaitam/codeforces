package main

import (
    "bytes"
    "fmt"
    "math/rand"
    "os"
    "os/exec"
    "strconv"
    "strings"
    "time"
)

func solve(A int) string {
	if A == 1 {
		return "1 1\n1"
	}
	N := 10 * (A - 1)
	return fmt.Sprintf("%d 2\n1 10", N)
}

func run(bin, input string) (string, error) {
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

func main() {
	if len(os.Args) == 3 && os.Args[1] == "--" {
		os.Args = append([]string{os.Args[0]}, os.Args[2])
	}
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		A := rng.Intn(1000) + 1
		input := fmt.Sprintf("%d\n", A)
        expected := solve(A)
        got, err := run(bin, input)
        if err != nil {
            fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
            os.Exit(1)
        }
        if strings.TrimSpace(got) != strings.TrimSpace(expected) {
            // Try to accept other valid constructions, e.g., M=2 with denominations {1, d}
            if !accept865Variant(A, got) {
                fmt.Fprintf(os.Stderr, "case %d failed\ninput:\n%s\nexpected:\n%s\n---\ngot:\n%s\n", i+1, input, expected, got)
                os.Exit(1)
            }
        }
    }
    fmt.Println("All tests passed")
}

func accept865Variant(A int, out string) bool {
    lines := strings.Split(strings.TrimSpace(out), "\n")
    if len(lines) < 2 {
        return false
    }
    first := strings.Fields(lines[0])
    if len(first) != 2 {
        return false
    }
    N, err1 := strconv.Atoi(first[0])
    M, err2 := strconv.Atoi(first[1])
    if err1 != nil || err2 != nil || N < 1 || N > 1_000_000 || M < 1 || M > 10 {
        return false
    }
    denFields := strings.Fields(lines[1])
    if len(denFields) != M {
        return false
    }
    dens := make([]int, M)
    seen := map[int]bool{}
    for i := 0; i < M; i++ {
        v, err := strconv.Atoi(denFields[i])
        if err != nil || v < 1 || v > 1_000_000 || seen[v] {
            return false
        }
        seen[v] = true
        dens[i] = v
    }
    // Special case A==1: allow "1 1" then single denom
    if A == 1 {
        if M == 1 {
            return true
        }
    }
    // Accept the common pattern: M==2 and denominations are {1, d}
    if M == 2 {
        hasOne := false
        other := 0
        for _, v := range dens {
            if v == 1 {
                hasOne = true
            } else {
                other = v
            }
        }
        if hasOne && other > 0 {
            ways := N/other + 1
            return ways == A
        }
    }
    return false
}

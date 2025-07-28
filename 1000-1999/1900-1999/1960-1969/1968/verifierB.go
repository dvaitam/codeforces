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

func solve(a, b string) int {
	i, j := 0, 0
	n, m := len(a), len(b)
	for i < n && j < m {
		if a[i] == b[j] {
			i++
		}
		j++
	}
	return i
}

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(10) + 1
	m := rng.Intn(10) + 1
	var sb strings.Builder
	aBytes := make([]byte, n)
	bBytes := make([]byte, m)
	for i := 0; i < n; i++ {
		if rng.Intn(2) == 0 {
			aBytes[i] = '0'
		} else {
			aBytes[i] = '1'
		}
	}
	for i := 0; i < m; i++ {
		if rng.Intn(2) == 0 {
			bBytes[i] = '0'
		} else {
			bBytes[i] = '1'
		}
	}
	aStr := string(aBytes)
	bStr := string(bBytes)
	sb.WriteString("1\n")
	fmt.Fprintf(&sb, "%d %d\n%s\n%s\n", n, m, aStr, bStr)
	expect := fmt.Sprint(solve(aStr, bStr))
	return sb.String(), expect
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
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	const cases = 100
	for i := 0; i < cases; i++ {
		inp, exp := genCase(rng)
		got, err := run(bin, inp)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "case %d failed\nexpected: %s\ngot: %s\ninput:\n%s", i+1, exp, got, inp)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", cases)
}

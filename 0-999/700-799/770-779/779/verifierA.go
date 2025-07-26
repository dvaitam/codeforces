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

func runCandidate(bin, input string) (string, error) {
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

func solveCaseA(n int, a, b []int) string {
	freqA := make([]int, 6)
	freqB := make([]int, 6)
	for _, v := range a {
		freqA[v]++
	}
	for _, v := range b {
		freqB[v]++
	}
	sumDiff := 0
	for i := 1; i <= 5; i++ {
		if (freqA[i]+freqB[i])%2 != 0 {
			return "-1"
		}
		d := freqA[i] - freqB[i]
		if d < 0 {
			d = -d
		}
		sumDiff += d
	}
	return fmt.Sprintf("%d", sumDiff/4)
}

func generateCaseA(rng *rand.Rand) (string, string) {
	n := rng.Intn(100) + 1
	a := make([]int, n)
	b := make([]int, n)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		a[i] = rng.Intn(5) + 1
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(a[i]))
	}
	sb.WriteByte('\n')
	for i := 0; i < n; i++ {
		b[i] = rng.Intn(5) + 1
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(b[i]))
	}
	sb.WriteByte('\n')
	input := sb.String()
	expect := solveCaseA(n, a, b)
	return input, expect
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
		in, exp := generateCaseA(rng)
		out, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(exp) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

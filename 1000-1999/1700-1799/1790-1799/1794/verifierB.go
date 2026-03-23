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

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
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

func buildCandidate(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp, err := os.CreateTemp("", "cand*")
		if err != nil {
			return "", nil, err
		}
		tmp.Close()
		cmd := exec.Command("go", "build", "-o", tmp.Name(), path)
		if out, err := cmd.CombinedOutput(); err != nil {
			os.Remove(tmp.Name())
			return "", nil, fmt.Errorf("build failed: %v\n%s", err, out)
		}
		return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
	}
	return path, func() {}, nil
}

func solveCase(n int, a []int) string {
	vals := make([]int64, n)
	for i := 0; i < n; i++ {
		vals[i] = int64(a[i])
		if vals[i] == 1 {
			vals[i]++
		}
	}
	for i := 0; i < n-1; i++ {
		if vals[i+1]%vals[i] == 0 {
			vals[i+1]++
		}
	}
	parts := make([]string, n)
	for i := 0; i < n; i++ {
		parts[i] = strconv.FormatInt(vals[i], 10)
	}
	return strings.Join(parts, " ")
}

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(20) + 1
	sb := strings.Builder{}
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", rng.Intn(50)+1))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func parseInput(input string) (int, []int) {
	tokens := strings.Fields(input)
	// tokens[0] = "1" (t), tokens[1] = n, tokens[2..] = array
	n, _ := strconv.Atoi(tokens[1])
	a := make([]int, n)
	for i := 0; i < n; i++ {
		a[i], _ = strconv.Atoi(tokens[2+i])
	}
	return n, a
}

func validateOutput(n int, a []int, output string) error {
	tokens := strings.Fields(output)
	if len(tokens) != n {
		return fmt.Errorf("expected %d values, got %d", n, len(tokens))
	}
	vals := make([]int64, n)
	totalOps := int64(0)
	for i := 0; i < n; i++ {
		v, err := strconv.ParseInt(tokens[i], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid number %q", tokens[i])
		}
		if v < int64(a[i]) {
			return fmt.Errorf("value %d decreased from %d", v, a[i])
		}
		totalOps += v - int64(a[i])
		vals[i] = v
	}
	if totalOps > int64(2*n) {
		return fmt.Errorf("too many operations: %d > %d", totalOps, 2*n)
	}
	for i := 0; i < n-1; i++ {
		if vals[i+1]%vals[i] == 0 {
			return fmt.Errorf("a[%d]=%d divides a[%d]=%d", i, vals[i], i+1, vals[i+1])
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	candidatePath, cleanup, err := buildCandidate(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 1; i <= 100; i++ {
		input := generateCase(rng)
		n, a := parseInput(input)
		got, err := run(candidatePath, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", i, err)
			os.Exit(1)
		}
		if err := validateOutput(n, a, got); err != nil {
			// Also check against embedded solver
			expect := solveCase(n, a)
			fmt.Printf("case %d failed: %v\ninput:\n%s\nexpected:\n%s\n\ngot:\n%s\n", i, err, input, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

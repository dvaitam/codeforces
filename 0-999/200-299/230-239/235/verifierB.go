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

func solveCase(p []float64) float64 {
	m := 0.0
	r := 0.0
	for _, prob := range p {
		r += prob * (2*m + 1)
		m = (m + 1) * prob
	}
	return r
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(10) + 1
	probs := make([]float64, n)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d ", n))
	for i := 0; i < n; i++ {
		probs[i] = rng.Float64()
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%.2f", probs[i]))
	}
	input := sb.String() + "\n"
	expect := fmt.Sprintf("%.10f", solveCase(probs))
	return input, expect
}

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

func parseProbs(parts []string) ([]float64, error) {
	res := make([]float64, len(parts))
	for i, s := range parts {
		v, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return nil, err
		}
		res[i] = v
	}
	return res, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		out, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, in)
			os.Exit(1)
		}
		if out != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

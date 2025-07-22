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
	return out.String(), nil
}

func expectedOutput(w [][]int) string {
	n := len(w)
	var sb strings.Builder
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			v := float64(w[i][j]+w[j][i]) / 2.0
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%.8f", v))
		}
		sb.WriteByte('\n')
	}
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			v := float64(w[i][j]-w[j][i]) / 2.0
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%.8f", v))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(5) + 1
	w := make([][]int, n)
	for i := 0; i < n; i++ {
		w[i] = make([]int, n)
		for j := 0; j < n; j++ {
			w[i][j] = rng.Intn(2001) - 1000
		}
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", w[i][j])
		}
		sb.WriteByte('\n')
	}
	return sb.String(), expectedOutput(w)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, expect := genCase(rng)
		out, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(expect) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected\n%s\n got\n%s\ninput:\n%s", i+1, expect, out, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

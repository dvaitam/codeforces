package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func buildOracle() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	oracle := filepath.Join(dir, "oracleD")
	cmd := exec.Command("go", "build", "-o", oracle, "386D.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return oracle, nil
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

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(6) + 3 // 3..8
	pos := rng.Perm(n)[:3]
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	fmt.Fprintf(&sb, "%d %d %d\n", pos[0]+1, pos[1]+1, pos[2]+1)
	matrix := make([][]byte, n)
	for i := 0; i < n; i++ {
		matrix[i] = make([]byte, n)
		for j := 0; j < n; j++ {
			if i == j {
				matrix[i][j] = '*'
			} else if j < i {
				matrix[i][j] = matrix[j][i]
			} else {
				matrix[i][j] = byte(rng.Intn(26)) + 'a'
			}
		}
	}
	for i := 0; i < n; i++ {
		sb.Write(matrix[i])
		sb.WriteByte('\n')
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	cases := []string{}
	// simple deterministic case
	cases = append(cases, func() string {
		n := 3
		var sb strings.Builder
		sb.WriteString("3\n1 2 3\n")
		sb.WriteString("*aa\n")
		sb.WriteString("a*a\n")
		sb.WriteString("aa*\n")
		return sb.String()
	}())

	for i := 0; i < 100; i++ {
		cases = append(cases, generateCase(rng))
	}

	for i, in := range cases {
		// run oracle
		cmdO := exec.Command(oracle)
		cmdO.Stdin = strings.NewReader(in)
		var outO bytes.Buffer
		cmdO.Stdout = &outO
		if err := cmdO.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "oracle run error: %v\n", err)
			os.Exit(1)
		}
		expected := strings.TrimSpace(outO.String())

		got, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expected, got, in)
			os.Exit(1)
		}
	}

	fmt.Println("All tests passed")
}

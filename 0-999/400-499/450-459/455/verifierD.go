package main

import (
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
	cmd := exec.Command("go", "build", "-o", oracle, "455D.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return oracle, nil
}

func runBin(bin string, input string) (string, error) {
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
	n := rng.Intn(8) + 1
	arr := make([]int, n)
	for i := range arr {
		arr[i] = rng.Intn(n) + 1
	}
	q := rng.Intn(8) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	sb.WriteString(fmt.Sprintf("%d\n", q))
	for i := 0; i < q; i++ {
		if rng.Intn(2) == 0 {
			l := rng.Intn(n) + 1
			r := rng.Intn(n) + 1
			sb.WriteString(fmt.Sprintf("1 %d %d\n", l, r))
		} else {
			l := rng.Intn(n) + 1
			r := rng.Intn(n) + 1
			k := rng.Intn(n) + 1
			sb.WriteString(fmt.Sprintf("2 %d %d %d\n", l, r, k))
		}
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
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

	cases := make([]string, 0, 101)
	for i := 0; i < 100; i++ {
		cases = append(cases, generateCase(rng))
	}

	for idx, input := range cases {
		want, err := runBin(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle error on case %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		got, err := runBin(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != want {
			fmt.Fprintf(os.Stderr, "case %d failed\nexpected:\n%s\n\ngot:\n%s\ninput:\n%s", idx+1, want, got, input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}

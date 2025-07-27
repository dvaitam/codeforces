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
	oracle := filepath.Join(dir, "oracleC")
	cmd := exec.Command("go", "build", "-o", oracle, "1387C.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return oracle, nil
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
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func ruleLine(rng *rand.Rand, a, G int) string {
	k := rng.Intn(3) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d", a, k)
	for j := 0; j < k; j++ {
		fmt.Fprintf(&sb, " %d", rng.Intn(G))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func genCase(rng *rand.Rand) string {
	G := rng.Intn(4) + 3 // 3..6
	minN := G - 2
	N := minN + rng.Intn(3)
	M := rng.Intn(3)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d\n", G, N, M)
	var rules []string
	for a := 2; a < G; a++ {
		rules = append(rules, ruleLine(rng, a, G))
	}
	for len(rules) < N {
		a := rng.Intn(G-2) + 2
		rules = append(rules, ruleLine(rng, a, G))
	}
	for _, r := range rules {
		sb.WriteString(r)
	}
	for i := 0; i < M; i++ {
		l := rng.Intn(3) + 1
		fmt.Fprintf(&sb, "%d", l)
		for j := 0; j < l; j++ {
			fmt.Fprintf(&sb, " %d", rng.Intn(2))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	const cases = 100
	for i := 1; i <= cases; i++ {
		input := genCase(rng)
		expect, err := run(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle error on case %d: %v\n", i, err)
			os.Exit(1)
		}
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i, expect, got, input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", cases)
}

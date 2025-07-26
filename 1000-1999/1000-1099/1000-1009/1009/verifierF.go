package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func buildOracle() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	oracle := filepath.Join(dir, "oracleF")
	cmd := exec.Command("go", "build", "-o", oracle, "1009F.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return oracle, nil
}

func runProg(bin, input string) (string, error) {
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

func genCaseF(rng *rand.Rand) (int, [][2]int) {
	n := rng.Intn(15) + 1
	if n == 1 {
		return n, nil
	}
	edges := make([][2]int, n-1)
	for i := 2; i <= n; i++ {
		p := rng.Intn(i-1) + 1
		edges[i-2] = [2]int{p, i}
	}
	return n, edges
}

func runCaseF(bin, oracle string, n int, edges [][2]int) error {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(n))
	sb.WriteByte('\n')
	for i, e := range edges {
		if i > 0 {
			sb.WriteByte('\n')
		}
		sb.WriteString(strconv.Itoa(e[0]))
		sb.WriteByte(' ')
		sb.WriteString(strconv.Itoa(e[1]))
	}
	if n > 1 {
		sb.WriteByte('\n')
	}
	input := sb.String()
	exp, err := runProg(oracle, input)
	if err != nil {
		return fmt.Errorf("oracle error: %v", err)
	}
	got, err := runProg(bin, input)
	if err != nil {
		return err
	}
	if got != exp {
		return fmt.Errorf("expected %s got %s", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
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
	for i := 0; i < 100; i++ {
		n, edges := genCaseF(rng)
		if err := runCaseF(bin, oracle, n, edges); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

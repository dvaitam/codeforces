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

func buildOracle() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	oracle := dir + "/oracleD"
	cmd := exec.Command("go", "build", "-o", oracle, "468D.go")
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

func generateCase(rng *rand.Rand) (int, [][3]int) {
	n := rng.Intn(6) + 2
	edges := make([][3]int, 0, n-1)
	for i := 2; i <= n; i++ {
		u := rng.Intn(i-1) + 1
		v := i
		w := rng.Intn(10) + 1
		edges = append(edges, [3]int{u, v, w})
	}
	return n, edges
}

func buildInput(n int, edges [][3]int) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for _, e := range edges {
		sb.WriteString(fmt.Sprintf("%d %d %d\n", e[0], e[1], e[2]))
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
	for i := 0; i < 100; i++ {
		n, edges := generateCase(rng)
		input := buildInput(n, edges)
		exp, err := run(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle error: %v\n", err)
			os.Exit(1)
		}
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if out != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, out, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

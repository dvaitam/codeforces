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
	oracle := dir + "/oracleE"
	cmd := exec.Command("go", "build", "-o", oracle, "468E.go")
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

func generateCase(rng *rand.Rand) (int, int, [][3]int64) {
	n := rng.Intn(6) + 1
	k := rng.Intn(5)
	edges := make([][3]int64, k)
	positions := make(map[[2]int]bool)
	for i := 0; i < k; i++ {
		for {
			x := rng.Intn(n) + 1
			y := rng.Intn(n) + 1
			if !positions[[2]int{x, y}] {
				positions[[2]int{x, y}] = true
				w := rng.Int63n(5)
				edges[i] = [3]int64{int64(x), int64(y), w}
				break
			}
		}
	}
	return n, k, edges
}

func buildInput(n, k int, edges [][3]int64) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, k))
	for _, e := range edges {
		sb.WriteString(fmt.Sprintf("%d %d %d\n", e[0], e[1], e[2]))
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
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
		n, k, edges := generateCase(rng)
		input := buildInput(n, k, edges)
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

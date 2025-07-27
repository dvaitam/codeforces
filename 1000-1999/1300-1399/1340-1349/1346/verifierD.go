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
	oracle := "oracleD"
	cmd := exec.Command("go", "build", "-o", oracle, "1346D.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle: %v\n%s", err, out)
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
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func genCase(rng *rand.Rand) string {
	t := rng.Intn(3) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", t)
	for i := 0; i < t; i++ {
		n := rng.Intn(4) + 2
		maxEdges := n * (n - 1) / 2
		m := rng.Intn(maxEdges-(n-1)+1) + (n - 1)
		fmt.Fprintf(&sb, "%d %d\n", n, m)
		edges := make(map[[2]int]bool)
		// ensure connectivity via chain
		for j := 1; j < n; j++ {
			u := j
			v := j + 1
			w := rng.Intn(100) + 1
			fmt.Fprintf(&sb, "%d %d %d\n", u, v, w)
			edges[[2]int{u, v}] = true
			edges[[2]int{v, u}] = true
		}
		extra := m - (n - 1)
		for extra > 0 {
			u := rng.Intn(n) + 1
			v := rng.Intn(n) + 1
			if u == v || edges[[2]int{u, v}] {
				continue
			}
			w := rng.Intn(100) + 1
			fmt.Fprintf(&sb, "%d %d %d\n", u, v, w)
			edges[[2]int{u, v}] = true
			edges[[2]int{v, u}] = true
			extra--
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

	for i := 0; i < 100; i++ {
		input := genCase(rng)
		exp, err := run(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failure on case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "case %d failed:\ninput:\n%sexpected:\n%s\ngot:\n%s\n", i+1, input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

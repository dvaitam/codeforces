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
	oracle := filepath.Join(dir, "oracleI")
	cmd := exec.Command("go", "build", "-o", oracle, "1533I.go")
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
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func generateCase(rng *rand.Rand) string {
	n1 := rng.Intn(3) + 1
	n2 := rng.Intn(3) + 1
	m := rng.Intn(n1*n2-n1-n2+2) + max(n1, n2)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d\n", n1, n2, m)
	for i := 0; i < n1; i++ {
		fmt.Fprintf(&sb, "%d ", rng.Intn(10)+1)
	}
	sb.WriteByte('\n')
	type pair struct{ x, y int }
	used := map[pair]bool{}
	ensure := []pair{}
	for i := 1; i <= n1; i++ {
		y := rng.Intn(n2) + 1
		used[pair{i, y}] = true
		ensure = append(ensure, pair{i, y})
	}
	for j := 1; j <= n2; j++ {
		x := rng.Intn(n1) + 1
		p := pair{x, j}
		if !used[p] {
			used[p] = true
			ensure = append(ensure, p)
		}
	}
	for len(used) < m {
		p := pair{rng.Intn(n1) + 1, rng.Intn(n2) + 1}
		used[p] = true
	}
	cnt := 0
	for p := range used {
		fmt.Fprintf(&sb, "%d %d\n", p.x, p.y)
		cnt++
		if cnt == m {
			break
		}
	}
	return sb.String()
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func main() {
	if len(os.Args) != 2 && !(len(os.Args) == 3 && os.Args[1] == "--") {
		fmt.Println("usage: go run verifierI.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	if bin == "--" {
		bin = os.Args[2]
	}
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in := generateCase(rng)
		exp, err := run(oracle, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle error on case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(exp) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, got, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

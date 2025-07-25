package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

func runBinary(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func buildOracle() (string, error) {
	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)
	exe := filepath.Join(dir, "oracleD")
	src := filepath.Join(dir, "778D.go")
	cmd := exec.Command("go", "build", "-o", exe, src)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle: %v\n%s", err, out)
	}
	return exe, nil
}

func genGrid(rng *rand.Rand, n, m int) []string {
	g := make([][]byte, n)
	for i := range g {
		g[i] = make([]byte, m)
	}
	for i := 0; i < n; i += 2 {
		for j := 0; j < m; j += 2 {
			if rng.Intn(2) == 0 {
				g[i][j], g[i+1][j] = 'U', 'D'
				g[i][j+1], g[i+1][j+1] = 'U', 'D'
			} else {
				g[i][j], g[i][j+1] = 'L', 'R'
				g[i+1][j], g[i+1][j+1] = 'L', 'R'
			}
		}
	}
	res := make([]string, n)
	for i := range g {
		res[i] = string(g[i])
	}
	return res
}

func genCase(rng *rand.Rand) string {
	n := rng.Intn(3)*2 + 2
	m := rng.Intn(3)*2 + 2
	g1 := genGrid(rng, n, m)
	g2 := genGrid(rng, n, m)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for i := 0; i < n; i++ {
		sb.WriteString(g1[i])
		sb.WriteByte('\n')
	}
	for i := 0; i < n; i++ {
		sb.WriteString(g2[i])
		sb.WriteByte('\n')
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for t := 0; t < 100; t++ {
		input := genCase(rng)
		candOut, cErr := runBinary(candidate, input)
		refOut, rErr := runBinary(oracle, input)
		if cErr != nil {
			fmt.Fprintf(os.Stderr, "case %d: candidate error: %v\ninput:\n%s", t+1, cErr, input)
			os.Exit(1)
		}
		if rErr != nil {
			fmt.Fprintf(os.Stderr, "case %d: oracle error: %v\ninput:\n%s", t+1, rErr, input)
			os.Exit(1)
		}
		if candOut != refOut {
			fmt.Fprintf(os.Stderr, "case %d failed\ninput:\n%s\nexpected:%s\nactual:%s\n", t+1, input, refOut, candOut)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

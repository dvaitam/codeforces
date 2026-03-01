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
)

func buildOracle() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	oracle := filepath.Join(dir, "oracleE")
	cmd := exec.Command("go", "build", "-o", oracle, "1264E.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return oracle, nil
}

func runProg(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func parseCase(input string) (int, [][2]int, error) {
	fields := strings.Fields(input)
	if len(fields) < 2 {
		return 0, nil, fmt.Errorf("bad input header")
	}
	n, err := strconv.Atoi(fields[0])
	if err != nil {
		return 0, nil, fmt.Errorf("bad n: %w", err)
	}
	m, err := strconv.Atoi(fields[1])
	if err != nil {
		return 0, nil, fmt.Errorf("bad m: %w", err)
	}
	need := 2 + 2*m
	if len(fields) != need {
		return 0, nil, fmt.Errorf("bad input length")
	}
	edges := make([][2]int, 0, m)
	for i := 0; i < m; i++ {
		u, err := strconv.Atoi(fields[2+2*i])
		if err != nil {
			return 0, nil, fmt.Errorf("bad u in edge %d: %w", i, err)
		}
		v, err := strconv.Atoi(fields[3+2*i])
		if err != nil {
			return 0, nil, fmt.Errorf("bad v in edge %d: %w", i, err)
		}
		edges = append(edges, [2]int{u - 1, v - 1})
	}
	return n, edges, nil
}

func parseMatrix(out string, n int) ([][]int, error) {
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != n {
		return nil, fmt.Errorf("expected %d lines, got %d", n, len(lines))
	}
	a := make([][]int, n)
	for i := 0; i < n; i++ {
		line := strings.TrimSpace(lines[i])
		if len(line) != n {
			return nil, fmt.Errorf("line %d has length %d, expected %d", i+1, len(line), n)
		}
		a[i] = make([]int, n)
		for j := 0; j < n; j++ {
			if line[j] != '0' && line[j] != '1' {
				return nil, fmt.Errorf("line %d contains non-binary char", i+1)
			}
			a[i][j] = int(line[j] - '0')
		}
	}
	return a, nil
}

func validateMatrix(a [][]int, fixed [][2]int) error {
	n := len(a)
	for i := 0; i < n; i++ {
		if a[i][i] != 0 {
			return fmt.Errorf("a[%d][%d] must be 0", i+1, i+1)
		}
		for j := i + 1; j < n; j++ {
			if a[i][j]+a[j][i] != 1 {
				return fmt.Errorf("a[%d][%d] + a[%d][%d] must be 1", i+1, j+1, j+1, i+1)
			}
		}
	}
	for _, e := range fixed {
		u, v := e[0], e[1]
		if a[u][v] != 1 || a[v][u] != 0 {
			return fmt.Errorf("fixed match %d -> %d was changed", u+1, v+1)
		}
	}
	return nil
}

func beauty(a [][]int) int {
	n := len(a)
	ans := 0
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if i == j || a[i][j] == 0 {
				continue
			}
			for k := 0; k < n; k++ {
				if k == i || k == j {
					continue
				}
				if a[j][k] == 1 && a[k][i] == 1 {
					ans++
				}
			}
		}
	}
	return ans
}

func genCase(rng *rand.Rand) string {
	n := rng.Intn(3) + 2 // 2..4
	maxEdges := n * (n - 1) / 2
	m := rng.Intn(maxEdges + 1)
	edges := make(map[[2]int]bool)
	used := make(map[[2]int]bool)
	for len(edges) < m {
		x := rng.Intn(n) + 1
		y := rng.Intn(n-1) + 1
		if y >= x {
			y++
		}
		key := [2]int{x, y}
		if x > y {
			key = [2]int{y, x}
		}
		if used[key] {
			continue
		}
		used[key] = true
		edges[[2]int{x, y}] = true
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n%d\n", n, m)
	for e := range edges {
		fmt.Fprintf(&sb, "%d %d\n", e[0], e[1])
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
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	rng := rand.New(rand.NewSource(1264))
	for i := 0; i < 100; i++ {
		input := genCase(rng)
		want, err := runProg(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle error on case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := runProg(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}

		n, fixed, err := parseCase(input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "verifier parse error on case %d: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		wantMatrix, err := parseMatrix(want, n)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle produced invalid matrix on case %d: %v\noutput:\n%s\ninput:\n%s", i+1, err, want, input)
			os.Exit(1)
		}
		gotMatrix, err := parseMatrix(got, n)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: invalid output: %v\noutput:\n%s\ninput:\n%s", i+1, err, got, input)
			os.Exit(1)
		}
		if err := validateMatrix(wantMatrix, fixed); err != nil {
			fmt.Fprintf(os.Stderr, "oracle output invalid on case %d: %v\noutput:\n%s\ninput:\n%s", i+1, err, want, input)
			os.Exit(1)
		}
		if err := validateMatrix(gotMatrix, fixed); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\noutput:\n%s\ninput:\n%s", i+1, err, got, input)
			os.Exit(1)
		}
		wantBeauty := beauty(wantMatrix)
		gotBeauty := beauty(gotMatrix)
		if gotBeauty != wantBeauty {
			fmt.Fprintf(os.Stderr, "case %d failed: expected beauty %d, got %d\ninput:\n%s", i+1, wantBeauty, gotBeauty, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

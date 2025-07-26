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
)

func buildOracle() (string, error) {
	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)
	src := filepath.Join(dir, "1025E.go")
	bin := filepath.Join(os.TempDir(), "oracle1025E.bin")
	cmd := exec.Command("go", "build", "-o", bin, src)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return bin, nil
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

func genCase(r *rand.Rand) string {
	n := r.Intn(4) + 2
	m := r.Intn(n) + 1
	type cell struct{ x, y int }
	used := make(map[[2]int]bool)
	cells := make([]cell, 0, n*n)
	for x := 1; x <= n; x++ {
		for y := 1; y <= n; y++ {
			cells = append(cells, cell{x, y})
		}
	}
	r.Shuffle(len(cells), func(i, j int) { cells[i], cells[j] = cells[j], cells[i] })
	start := make([]cell, m)
	target := make([]cell, m)
	idx := 0
	for i := 0; i < m; i++ {
		start[i] = cells[idx]
		used[[2]int{cells[idx].x, cells[idx].y}] = true
		idx++
	}
	for i := 0; i < m; i++ {
		// choose unused cell for target
		for used[[2]int{cells[idx].x, cells[idx].y}] {
			idx++
		}
		target[i] = cells[idx]
		used[[2]int{cells[idx].x, cells[idx].y}] = true
		idx++
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, m)
	for _, c := range start {
		fmt.Fprintf(&sb, "%d %d\n", c.x, c.y)
	}
	for _, c := range target {
		fmt.Fprintf(&sb, "%d %d\n", c.x, c.y)
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	userBin := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(oracle)
	r := rand.New(rand.NewSource(1))
	for i := 0; i < 100; i++ {
		input := genCase(r)
		want, err := run(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failed on test %d: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		got, err := run(userBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if want != got {
			fmt.Printf("test %d failed\ninput:\n%sexpected: %s\ngot: %s\n", i+1, input, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", 100)
}

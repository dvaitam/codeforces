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
	oracle := filepath.Join(dir, "oracleF1")
	cmd := exec.Command("go", "build", "-o", oracle, "513F1.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return oracle, nil
}

func run(bin, input string) (string, error) {
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

func randomCell(rng *rand.Rand, cells [][2]int) (int, int) {
	c := cells[rng.Intn(len(cells))]
	return c[0], c[1]
}

func genCase(rng *rand.Rand) string {
	n := rng.Intn(4) + 2
	m := rng.Intn(4) + 2
	grid := make([][]byte, n)
	cells := make([][2]int, 0)
	for i := 0; i < n; i++ {
		row := make([]byte, m)
		for j := 0; j < m; j++ {
			if rng.Intn(4) == 0 {
				row[j] = '#'
			} else {
				row[j] = '.'
				cells = append(cells, [2]int{i, j})
			}
		}
		grid[i] = row
	}
	if len(cells) == 0 {
		// ensure at least one free cell
		grid[0][0] = '.'
		cells = append(cells, [2]int{0, 0})
	}
	males := rng.Intn(3)
	females := rng.Intn(3)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d %d\n", n, m, males, females)
	for i := 0; i < n; i++ {
		sb.WriteString(string(grid[i]))
		sb.WriteByte('\n')
	}
	br, bc := randomCell(rng, cells)
	bt := rng.Intn(3) + 1
	fmt.Fprintf(&sb, "%d %d %d\n", br+1, bc+1, bt)
	for i := 0; i < males; i++ {
		r, c := randomCell(rng, cells)
		t := rng.Intn(3) + 1
		fmt.Fprintf(&sb, "%d %d %d\n", r+1, c+1, t)
	}
	for i := 0; i < females; i++ {
		r, c := randomCell(rng, cells)
		t := rng.Intn(3) + 1
		fmt.Fprintf(&sb, "%d %d %d\n", r+1, c+1, t)
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF1.go /path/to/binary")
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
			fmt.Fprintf(os.Stderr, "case %d: %v\n", i, err)
			os.Exit(1)
		}
		if got != expect {
			fmt.Printf("case %d failed\nexpected:\n%s\n\ngot:\n%s\n", i, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", cases)
}

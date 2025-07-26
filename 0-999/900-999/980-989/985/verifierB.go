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

func expectedB(n, m int, grid [][]byte) string {
	cnt := make([]int, m)
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if grid[i][j] == '1' {
				cnt[j]++
			}
		}
	}
	for i := 0; i < n; i++ {
		ok := true
		for j := 0; j < m; j++ {
			if cnt[j]-int(grid[i][j]-'0') == 0 {
				ok = false
				break
			}
		}
		if ok {
			return "YES"
		}
	}
	return "NO"
}

func generateCaseB(rng *rand.Rand) (int, int, [][]byte) {
	n := rng.Intn(5) + 1
	m := rng.Intn(5) + 1
	grid := make([][]byte, n)
	for i := 0; i < n; i++ {
		row := make([]byte, m)
		for j := 0; j < m; j++ {
			if rng.Intn(2) == 0 {
				row[j] = '0'
			} else {
				row[j] = '1'
			}
		}
		grid[i] = row
	}
	// ensure each column has at least one '1'
	for j := 0; j < m; j++ {
		has := false
		for i := 0; i < n; i++ {
			if grid[i][j] == '1' {
				has = true
				break
			}
		}
		if !has {
			grid[rng.Intn(n)][j] = '1'
		}
	}
	return n, m, grid
}

func runCaseB(bin string, n, m int, grid [][]byte) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for i := 0; i < n; i++ {
		sb.WriteString(string(grid[i]))
		sb.WriteByte('\n')
	}
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	expected := expectedB(n, m, grid)
	if got != expected {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n, m, grid := generateCaseB(rng)
		if err := runCaseB(bin, n, m, grid); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			for _, row := range grid {
				fmt.Fprintln(os.Stderr, string(row))
			}
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

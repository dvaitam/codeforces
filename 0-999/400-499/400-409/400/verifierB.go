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

func expected(n, m int, grid []string) string {
	dist := make(map[int]struct{})
	for i := 0; i < n; i++ {
		line := grid[i]
		posG, posS := -1, -1
		for j := 0; j < m; j++ {
			switch line[j] {
			case 'G':
				posG = j
			case 'S':
				posS = j
			}
		}
		if posG > posS {
			return "-1"
		}
		d := posS - posG
		if d > 0 {
			dist[d] = struct{}{}
		}
	}
	return fmt.Sprintf("%d", len(dist))
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

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for tcase := 0; tcase < 100; tcase++ {
		n := rng.Intn(5) + 1
		m := rng.Intn(5) + 2 // ensure m>=2
		grid := make([]string, n)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
		for i := 0; i < n; i++ {
			posG := rng.Intn(m)
			posS := rng.Intn(m)
			for posS == posG {
				posS = rng.Intn(m)
			}
			if posG > posS && rng.Intn(2) == 0 {
				posG, posS = posS, posG // sometimes make G before S
			}
			row := make([]byte, m)
			for j := 0; j < m; j++ {
				row[j] = '*'
			}
			row[posG] = 'G'
			row[posS] = 'S'
			grid[i] = string(row)
			sb.WriteString(grid[i])
			sb.WriteByte('\n')
		}
		input := sb.String()
		expectedOut := expected(n, m, grid)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", tcase+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expectedOut {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", tcase+1, expectedOut, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

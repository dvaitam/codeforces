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

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(100) + 1
	board := make([][]byte, n)
	for i := range board {
		row := make([]byte, n)
		for j := range row {
			if rng.Intn(2) == 0 {
				row[j] = 'x'
			} else {
				row[j] = 'o'
			}
		}
		board[i] = row
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		sb.Write(board[i])
		sb.WriteByte('\n')
	}
	// compute expected
	dirs := [][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}
	even := true
	for i := 0; i < n && even; i++ {
		for j := 0; j < n && even; j++ {
			cnt := 0
			for _, d := range dirs {
				ni, nj := i+d[0], j+d[1]
				if ni >= 0 && ni < n && nj >= 0 && nj < n && board[ni][nj] == 'o' {
					cnt++
				}
			}
			if cnt%2 != 0 {
				even = false
			}
		}
	}
	if even {
		return sb.String(), "YES"
	}
	return sb.String(), "NO"
}

func runCase(bin, input, expected string) error {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	outStr := strings.TrimSpace(out.String())
	exp := strings.TrimSpace(expected)
	if outStr != exp {
		return fmt.Errorf("expected %q got %q", exp, outStr)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

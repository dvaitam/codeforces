package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func expectedB(n int, moves [][2]int) string {
	rows := make([]bool, n+1)
	cols := make([]bool, n+1)
	remRows, remCols := n, n
	res := make([]string, len(moves))
	for i, mv := range moves {
		x, y := mv[0], mv[1]
		if !rows[x] {
			rows[x] = true
			remRows--
		}
		if !cols[y] {
			cols[y] = true
			remCols--
		}
		res[i] = fmt.Sprint(remRows * remCols)
	}
	return strings.Join(res, " ")
}

func generateCaseB(rng *rand.Rand) (int, [][2]int) {
	n := rng.Intn(10) + 1
	maxM := n * n
	m := rng.Intn(maxM) + 1
	used := make(map[[2]int]bool)
	moves := make([][2]int, 0, m)
	for len(moves) < m {
		x := rng.Intn(n) + 1
		y := rng.Intn(n) + 1
		key := [2]int{x, y}
		if !used[key] {
			used[key] = true
			moves = append(moves, key)
		}
	}
	return n, moves
}

func runCaseB(bin string, n int, moves [][2]int) error {
	var input strings.Builder
	input.WriteString(fmt.Sprintf("%d %d\n", n, len(moves)))
	for _, mv := range moves {
		input.WriteString(fmt.Sprintf("%d %d\n", mv[0], mv[1]))
	}
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	expect := expectedB(n, moves)
	if got != expect {
		return fmt.Errorf("expected %s got %s", expect, got)
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
		n, moves := generateCaseB(rng)
		if err := runCaseB(bin, n, moves); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

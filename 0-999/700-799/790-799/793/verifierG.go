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

func run(bin, input string) (string, error) {
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

func maxMatching(board [][]bool) int {
	n := len(board)
	matchV := make([]int, n)
	for i := range matchV {
		matchV[i] = -1
	}
	var dfs func(int, []bool) bool
	dfs = func(u int, vis []bool) bool {
		for v := 0; v < n; v++ {
			if board[u][v] && !vis[v] {
				vis[v] = true
				if matchV[v] == -1 || dfs(matchV[v], vis) {
					matchV[v] = u
					return true
				}
			}
		}
		return false
	}
	result := 0
	for u := 0; u < n; u++ {
		vis := make([]bool, n)
		if dfs(u, vis) {
			result++
		}
	}
	return result
}

func solveG(n int, removed [][4]int) string {
	board := make([][]bool, n)
	for i := 0; i < n; i++ {
		row := make([]bool, n)
		for j := 0; j < n; j++ {
			row[j] = true
		}
		board[i] = row
	}
	for _, rec := range removed {
		x1, y1, x2, y2 := rec[0]-1, rec[1]-1, rec[2]-1, rec[3]-1
		for i := x1; i <= x2; i++ {
			for j := y1; j <= y2; j++ {
				board[i][j] = false
			}
		}
	}
	ans := maxMatching(board)
	return fmt.Sprint(ans)
}

func genCase(rng *rand.Rand) (int, [][4]int) {
	n := rng.Intn(6) + 1
	q := rng.Intn(4)
	recs := make([][4]int, q)
	used := make([][]bool, n)
	for i := 0; i < n; i++ {
		used[i] = make([]bool, n)
	}
	for idx := 0; idx < q; idx++ {
		x1 := rng.Intn(n) + 1
		y1 := rng.Intn(n) + 1
		x2 := x1 + rng.Intn(n-x1+1)
		y2 := y1 + rng.Intn(n-y1+1)
		recs[idx] = [4]int{x1, y1, x2, y2}
		for i := x1 - 1; i <= x2-1; i++ {
			for j := y1 - 1; j <= y2-1; j++ {
				used[i][j] = true
			}
		}
	}
	return n, recs
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n, recs := genCase(rng)
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d\n", n)
		fmt.Fprintf(&sb, "%d\n", len(recs))
		for _, r := range recs {
			fmt.Fprintf(&sb, "%d %d %d %d\n", r[0], r[1], r[2], r[3])
		}
		expect := solveG(n, recs)
		got, err := run(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expect {
			fmt.Fprintf(os.Stderr, "case %d failed\ninput:\n%sexpected %s got %s\n", i+1, sb.String(), expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

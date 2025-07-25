package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

const MOD = 1000000007

func runCandidate(bin string, input string) (string, error) {
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

func countPalPaths(grid [][]byte) int {
	n := len(grid)
	m := len(grid[0])
	path := make([]byte, n+m-1)
	var cnt int
	var dfs func(r, c, idx int)
	dfs = func(r, c, idx int) {
		path[idx] = grid[r][c]
		if r == n-1 && c == m-1 {
			if isPalindrome(path) {
				cnt++
			}
			return
		}
		if r+1 < n {
			dfs(r+1, c, idx+1)
		}
		if c+1 < m {
			dfs(r, c+1, idx+1)
		}
	}
	dfs(0, 0, 0)
	return cnt % MOD
}

func isPalindrome(s []byte) bool {
	i, j := 0, len(s)-1
	for i < j {
		if s[i] != s[j] {
			return false
		}
		i++
		j--
	}
	return true
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(1))
	letters := []byte("ab")
	for i := 0; i < 100; i++ {
		n := rng.Intn(4) + 2
		m := rng.Intn(4) + 2
		grid := make([][]byte, n)
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d\n", n, m)
		for r := 0; r < n; r++ {
			row := make([]byte, m)
			for c := 0; c < m; c++ {
				row[c] = letters[rng.Intn(len(letters))]
			}
			grid[r] = row
			sb.Write(row)
			sb.WriteByte('\n')
		}
		input := sb.String()
		want := countPalPaths(grid)
		out, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		var got int
		if _, err := fmt.Sscanf(out, "%d", &got); err != nil || got%MOD != want {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %s\n", i+1, want, strings.TrimSpace(out))
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

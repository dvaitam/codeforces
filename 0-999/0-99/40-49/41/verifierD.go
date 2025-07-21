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

// solve uses the DP algorithm from 41D.go to compute expected result
func solve(n, m, k int, grid []string) (int, int, string) {
	K := k + 1
	mp := make([][]int, n+1)
	for i := 0; i <= n; i++ {
		mp[i] = make([]int, m+1)
	}
	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			mp[i][j] = int(grid[i-1][j-1] - '0')
		}
	}
	f := make([][][]int, n+1)
	dir := make([][][]byte, n+1)
	for i := 0; i <= n; i++ {
		f[i] = make([][]int, m+1)
		dir[i] = make([][]byte, m+1)
		for j := 0; j <= m; j++ {
			f[i][j] = make([]int, K)
			dir[i][j] = make([]byte, K)
			for u := 0; u < K; u++ {
				f[i][j][u] = -1
			}
		}
	}
	for j := 1; j <= m; j++ {
		r := mp[n][j] % K
		f[n][j][r] = mp[n][j]
	}
	for i := n - 1; i >= 1; i-- {
		for j := 1; j <= m; j++ {
			for u := 0; u < K; u++ {
				u0 := (u + mp[i][j]) % K
				if j > 1 && f[i+1][j-1][u] >= 0 {
					v := f[i+1][j-1][u] + mp[i][j]
					if v > f[i][j][u0] {
						f[i][j][u0] = v
						dir[i][j][u0] = 'R'
					}
				}
				if j < m && f[i+1][j+1][u] >= 0 {
					v := f[i+1][j+1][u] + mp[i][j]
					if v > f[i][j][u0] {
						f[i][j][u0] = v
						dir[i][j][u0] = 'L'
					}
				}
			}
		}
	}
	ans := -1
	startCol := 1
	for j := 1; j <= m; j++ {
		if f[1][j][0] > ans {
			ans = f[1][j][0]
			startCol = j
		}
	}
	if ans < 0 {
		return -1, 0, ""
	}
	var path []byte
	var start int
	var find func(i, j, u int)
	find = func(i, j, u int) {
		if i == n {
			start = j
			return
		}
		sum := f[i][j][u]
		u0 := (sum - mp[i][j]) % K
		if u0 < 0 {
			u0 += K
		}
		d := dir[i][j][u]
		if d == 'R' {
			find(i+1, j-1, u0)
		} else {
			find(i+1, j+1, u0)
		}
		path = append(path, d)
	}
	find(1, startCol, 0)
	return ans, start, string(path)
}

func generateCase(rng *rand.Rand) (string, int, int, string) {
	n := rng.Intn(4) + 2
	m := rng.Intn(4) + 2
	k := rng.Intn(3)
	grid := make([]string, n)
	for i := 0; i < n; i++ {
		b := make([]byte, m)
		for j := 0; j < m; j++ {
			b[j] = byte('0' + rng.Intn(10))
		}
		grid[i] = string(b)
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, m, k))
	for i := 0; i < n; i++ {
		sb.WriteString(grid[i] + "\n")
	}
	ans, start, path := solve(n, m, k, grid)
	if ans < 0 {
		return sb.String(), -1, 0, ""
	}
	return sb.String(), ans, start, path
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
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, expAns, expStart, expPath := generateCase(rng)
		out, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		lines := strings.Split(out, "\n")
		if expAns < 0 {
			if strings.TrimSpace(lines[0]) != "-1" {
				fmt.Fprintf(os.Stderr, "case %d failed: expected -1 got %s\ninput:\n%s", i+1, out, in)
				os.Exit(1)
			}
			continue
		}
		if len(lines) < 3 {
			fmt.Fprintf(os.Stderr, "case %d failed: expected 3 lines got %d\ninput:\n%s", i+1, len(lines), in)
			os.Exit(1)
		}
		gotAns, err1 := strconv.Atoi(strings.TrimSpace(lines[0]))
		gotStart, err2 := strconv.Atoi(strings.TrimSpace(lines[1]))
		gotPath := strings.TrimSpace(lines[2])
		if err1 != nil || err2 != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: invalid integers\ninput:\n%s", i+1, in)
			os.Exit(1)
		}
		if gotAns != expAns || gotStart != expStart || gotPath != expPath {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d %d %q got %d %d %q\ninput:\n%s", i+1, expAns, expStart, expPath, gotAns, gotStart, gotPath, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

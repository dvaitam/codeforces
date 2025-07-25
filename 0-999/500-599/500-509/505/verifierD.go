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

func expected(n int, pairs [][2]int) int {
	reach := make([][]bool, n+1)
	for i := range reach {
		reach[i] = make([]bool, n+1)
		if i > 0 {
			reach[i][i] = true
		}
	}
	for _, p := range pairs {
		reach[p[0]][p[1]] = true
	}
	// transitive closure
	for k := 1; k <= n; k++ {
		for i := 1; i <= n; i++ {
			if reach[i][k] {
				for j := 1; j <= n; j++ {
					if reach[k][j] {
						reach[i][j] = true
					}
				}
			}
		}
	}
	comp := make([]int, n+1)
	for i := range comp {
		comp[i] = -1
	}
	cid := 0
	for v := 1; v <= n; v++ {
		if comp[v] != -1 {
			continue
		}
		// BFS
		queue := []int{v}
		comp[v] = cid
		for len(queue) > 0 {
			x := queue[0]
			queue = queue[1:]
			for y := 1; y <= n; y++ {
				if comp[y] == -1 && reach[x][y] && reach[y][x] {
					comp[y] = cid
					queue = append(queue, y)
				}
			}
		}
		cid++
	}
	C := cid
	size := make([]int, C)
	for i := 1; i <= n; i++ {
		size[comp[i]]++
	}
	closureC := make([][]bool, C)
	for i := range closureC {
		closureC[i] = make([]bool, C)
	}
	for i := 1; i <= n; i++ {
		for j := 1; j <= n; j++ {
			if reach[i][j] && comp[i] != comp[j] {
				closureC[comp[i]][comp[j]] = true
			}
		}
	}
	// transitive closure on components already satisfied due to previous step
	// compute transitive reduction edge count
	edgeCnt := 0
	for i := 0; i < C; i++ {
		for j := 0; j < C; j++ {
			if i == j || !closureC[i][j] {
				continue
			}
			redundant := false
			for k := 0; k < C; k++ {
				if k != i && k != j && closureC[i][k] && closureC[k][j] {
					redundant = true
					break
				}
			}
			if !redundant {
				edgeCnt++
			}
		}
	}
	ans := edgeCnt
	for i := 0; i < C; i++ {
		if size[i] > 1 {
			ans += size[i]
		}
	}
	return ans
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for tc := 0; tc < 100; tc++ {
		n := rng.Intn(6) + 2
		maxEdges := n * (n - 1)
		limit := 8
		if maxEdges < limit {
			limit = maxEdges
		}
		m := rng.Intn(limit) + 1
		used := map[[2]int]struct{}{}
		pairs := make([][2]int, 0, m)
		for len(pairs) < m {
			a := rng.Intn(n) + 1
			b := rng.Intn(n) + 1
			if a == b {
				continue
			}
			p := [2]int{a, b}
			if _, ok := used[p]; ok {
				continue
			}
			used[p] = struct{}{}
			pairs = append(pairs, p)
		}
		input := fmt.Sprintf("%d %d\n", n, m)
		for _, p := range pairs {
			input += fmt.Sprintf("%d %d\n", p[0], p[1])
		}
		exp := expected(n, pairs)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", tc+1, err, input)
			os.Exit(1)
		}
		if out != fmt.Sprintf("%d", exp) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %s\ninput:\n%s", tc+1, exp, out, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

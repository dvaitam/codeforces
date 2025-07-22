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

type pair struct {
	v int
	p int
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func expected(n, m, k, w int, levels [][]string) string {
	maxCost := n * m
	dist := make([][]int, k)
	for i := 0; i < k; i++ {
		dist[i] = make([]int, k)
	}
	for i := 0; i < k; i++ {
		for j := i + 1; j < k; j++ {
			diff := 0
			for r := 0; r < n; r++ {
				a := levels[i][r]
				b := levels[j][r]
				for c := 0; c < m; c++ {
					if a[c] != b[c] {
						diff++
					}
				}
			}
			cost := min(maxCost, diff*w)
			dist[i][j] = cost
			dist[j][i] = cost
		}
	}
	used := make([]bool, k)
	minCost := make([]int, k)
	prev := make([]int, k)
	used[0] = true
	total := maxCost
	ans := []pair{{0, -1}}
	for i := 1; i < k; i++ {
		minCost[i] = dist[0][i]
		prev[i] = 0
	}
	for cnt := 1; cnt < k; cnt++ {
		v := -1
		for i := 0; i < k; i++ {
			if !used[i] && (v == -1 || minCost[i] < minCost[v]) {
				v = i
			}
		}
		if minCost[v] >= maxCost {
			total += maxCost
			ans = append(ans, pair{v, -1})
		} else {
			total += minCost[v]
			ans = append(ans, pair{v, prev[v]})
		}
		used[v] = true
		for u := 0; u < k; u++ {
			if !used[u] && dist[v][u] < minCost[u] {
				minCost[u] = dist[v][u]
				prev[u] = v
			}
		}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", total))
	for _, pr := range ans {
		x := pr.v + 1
		y := pr.p + 1
		if pr.p < 0 {
			y = 0
		}
		sb.WriteString(fmt.Sprintf("%d %d\n", x, y))
	}
	return strings.TrimSpace(sb.String())
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
	got := strings.TrimSpace(out.String())
	if got != expected {
		return fmt.Errorf("expected:\n%s\ngot:\n%s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n := rng.Intn(2) + 1
		m := rng.Intn(2) + 1
		k := rng.Intn(2) + 1
		w := rng.Intn(3) + 1
		levels := make([][]string, k)
		for lvl := 0; lvl < k; lvl++ {
			levels[lvl] = make([]string, n)
			for r := 0; r < n; r++ {
				rowBytes := make([]byte, m)
				for c := 0; c < m; c++ {
					rowBytes[c] = byte('a' + rng.Intn(3))
				}
				levels[lvl][r] = string(rowBytes)
			}
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d %d %d\n", n, m, k, w))
		for lvl := 0; lvl < k; lvl++ {
			for r := 0; r < n; r++ {
				sb.WriteString(levels[lvl][r])
				sb.WriteByte('\n')
			}
		}
		input := sb.String()
		exp := expected(n, m, k, w, levels)
		if err := runCase(bin, input, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"
)

type node struct{ val, no int }

func solveD(n, m int, mat [][]int) []int {
	rec := make([]node, m)
	maxNodes := m*2 + n*2 + 5
	join := make([]int, maxNodes)
	v := make([][]int, maxNodes)
	cnt := 0
	ans := make([]int, 0, m)
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			rec[j].val = mat[i][j]
			rec[j].no = j
		}
		sort.Slice(rec, func(a, b int) bool { return rec[a].val < rec[b].val })
		for j := 0; j < m; j++ {
			if rec[j].val < 0 {
				continue
			}
			if j == 0 || rec[j].val != rec[j-1].val {
				cnt++
			}
			end := m + cnt + 1
			v[rec[j].no] = append(v[rec[j].no], end)
			join[end]++
			start := m + cnt
			v[start] = append(v[start], rec[j].no)
			join[rec[j].no]++
		}
		cnt++
	}
	totalNodes := m + cnt + 3
	q := make([]int, 0, totalNodes)
	for i := 0; i < totalNodes && i < len(join); i++ {
		if join[i] == 0 {
			q = append(q, i)
		}
	}
	head := 0
	for head < len(q) {
		t := q[head]
		head++
		if t < m {
			ans = append(ans, t)
		}
		for _, to := range v[t] {
			join[to]--
			if join[to] == 0 {
				q = append(q, to)
			}
		}
	}
	if len(ans) < m {
		return nil
	}
	for i := 0; i < m; i++ {
		ans[i]++
	}
	return ans[:m]
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(3) + 1
	m := rng.Intn(3) + 1
	mat := make([][]int, n)
	for i := 0; i < n; i++ {
		row := make([]int, m)
		for j := 0; j < m; j++ {
			if rng.Float64() < 0.3 {
				row[j] = -1
			} else {
				row[j] = rng.Intn(10)
			}
		}
		mat[i] = row
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(mat[i][j]))
		}
		sb.WriteByte('\n')
	}
	res := solveD(n, m, mat)
	if res == nil {
		return sb.String(), "-1\n"
	}
	var out strings.Builder
	for i, v := range res {
		if i > 0 {
			out.WriteByte(' ')
		}
		out.WriteString(strconv.Itoa(v))
	}
	out.WriteByte('\n')
	return sb.String(), out.String()
}

func runCase(exe string, input, expected string) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	exp := strings.TrimSpace(expected)
	if got != exp {
		return fmt.Errorf("expected %q got %q", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(exe, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

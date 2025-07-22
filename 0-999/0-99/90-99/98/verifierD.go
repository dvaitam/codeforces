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

type pair struct{ x, y int }

func dfs0(a []int, n int, m, x, y, z int) int {
	if m == 0 {
		return 0
	}
	i := m
	for i > 0 && a[i] == a[m] {
		i--
	}
	if z != 0 && i+1 < m {
		t1 := dfs0(a, n, i, x, y, 0) + (m - i) + dfs0(a, n, i, y, x, 0) + (m - i) + dfs0(a, n, i, x, y, 1)
		t2 := dfs0(a, n, n-1, x, 6-x-y, 0) + 1 + dfs0(a, n, n-1, 6-x-y, y, 0)
		if t1 < t2 {
			return t1
		}
		return t2
	}
	return dfs0(a, n, i, x, 6-x-y, 0) + (m - i) + dfs0(a, n, i, 6-x-y, y, 0)
}

func dfs(a []int, n int, m, x, y, z int, out *[]pair) int {
	if m == 0 {
		return 0
	}
	i := m
	for i > 0 && a[i] == a[m] {
		i--
	}
	if z != 0 && i+1 < m {
		t1 := dfs0(a, n, i, x, y, 0) + (m - i) + dfs0(a, n, i, y, x, 0) + (m - i) + dfs0(a, n, i, x, y, 1)
		t2 := dfs0(a, n, n-1, x, 6-x-y, 0) + 1 + dfs0(a, n, n-1, 6-x-y, y, 0)
		if t1 < t2 {
			dfs(a, n, i, x, y, 0, out)
			for j := m; j > i; j-- {
				*out = append(*out, pair{x, 6 - x - y})
			}
			dfs(a, n, i, y, x, 0, out)
			for j := m; j > i; j-- {
				*out = append(*out, pair{6 - x - y, y})
			}
			dfs(a, n, i, x, y, 1, out)
			return t1
		}
		dfs(a, n, n-1, x, 6-x-y, 0, out)
		*out = append(*out, pair{x, y})
		dfs(a, n, n-1, 6-x-y, y, 0, out)
		return t2
	}
	t := dfs(a, n, i, x, 6-x-y, 0, out) + (m - i)
	for j := m; j > i; j-- {
		*out = append(*out, pair{x, y})
	}
	t += dfs(a, n, i, 6-x-y, y, 0, out)
	return t
}

func solve(input string) string {
	fields := strings.Fields(input)
	if len(fields) == 0 {
		return ""
	}
	var n int
	fmt.Sscan(fields[0], &n)
	a := make([]int, n+1)
	idx := 1
	for i := n; i >= 1; i-- {
		fmt.Sscan(fields[idx], &a[i])
		idx++
	}
	out := []pair{}
	cnt := dfs(a, n, n, 1, 3, 1, &out)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", cnt))
	for _, p := range out {
		sb.WriteString(fmt.Sprintf("%d %d\n", p.x, p.y))
	}
	return strings.TrimRight(sb.String(), "\n")
}

func genCase(rng *rand.Rand) string {
	n := rng.Intn(6) + 1
	arr := make([]int, n)
	arr[0] = rng.Intn(20) + 1
	for i := 1; i < n; i++ {
		arr[i] = rng.Intn(arr[i-1]) + 1
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		sb.WriteString(fmt.Sprintf("%d ", arr[i]))
	}
	sb.WriteString("\n")
	return sb.String()
}

func runCase(bin string, input, expected string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != expected {
		return fmt.Errorf("expected\n%s\ngot\n%s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input := genCase(rng)
		expected := solve(strings.TrimSpace(input))
		if err := runCase(bin, input, expected); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

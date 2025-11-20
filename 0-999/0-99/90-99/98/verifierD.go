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

func parseInputCase(input string) (int, []int, error) {
	reader := strings.NewReader(input)
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return 0, nil, fmt.Errorf("failed to read n: %v", err)
	}
	diam := make([]int, n)
	for i := 0; i < n; i++ {
		if _, err := fmt.Fscan(reader, &diam[i]); err != nil {
			return 0, nil, fmt.Errorf("failed to read diameter %d: %v", i+1, err)
		}
	}
	return n, diam, nil
}

func simulateMoves(n int, diam []int, moves []pair) error {
	pegs := make([][]int, 3)
	for i := 0; i < n; i++ {
		pegs[0] = append(pegs[0], i)
	}
	for idx, mv := range moves {
		s, t := mv.x, mv.y
		if s < 1 || s > 3 || t < 1 || t > 3 {
			return fmt.Errorf("move %d: invalid pillar %d %d", idx+1, s, t)
		}
		if s == t {
			return fmt.Errorf("move %d: source equals destination", idx+1)
		}
		src := s - 1
		dst := t - 1
		if len(pegs[src]) == 0 {
			return fmt.Errorf("move %d: source pillar empty", idx+1)
		}
		disk := pegs[src][len(pegs[src])-1]
		pegs[src] = pegs[src][:len(pegs[src])-1]
		if len(pegs[dst]) > 0 {
			top := pegs[dst][len(pegs[dst])-1]
			if diam[top] < diam[disk] {
				return fmt.Errorf("move %d: cannot place larger disk on smaller", idx+1)
			}
		}
		pegs[dst] = append(pegs[dst], disk)
	}
	if len(pegs[0]) != 0 || len(pegs[1]) != 0 || len(pegs[2]) != n {
		return fmt.Errorf("not all disks moved to third pillar")
	}
	for i, disk := range pegs[2] {
		if disk != i {
			return fmt.Errorf("final order incorrect")
		}
	}
	return nil
}

func validateOutput(output, input string, expectedMoves int) error {
	output = strings.TrimSpace(output)
	if output == "" {
		return fmt.Errorf("empty output")
	}
	fields := strings.Fields(output)
	wantTokens := 1 + 2*expectedMoves
	if len(fields) != wantTokens {
		return fmt.Errorf("expected %d move pairs (%d tokens) but got %d tokens", expectedMoves, wantTokens, len(fields))
	}
	movesCount, err := strconv.Atoi(fields[0])
	if err != nil {
		return fmt.Errorf("invalid move count: %v", err)
	}
	if movesCount != expectedMoves {
		return fmt.Errorf("expected %d moves but got %d", expectedMoves, movesCount)
	}
	moves := make([]pair, movesCount)
	idx := 1
	for i := 0; i < movesCount; i++ {
		s, err := strconv.Atoi(fields[idx])
		if err != nil {
			return fmt.Errorf("invalid source in move %d: %v", i+1, err)
		}
		idx++
		t, err := strconv.Atoi(fields[idx])
		if err != nil {
			return fmt.Errorf("invalid destination in move %d: %v", i+1, err)
		}
		idx++
		moves[i] = pair{s, t}
	}
	n, diam, err := parseInputCase(input)
	if err != nil {
		return err
	}
	return simulateMoves(n, diam, moves)
}

func runCase(bin string, input string) error {
	expected := solve(strings.TrimSpace(input))
	parts := strings.Fields(expected)
	if len(parts) == 0 {
		return fmt.Errorf("reference produced empty output")
	}
	expectedMoves, err := strconv.Atoi(parts[0])
	if err != nil {
		return fmt.Errorf("invalid reference move count: %v", err)
	}
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return validateOutput(out.String(), input, expectedMoves)
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
		if err := runCase(bin, input); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

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

type testCase struct {
	input    string
	expected string
	n, k     int
	arr      []int
}

func generateCase(rng *rand.Rand) testCase {
	n := rng.Intn(15) + 1
	k := rng.Intn(n) + 1
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = rng.Intn(10) + 1
	}
	// build input string
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, k))
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	expected := feasible(n, k, arr)
	return testCase{input: sb.String(), expected: expected, n: n, k: k, arr: arr}
}

func feasible(n, k int, a []int) string {
	const maxV = 5000
	if n < k {
		return "NO"
	}
	pos := make([][]int, maxV+1)
	for i, v := range a {
		if v >= 1 && v <= maxV {
			pos[v] = append(pos[v], i)
			if len(pos[v]) > k {
				return "NO"
			}
		} else {
			return "NO" // out of range
		}
	}
	ans := make([]int, n)
	idx := 0
	for v := 1; v <= maxV; v++ {
		for _, p := range pos[v] {
			ans[p] = (idx % k) + 1
			idx++
		}
	}
	if idx < k {
		return "NO"
	}
	var sb strings.Builder
	sb.WriteString("YES\n")
	for i, c := range ans {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", c))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func runCase(bin string, t testCase) error {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(t.input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	outStr := strings.TrimSpace(out.String())
	expStr := strings.TrimSpace(t.expected)
	if expStr == "NO" {
		if outStr != "NO" {
			return fmt.Errorf("expected NO got %q", outStr)
		}
		return nil
	}
	// expected YES coloring
	lines := strings.Split(outStr, "\n")
	if len(lines) != 2 || strings.TrimSpace(lines[0]) != "YES" {
		return fmt.Errorf("expected YES with coloring, got %q", outStr)
	}
	cols := strings.Fields(lines[1])
	if len(cols) != t.n {
		return fmt.Errorf("expected %d colors, got %d", t.n, len(cols))
	}
	// verify validity
	k := t.k
	a := t.arr
	usedColor := make([]bool, k+1)
	countsByValue := make(map[int]map[int]bool)
	for i, colStr := range cols {
		var c int
		fmt.Sscanf(colStr, "%d", &c)
		if c < 1 || c > k {
			return fmt.Errorf("color out of range")
		}
		usedColor[c] = true
		v := a[i]
		if countsByValue[v] == nil {
			countsByValue[v] = make(map[int]bool)
		}
		if countsByValue[v][c] {
			return fmt.Errorf("duplicate color for value")
		}
		countsByValue[v][c] = true
	}
	for c := 1; c <= k; c++ {
		if !usedColor[c] {
			return fmt.Errorf("color %d unused", c)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := generateCase(rng)
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

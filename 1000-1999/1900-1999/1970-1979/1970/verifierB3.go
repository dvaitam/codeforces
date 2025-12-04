package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

func generateCase(rng *rand.Rand) (string, []int) {
	n := rng.Intn(10) + 2
	arr := make([]int, n)
	for i := range arr {
		arr[i] = rng.Intn(n + 1)
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	return sb.String(), arr
}

func isPossible(a []int) bool {
	n := len(a)
	b := make([]int, n)
	copy(b, a)
	sort.Ints(b)

	hasDuplicate := false
	for i := 0; i < n-1; i++ {
		if b[i] == b[i+1] {
			hasDuplicate = true
			break
		}
	}

	if hasDuplicate || b[0] == 0 {
		return true
	}
	if n == 2 {
		return false
	}
	return true
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func verify(input string, output string) error {
	lines := strings.Fields(input)
	var n int
	fmt.Sscanf(lines[0], "%d", &n)
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Sscanf(lines[i+1], "%d", &a[i])
	}

	outFields := strings.Fields(output)
	if len(outFields) == 0 {
		return fmt.Errorf("no output")
	}

	ans := outFields[0]
	possible := isPossible(a)

	if ans == "NO" {
		if possible {
			return fmt.Errorf("output NO but solution likely exists")
		}
		return nil
	}
	if ans != "YES" {
		return fmt.Errorf("expected YES or NO, got %s", ans)
	}
	if !possible {
		return fmt.Errorf("output YES but solution impossible")
	}

	if len(outFields) < 1+2*n+n {
		return fmt.Errorf("insufficient output fields")
	}

	xs := make([]int, n)
	ys := make([]int, n)
	xSet := make(map[int]bool)

	idx := 1
	for i := 0; i < n; i++ {
		fmt.Sscanf(outFields[idx], "%d", &xs[i])
		idx++
		fmt.Sscanf(outFields[idx], "%d", &ys[i])
		idx++

		if xs[i] < 1 || xs[i] > n || ys[i] < 1 || ys[i] > n {
			return fmt.Errorf("coordinates out of bounds: (%d, %d)", xs[i], ys[i])
		}
		if xSet[xs[i]] {
			return fmt.Errorf("duplicate x coordinate (column): %d", xs[i])
		}
		xSet[xs[i]] = true
	}

	for i := 0; i < n; i++ {
		var targetIdx int
		fmt.Sscanf(outFields[idx], "%d", &targetIdx)
		idx++
		if targetIdx < 1 || targetIdx > n {
			return fmt.Errorf("target index out of bounds: %d", targetIdx)
		}

		// 0-indexed for array access
		t := targetIdx - 1
		dist := abs(xs[i]-xs[t]) + abs(ys[i]-ys[t])
		if dist != a[i] {
			return fmt.Errorf("wizard %d distance mismatch: required %d, got %d (target %d, pos (%d,%d) -> (%d,%d))", 
				i+1, a[i], dist, targetIdx, xs[i], ys[i], xs[t], ys[t])
		}
	}

	return nil
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v", err)
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB3.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, _ := generateCase(rng)
		out, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if err := verify(in, out); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s\noutput:\n%s", i+1, err, in, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

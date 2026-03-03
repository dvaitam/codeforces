package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// Oracle solver using DFS with pruning
func solve(n, a, b int, h []int) int {
	hCopy := make([]int, n+2)
	for i := 1; i <= n; i++ {
		hCopy[i] = h[i-1]
	}
	ans := 1000

	var dfs func(int, int)
	dfs = func(idx, sum int) {
		if sum >= ans {
			return
		}
		if idx == n {
			if hCopy[n-1] < 0 && hCopy[n] < 0 {
				ans = sum
			}
			return
		}

		needed := 0
		if hCopy[idx-1] >= 0 {
			needed = hCopy[idx-1]/b + 1
		}

		oldPrev := hCopy[idx-1]
		oldCur := hCopy[idx]
		oldNext := hCopy[idx+1]

		for i := needed; i <= 16; i++ {
			hCopy[idx-1] -= i * b
			hCopy[idx] -= i * a
			hCopy[idx+1] -= i * b

			dfs(idx+1, sum+i)

			hCopy[idx-1] = oldPrev
			hCopy[idx] = oldCur
			hCopy[idx+1] = oldNext

			if hCopy[idx-1] < 0 && hCopy[idx] < 0 && hCopy[idx+1] < 0 && idx < n-1 {
				break
			}
		}
	}
	dfs(2, 0)
	return ans
}

func runCase(exe string, n, a, b int, h []int) error {
	input := fmt.Sprintf("%d %d %d\n", n, a, b)
	for i := 0; i < n; i++ {
		input += fmt.Sprintf("%d", h[i])
		if i+1 < n {
			input += " "
		}
	}
	input += "\n"

	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, string(out))
	}

	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	if len(lines) < 2 {
		return fmt.Errorf("output should have at least 2 lines, got %d", len(lines))
	}

	gotT, err := strconv.Atoi(strings.TrimSpace(lines[0]))
	if err != nil {
		return fmt.Errorf("invalid t: %v", err)
	}

	expectedT := solve(n, a, b, h)
	if gotT != expectedT {
		return fmt.Errorf("incorrect minimum shots: expected t=%d, got %d", expectedT, gotT)
	}

	fields := strings.Fields(lines[1])
	if len(fields) != gotT {
		return fmt.Errorf("shot count mismatch: expected %d shots, got %d", gotT, len(fields))
	}

	curH := make([]int, n)
	copy(curH, h)
	for _, s := range fields {
		idx, err := strconv.Atoi(s)
		if err != nil {
			return fmt.Errorf("invalid shot index %q: %v", s, err)
		}
		if idx < 2 || idx > n-1 {
			return fmt.Errorf("shot index %d out of range [2, %d]", idx, n-1)
		}
		curH[idx-1] -= a
		curH[idx-2] -= b
		curH[idx] -= b
	}

	for i := 0; i < n; i++ {
		if curH[i] >= 0 {
			return fmt.Errorf("archer %d still alive with health %d", i+1, curH[i])
		}
	}

	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Hardcoded failing case: n=5, a=3, b=2, h=[2, 4, 3, 3, 1]
	// Optimal t is 4 (e.g., 2 2 4 4). 3 is impossible because archer 4 health becomes 0.
	if err := runCase(exe, 5, 3, 2, []int{2, 4, 3, 3, 1}); err != nil {
		fmt.Fprintf(os.Stderr, "Case [2 4 3 3 1] failed: %v\n", err)
		os.Exit(1)
	}

	for i := 0; i < 50; i++ {
		n := rng.Intn(3) + 3 
		if i > 40 {
			n = rng.Intn(8) + 3
		}
		a := rng.Intn(5) + 2
		b := rng.Intn(a-1) + 1
		h := make([]int, n)
		for j := 0; j < n; j++ {
			h[j] = rng.Intn(15) + 1
		}
		if err := runCase(exe, n, a, b, h); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%d %d %d\n%v\n", i+1, err, n, a, b, h)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

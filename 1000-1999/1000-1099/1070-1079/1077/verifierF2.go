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

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func expected(n, k, x int, a []int) int64 {
	if n/k > x {
		return -1
	}
	g := make([]int64, n+2)
	f := make([]int64, n+2)
	for i := 1; i <= k && i <= n; i++ {
		g[i] = int64(a[i])
	}
	dui := make([]int, n+2)
	for l := 2; l <= x; l++ {
		for i := 1; i <= n; i++ {
			f[i] = 0
		}
		t1, t2 := 0, -1
		upper := l * k
		if upper > n {
			upper = n
		}
		for i := l; i <= upper; i++ {
			for t1 <= t2 && g[i-1] >= g[dui[t2]] {
				t2--
			}
			t2++
			dui[t2] = i - 1
			for t1 <= t2 && dui[t1] < i-k {
				t1++
			}
			if t1 <= t2 {
				f[i] = g[dui[t1]] + int64(a[i])
			}
		}
		for i := 1; i <= n; i++ {
			g[i] = f[i]
		}
	}
	var ans int64
	start := max(1, n-k+1)
	for i := start; i <= n; i++ {
		if g[i] > ans {
			ans = g[i]
		}
	}
	if ans == 0 {
		return -1
	}
	return ans
}

func runCase(exe string, input string, exp int64) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	valStr := strings.TrimSpace(out.String())
	got, err := strconv.ParseInt(valStr, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid output %q", valStr)
	}
	if got != exp {
		return fmt.Errorf("expected %d got %d", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF2.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for tcase := 0; tcase < 100; tcase++ {
		n := rng.Intn(200) + 1
		k := rng.Intn(n) + 1
		x := rng.Intn(n) + 1
		a := make([]int, n+1)
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d %d\n", n, k, x)
		for i := 1; i <= n; i++ {
			a[i] = rng.Intn(1000) + 1
			if i > 1 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(a[i]))
		}
		sb.WriteByte('\n')
		input := sb.String()
		exp := expected(n, k, x, a)
		if err := runCase(exe, input, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", tcase+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

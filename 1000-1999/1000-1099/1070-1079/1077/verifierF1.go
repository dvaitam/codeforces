package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func expected(n, k, x int, a []int64) int64 {
	const INF int64 = -1 << 60
	d := make([][]int64, n+1)
	for i := range d {
		d[i] = make([]int64, x+1)
		for j := range d[i] {
			d[i][j] = INF
		}
	}
	d[0][0] = 0
	for i := 0; i < n; i++ {
		for j := 0; j < x; j++ {
			if d[i][j] != INF {
				for z := i + 1; z <= i+k && z <= n; z++ {
					val := d[i][j] + a[z]
					if val > d[z][j+1] {
						d[z][j+1] = val
					}
				}
			}
		}
	}
	ans := INF
	for i := n - k + 1; i <= n; i++ {
		if i >= 0 && i <= n && d[i][x] > ans {
			ans = d[i][x]
		}
	}
	if ans == INF {
		return math.MinInt64
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
	if exp == math.MinInt64 {
		if got != -1 {
			return fmt.Errorf("expected -1 got %d", got)
		}
	} else if got != exp {
		return fmt.Errorf("expected %d got %d", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF1.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for tcase := 0; tcase < 100; tcase++ {
		n := rng.Intn(200) + 1
		k := rng.Intn(n) + 1
		x := rng.Intn(n) + 1
		a := make([]int64, n+1)
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d %d\n", n, k, x)
		for i := 1; i <= n; i++ {
			a[i] = rng.Int63n(1000) + 1
			if i > 1 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.FormatInt(a[i], 10))
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

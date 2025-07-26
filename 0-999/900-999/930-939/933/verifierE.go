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

func min64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func expectedCost(n int, a []int64) int64 {
	f := make([]int64, n+2)
	for i := 1; i <= n; i++ {
		fo := a[i]
		if i-2 >= 0 {
			fo += f[i-2]
		}
		var fe int64
		if i-3 >= 0 {
			if a[i] > a[i-1] {
				fe = a[i]
			} else {
				fe = a[i-1]
			}
			fe += f[i-3]
		} else {
			if a[i] > a[i-1] {
				fe = a[i]
			} else {
				fe = a[i-1]
			}
		}
		if fo <= fe {
			f[i] = fo
		} else {
			f[i] = fe
		}
	}
	if n == 0 {
		return 0
	}
	if f[n-1] < f[n] {
		return f[n-1]
	}
	return f[n]
}

func generateCaseE(rng *rand.Rand) (int, []int64) {
	n := rng.Intn(20) + 1
	arr := make([]int64, n)
	for i := range arr {
		arr[i] = rng.Int63n(20)
	}
	return n, arr
}

func runCaseE(bin string, n int, arr []int64) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	sb.WriteByte('\n')
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(lines) == 0 {
		return fmt.Errorf("no output")
	}
	var m int
	if _, err := fmt.Sscan(lines[0], &m); err != nil {
		return fmt.Errorf("failed to read m: %v", err)
	}
	if m != len(lines)-1 {
		return fmt.Errorf("line count mismatch")
	}
	ops := make([]int, m)
	for i := 0; i < m; i++ {
		if _, err := fmt.Sscan(lines[i+1], &ops[i]); err != nil {
			return fmt.Errorf("bad operation index on line %d", i+2)
		}
		if ops[i] < 1 || ops[i] >= n {
			return fmt.Errorf("invalid operation index %d", ops[i])
		}
	}
	// simulate
	arrCopy := append([]int64(nil), arr...)
	var cost int64
	for _, op := range ops {
		if arrCopy[op-1] <= 0 || arrCopy[op] <= 0 {
			return fmt.Errorf("operation on non-positive numbers")
		}
		d := min64(arrCopy[op-1], arrCopy[op])
		arrCopy[op-1] -= d
		arrCopy[op] -= d
		cost += d
	}
	for i := 0; i < n-1; i++ {
		if arrCopy[i] > 0 && arrCopy[i+1] > 0 {
			return fmt.Errorf("game not finished")
		}
		if arrCopy[i] < 0 {
			return fmt.Errorf("negative value")
		}
	}
	if arrCopy[n-1] < 0 {
		return fmt.Errorf("negative value")
	}
	expected := expectedCost(n, append([]int64{0}, arr...))
	if cost != expected {
		return fmt.Errorf("expected cost %d got %d", expected, cost)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n, arr := generateCaseE(rng)
		if err := runCaseE(bin, n, arr); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

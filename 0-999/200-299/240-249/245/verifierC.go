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

func solveC(a []int) int64 {
	n := len(a) - 1
	maxA := 0
	for i := 1; i <= n; i++ {
		if a[i] > maxA {
			maxA = a[i]
		}
	}
	MAX := maxA
	m := (n - 1) / 2
	const INF int64 = 1e18
	f := make([][]int64, n+2)
	for v := n; v >= 1; v-- {
		c1 := 2 * v
		c2 := 2*v + 1
		var f1, f2 []int64
		if c1 <= n {
			f1 = f[c1]
		}
		if c2 <= n {
			f2 = f[c2]
		}
		hasVar := v <= m
		s := make([]int64, MAX+1)
		for k := 0; k <= MAX; k++ {
			if !hasVar && k > 0 {
				s[k] = INF
				continue
			}
			cost := int64(k)
			if f1 != nil {
				cost += f1[k]
			}
			if f2 != nil {
				cost += f2[k]
			}
			s[k] = cost
		}
		for k := MAX - 1; k >= 0; k-- {
			if s[k+1] < s[k] {
				s[k] = s[k+1]
			}
		}
		fv := make([]int64, MAX+1)
		for p := 0; p <= MAX; p++ {
			lb := a[v] - p
			if lb < 0 {
				lb = 0
			}
			if lb > MAX {
				fv[p] = INF
			} else {
				fv[p] = s[lb]
			}
		}
		f[v] = fv
	}
	res := f[1][0]
	if res >= int64(1e17) {
		return -1
	}
	return res
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(10) + 1
	a := make([]int, n+1)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 1; i <= n; i++ {
		a[i] = rng.Intn(5)
		sb.WriteString(fmt.Sprintf("%d", a[i]))
		if i < n {
			sb.WriteByte(' ')
		}
	}
	sb.WriteByte('\n')
	ans := solveC(a)
	return sb.String(), fmt.Sprint(ans)
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
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

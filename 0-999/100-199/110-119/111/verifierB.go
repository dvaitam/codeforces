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

const maxV = 100000

var spf [maxV + 1]int

func initSieve() {
	for i := 2; i <= maxV; i++ {
		if spf[i] == 0 {
			spf[i] = i
			if i <= maxV/i {
				for j := i * i; j <= maxV; j += i {
					if spf[j] == 0 {
						spf[j] = i
					}
				}
			}
		}
	}
}

func solveB(n int, xs, ys []int) []int {
	last := make([]int, maxV+1)
	for i := range last {
		last[i] = -1
	}
	ans := make([]int, n)
	for idx := 0; idx < n; idx++ {
		x := xs[idx]
		y := ys[idx]
		ps := make([]int, 0)
		cs := make([]int, 0)
		xx := x
		for xx > 1 {
			p := spf[xx]
			cnt := 0
			for xx%p == 0 {
				xx /= p
				cnt++
			}
			ps = append(ps, p)
			cs = append(cs, cnt)
		}
		l := idx - y
		r := idx - 1
		res := 0
		var dfs func(int, int)
		dfs = func(step, mul int) {
			if step < 0 {
				if last[mul] < l {
					res++
				}
				last[mul] = r + 1
				return
			}
			cur := mul
			for i := 0; i <= cs[step]; i++ {
				dfs(step-1, cur)
				cur *= ps[step]
			}
		}
		dfs(len(ps)-1, 1)
		ans[idx] = res
	}
	return ans
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(20) + 1
	xs := make([]int, n)
	ys := make([]int, n)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		xs[i] = rng.Intn(100000) + 1
		ys[i] = rng.Intn(i + 1)
		sb.WriteString(fmt.Sprintf("%d %d\n", xs[i], ys[i]))
	}
	out := solveB(n, xs, ys)
	var sbAns strings.Builder
	for _, v := range out {
		sbAns.WriteString(fmt.Sprintf("%d\n", v))
	}
	return sb.String(), sbAns.String()
}

func runCase(exe, input, expected string) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	if strings.TrimSpace(out.String()) != strings.TrimSpace(expected) {
		return fmt.Errorf("expected %q got %q", strings.TrimSpace(expected), strings.TrimSpace(out.String()))
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	initSieve()
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

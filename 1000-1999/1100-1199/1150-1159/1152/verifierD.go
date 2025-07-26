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

const MOD = 1000000007

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func expected(n int) int {
	N := 2 * n
	prev0 := make([]int, n+2)
	prev1 := make([]int, n+2)
	for pos := N - 1; pos >= 0; pos-- {
		maxb := min(pos, N-pos)
		cur0 := make([]int, n+2)
		cur1 := make([]int, n+2)
		for b := 0; b <= maxb; b++ {
			sumAll := 0
			bestVal := 0
			if b+1 <= N-pos-1 {
				p0 := prev0[b+1]
				p1 := prev1[b+1]
				mx := p0
				if p1 > mx {
					mx = p1
				}
				sumAll += mx
				val := 1 + p1 - mx
				if val > bestVal {
					bestVal = val
				}
			}
			if b > 0 {
				p0 := prev0[b-1]
				p1 := prev1[b-1]
				mx := p0
				if p1 > mx {
					mx = p1
				}
				sumAll += mx
				val := 1 + p1 - mx
				if val > bestVal {
					bestVal = val
				}
			}
			cur1[b] = sumAll
			if bestVal > 0 {
				cur0[b] = sumAll + bestVal
			} else {
				cur0[b] = sumAll
			}
		}
		prev0 = cur0
		prev1 = cur1
	}
	res := prev0[0] % MOD
	return res
}

func generateCase(rng *rand.Rand) (string, int) {
	n := rng.Intn(10) + 1
	inp := fmt.Sprintf("%d\n", n)
	return inp, expected(n)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	fixed := []int{1, 2, 3, 10}
	idx := 0
	for ; idx < len(fixed); idx++ {
		n := fixed[idx]
		inp := fmt.Sprintf("%d\n", n)
		exp := strconv.Itoa(expected(n))
		out, err := run(bin, inp)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", idx+1, err, inp)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:%s", idx+1, exp, out, inp)
			os.Exit(1)
		}
	}
	for ; idx < 100; idx++ {
		inp, expVal := generateCase(rng)
		out, err := run(bin, inp)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", idx+1, err, inp)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strconv.Itoa(expVal) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %s\ninput:%s", idx+1, expVal, out, inp)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", idx)
}

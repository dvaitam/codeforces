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

const MOD = 998244353

func expectedCount(n, m int, c []int) []int {
	F := make([]int, m+1)
	F[0] = 1
	for s := 1; s <= m; s++ {
		total := 0
		for _, w := range c {
			if w > s {
				continue
			}
			rem := s - w
			for l := 0; l <= rem; l++ {
				r := rem - l
				total += (F[l] * F[r]) % MOD
			}
		}
		F[s] = total % MOD
	}
	res := make([]int, m)
	for i := 1; i <= m; i++ {
		res[i-1] = F[i] % MOD
	}
	return res
}

type testCase struct {
	n int
	m int
	c []int
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(3) + 1
	m := rng.Intn(5) + 1
	vals := make([]int, n)
	used := make(map[int]struct{})
	for i := 0; i < n; i++ {
		v := rng.Intn(5) + 1
		for {
			if _, ok := used[v]; !ok {
				break
			}
			v = rng.Intn(5) + 1
		}
		used[v] = struct{}{}
		vals[i] = v
	}
	inputSb := strings.Builder{}
	inputSb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for i, v := range vals {
		if i > 0 {
			inputSb.WriteByte(' ')
		}
		inputSb.WriteString(fmt.Sprintf("%d", v))
	}
	inputSb.WriteByte('\n')
	outVals := expectedCount(n, m, vals)
	outSb := strings.Builder{}
	for i, v := range outVals {
		if i > 0 {
			outSb.WriteByte('\n')
		}
		outSb.WriteString(fmt.Sprintf("%d", v))
	}
	outSb.WriteByte('\n')
	return inputSb.String(), outSb.String()
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
	exp := strings.TrimSpace(expected)
	if got != exp {
		return fmt.Errorf("expected '%s' got '%s'", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

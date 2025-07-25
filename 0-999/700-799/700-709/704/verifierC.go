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

const mod = 1000000007

func solveC(n, m int, clauses [][]int) string {
	count := 0
	total := 1 << m
	for mask := 0; mask < total; mask++ {
		xorVal := 0
		for _, cl := range clauses {
			orVal := 0
			for _, lit := range cl {
				v := lit
				val := 0
				if v > 0 {
					if mask&(1<<(v-1)) != 0 {
						val = 1
					}
				} else {
					if mask&(1<<(-v-1)) == 0 {
						val = 1
					}
				}
				if val == 1 {
					orVal = 1
					break
				}
			}
			xorVal ^= orVal
		}
		if xorVal == 1 {
			count++
		}
	}
	return fmt.Sprintf("%d", count%mod)
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(3) + 1
	m := rng.Intn(4) + 1
	counts := make([]int, m+1)
	clauses := make([][]int, n)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for i := 0; i < n; i++ {
		k := rng.Intn(2) + 1
		clause := make([]int, k)
		for j := 0; j < k; j++ {
			var v int
			attempts := 0
			for {
				attempts++
				v = rng.Intn(m) + 1
				if counts[v] < 2 || attempts > 10 {
					break
				}
			}
			counts[v]++
			if rng.Intn(2) == 0 {
				clause[j] = v
			} else {
				clause[j] = -v
			}
		}
		clauses[i] = clause
		sb.WriteString(fmt.Sprintf("%d", k))
		for _, lit := range clause {
			sb.WriteString(fmt.Sprintf(" %d", lit))
		}
		sb.WriteByte('\n')
	}
	expected := solveC(n, m, clauses)
	return sb.String(), expected
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
	outStr := strings.TrimSpace(out.String())
	if outStr != expected {
		return fmt.Errorf("expected %s got %s", expected, outStr)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
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

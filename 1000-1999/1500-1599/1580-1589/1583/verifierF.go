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

type caseF struct{ n, k int }

func genCase(rng *rand.Rand) caseF {
	n := rng.Intn(8) + 2
	k := rng.Intn(n-1) + 2
	return caseF{n, k}
}

func expected(tc caseF) []string {
	n, k := tc.n, tc.k
	c := 0
	power := 1
	for power < n {
		power *= k
		c++
	}
	labels := make([][]int, n)
	for i := 0; i < n; i++ {
		labels[i] = make([]int, c)
		x := i
		for j := c - 1; j >= 0; j-- {
			labels[i][j] = x % k
			x /= k
		}
	}
	res := []string{fmt.Sprintf("%d", c)}
	m := n * (n - 1) / 2
	colors := make([]int, 0, m)
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			color := 1
			for d := 0; d < c; d++ {
				if labels[i][d] != labels[j][d] {
					color = d + 1
					break
				}
			}
			colors = append(colors, color)
		}
	}
	var sb strings.Builder
	for i, v := range colors {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	res = append(res, sb.String())
	return res
}

func runCase(bin string, tc caseF) error {
	input := fmt.Sprintf("%d %d\n", tc.n, tc.k)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	fields := strings.Fields(strings.TrimSpace(out.String()))
	exp := expected(tc)
	need := 1 + tc.n*(tc.n-1)/2
	if len(fields) != need {
		return fmt.Errorf("expected %d numbers got %d", need, len(fields))
	}
	if fields[0] != exp[0] {
		return fmt.Errorf("expected c=%s got %s", exp[0], fields[0])
	}
	expColors := strings.Fields(exp[1])
	for i := 1; i < need; i++ {
		if fields[i] != expColors[i-1] {
			return fmt.Errorf("color %d expected %s got %s", i, expColors[i-1], fields[i])
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	bin := os.Args[1]
	for i := 0; i < 100; i++ {
		tc := genCase(rng)
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

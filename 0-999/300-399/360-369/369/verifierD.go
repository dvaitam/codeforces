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

type caseD struct {
	n, k   int
	p      []int
	input  string
	expect int
}

func computeD(n, k int, p []int) int {
	sufPos := make([]bool, n+3)
	sufSafe := make([]bool, n+3)
	for i := n; i >= 1; i-- {
		sufPos[i] = sufPos[i+1] || (p[i-1] > 0)
		sufSafe[i] = sufSafe[i+1] && (p[i-1] < 100)
	}
	stride := k + 2
	visited := make([]bool, (n+2)*stride)
	type state struct{ i, d int }
	curr := []state{{1, func() int {
		if 2 > n+1 {
			return n + 1
		}
		return 2
	}() - 1}}
	visited[1*stride+curr[0].d] = true
	count := 1
	for round := 0; round < k && len(curr) > 0; round++ {
		next := []state{}
		for _, st := range curr {
			i := st.i
			j := i + st.d
			if i > n || j > n {
				continue
			}
			kill1 := sufPos[j]
			surv1 := sufSafe[j]
			kill2 := p[i-1] > 0
			surv2 := p[i-1] < 100
			if kill1 && kill2 {
				ni, nd := j+1, 1
				idx := ni*stride + nd
				if !visited[idx] {
					visited[idx] = true
					next = append(next, state{ni, nd})
					count++
				}
			}
			if kill1 && surv2 {
				ni, nd := j, 1
				idx := ni*stride + nd
				if !visited[idx] {
					visited[idx] = true
					next = append(next, state{ni, nd})
					count++
				}
			}
			if surv1 && kill2 {
				ni, nd := i, st.d+1
				idx := ni*stride + nd
				if !visited[idx] {
					visited[idx] = true
					next = append(next, state{ni, nd})
					count++
				}
			}
		}
		curr = next
	}
	return count
}

func generateCase(rng *rand.Rand) caseD {
	n := rng.Intn(5) + 1
	k := rng.Intn(5) + 1
	p := make([]int, n)
	for i := 0; i < n; i++ {
		p[i] = rng.Intn(101)
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, k))
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", p[i]))
	}
	sb.WriteByte('\n')
	exp := computeD(n, k, p)
	return caseD{n, k, p, sb.String(), exp}
}

func runCase(bin string, c caseD) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(c.input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var got int
	if _, err := fmt.Fscan(strings.NewReader(out.String()), &got); err != nil {
		return fmt.Errorf("failed to parse output: %v\n%s", err, out.String())
	}
	if got != c.expect {
		return fmt.Errorf("expected %d got %d", c.expect, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		c := generateCase(rng)
		if err := runCase(bin, c); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, c.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

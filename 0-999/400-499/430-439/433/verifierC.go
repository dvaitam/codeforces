package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
)

type test struct {
	input  string
	output string
}

func solve(n, m int, a []int) string {
	neighbors := make([][]int, n+1)
	var total int64
	for i := 1; i < m; i++ {
		u := a[i-1]
		v := a[i]
		if u > v {
			total += int64(u - v)
		} else {
			total += int64(v - u)
		}
		neighbors[u] = append(neighbors[u], v)
		neighbors[v] = append(neighbors[v], u)
	}
	bestGain := int64(0)
	for p := 1; p <= n; p++ {
		nb := neighbors[p]
		if len(nb) == 0 {
			continue
		}
		var curr int64
		for _, t := range nb {
			if p > t {
				curr += int64(p - t)
			} else {
				curr += int64(t - p)
			}
		}
		sort.Ints(nb)
		med := nb[len(nb)/2]
		var newCost int64
		for _, t := range nb {
			if med > t {
				newCost += int64(med - t)
			} else {
				newCost += int64(t - med)
			}
		}
		gain := curr - newCost
		if gain > bestGain {
			bestGain = gain
		}
	}
	res := total - bestGain
	return fmt.Sprintf("%d", res)
}

func generateTests() []test {
	rand.Seed(3)
	var tests []test
	for i := 0; i < 100; i++ {
		n := rand.Intn(10) + 1
		m := rand.Intn(15) + 1
		a := make([]int, m)
		for j := 0; j < m; j++ {
			a[j] = rand.Intn(n) + 1
		}
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d\n", n, m)
		for j := 0; j < m; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", a[j])
		}
		sb.WriteByte('\n')
		out := solve(n, m, a)
		tests = append(tests, test{sb.String(), out})
	}
	return tests
}

func runBinary(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = new(bytes.Buffer)
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return strings.TrimSpace(stdout.String()), nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		got, err := runBinary(bin, t.input)
		if err != nil {
			fmt.Printf("Test %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != t.output {
			fmt.Printf("Test %d failed:\ninput:\n%sexpected: %s\ngot: %s\n", i+1, t.input, t.output, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

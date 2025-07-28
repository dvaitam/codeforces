package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

const mod = 998244353

func countUnfortunate(n, m, k int) int {
	positions := make([]int, n)
	for i := 0; i < n; i++ {
		positions[i] = i + 1
	}
	var result int
	var rec func(int, int, []int)
	rec = func(start, need int, cur []int) {
		if len(cur) == need {
			result = (result + countOrientations(cur, k, n)) % mod
			return
		}
		for i := start; i < n; i++ {
			cur = append(cur, positions[i])
			rec(i+1, need, cur)
			cur = cur[:len(cur)-1]
		}
	}
	rec(0, m, []int{})
	return result
}

func countOrientations(pos []int, k, n int) int {
	m := len(pos)
	total := 0
	maskMax := 1 << m
	for mask := 0; mask < maskMax; mask++ {
		intervals := make([][2]int, m)
		ok := true
		for i, p := range pos {
			if (mask>>i)&1 == 0 { // left
				l := p - k
				r := p
				if l < 1 || r > n {
					ok = false
					break
				}
				intervals[i] = [2]int{l, r}
			} else { // right
				l := p
				r := p + k
				if l < 1 || r > n {
					ok = false
					break
				}
				intervals[i] = [2]int{l, r}
			}
		}
		if !ok {
			continue
		}
		sort.Slice(intervals, func(i, j int) bool { return intervals[i][0] < intervals[j][0] })
		for i := 1; i < m; i++ {
			if intervals[i][0] <= intervals[i-1][1] {
				ok = false
				break
			}
		}
		if ok {
			total++
		}
	}
	return total
}

func genCaseF(rng *rand.Rand) (int, int, int) {
	n := rng.Intn(6) + 1
	m := rng.Intn(n) + 1
	k := rng.Intn(n) + 1
	return n, m, k
}

func runCandidate(bin, input string) (string, error) {
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
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n, m, k := genCaseF(rng)
		input := fmt.Sprintf("%d %d %d\n", n, m, k)
		expect := fmt.Sprintf("%d", countUnfortunate(n, m, k))
		out, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if out != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", i+1, expect, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

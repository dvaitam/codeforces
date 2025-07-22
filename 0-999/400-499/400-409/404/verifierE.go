package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type testCaseE struct{ s string }

func generateTestsE() []testCaseE {
	r := rand.New(rand.NewSource(5))
	tests := make([]testCaseE, 0, 100)
	for len(tests) < 100 {
		n := 1 + r.Intn(6)
		var b strings.Builder
		for i := 0; i < n; i++ {
			if r.Intn(2) == 0 {
				b.WriteByte('L')
			} else {
				b.WriteByte('R')
			}
		}
		tests = append(tests, testCaseE{s: b.String()})
	}
	return tests
}

func simulate(s string, finish int, obstacles map[int]bool) (bool, bool) {
	pos := 0
	visitedFinish := false
	for i := 0; i < len(s); i++ {
		move := 0
		if s[i] == 'L' {
			move = -1
		} else {
			move = 1
		}
		next := pos + move
		if obstacles[next] {
			continue
		}
		pos = next
		if pos == finish {
			if i != len(s)-1 || visitedFinish {
				return false, false
			}
			visitedFinish = true
		}
	}
	return visitedFinish && pos == finish, visitedFinish
}

func bruteForceE(s string) int64 {
	n := len(s)
	minObs := 1<<31 - 1
	var count int64
	posRange := 2*n + 1
	cells := make([]int, 0, posRange)
	for i := -n; i <= n; i++ {
		if i != 0 {
			cells = append(cells, i)
		}
	}
	totalCells := len(cells)
	for obsMask := 0; obsMask < 1<<uint(totalCells); obsMask++ {
		obstacles := make(map[int]bool)
		for j, c := range cells {
			if obsMask&(1<<uint(j)) != 0 {
				obstacles[c] = true
			}
		}
		for finish := -n; finish <= n; finish++ {
			if finish == 0 {
				continue
			}
			ok, visited := simulate(s, finish, obstacles)
			if ok {
				k := len(obstacles)
				if k < minObs {
					minObs = k
					count = 1
				} else if k == minObs {
					count++
				}
			} else if visited {
				// visited but not valid -> ignore
			}
		}
	}
	if minObs == 1<<31-1 {
		return 0
	}
	return count
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("run failed: %v\n%s", err, errb.String())
	}
	return out.String(), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTestsE()
	for i, t := range tests {
		out, err := runBinary(bin, t.s+"\n")
		if err != nil {
			fmt.Printf("test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		expect := bruteForceE(t.s)
		got := strings.TrimSpace(out)
		if fmt.Sprint(expect) != got {
			fmt.Printf("test %d failed: expected %d got %s\n", i+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

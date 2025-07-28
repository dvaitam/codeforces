package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type test struct {
	input    string
	expected string
}

type cell struct {
	weights []int
	set     map[int]struct{}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func addCell(c *cell, posCount, negCount []int, posPresent, negPresent []bool, rFirst, nFirst *int, limit int) {
	for _, w := range c.weights {
		if w > 0 {
			if w > limit {
				continue
			}
			if posCount[w] == 0 {
				posPresent[w] = true
			}
			posCount[w]++
			if w == *rFirst {
				for *rFirst <= limit && posPresent[*rFirst] {
					*rFirst++
				}
			}
		} else {
			idx := -w
			if idx > limit {
				continue
			}
			if negCount[idx] == 0 {
				negPresent[idx] = true
			}
			negCount[idx]++
			if idx == *nFirst {
				for *nFirst <= limit && negPresent[*nFirst] {
					*nFirst++
				}
			}
		}
	}
}

func removeCell(c *cell, posCount, negCount []int, posPresent, negPresent []bool, rFirst, nFirst *int) {
	for _, w := range c.weights {
		if w > 0 {
			if posCount[w] > 0 {
				posCount[w]--
				if posCount[w] == 0 {
					posPresent[w] = false
					if w < *rFirst {
						*rFirst = w
					}
				}
			}
		} else {
			idx := -w
			if negCount[idx] > 0 {
				negCount[idx]--
				if negCount[idx] == 0 {
					negPresent[idx] = false
					if idx < *nFirst {
						*nFirst = idx
					}
				}
			}
		}
	}
}

func solveCase(n, m, k, t int, studs [][3]int) string {
	trans := false
	if n > m {
		n, m = m, n
		trans = true
	}
	limit := t - 1
	cells := make([][]cell, n)
	for i := range cells {
		cells[i] = make([]cell, m)
	}
	for _, st := range studs {
		x, y, w := st[0], st[1], st[2]
		if trans {
			x, y = y, x
		}
		if abs(w) > limit {
			continue
		}
		c := &cells[x-1][y-1]
		if c.set == nil {
			c.set = make(map[int]struct{})
		}
		if _, ok := c.set[w]; !ok {
			c.set[w] = struct{}{}
			c.weights = append(c.weights, w)
		}
	}
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			cells[i][j].set = nil
		}
	}
	if t <= 1 {
		ans := 0
		minnm := n
		if m < n {
			minnm = m
		}
		for L := 1; L <= minnm; L++ {
			ans += (n - L + 1) * (m - L + 1)
		}
		return fmt.Sprintf("%d", ans)
	}
	posCount := make([]int, limit+1)
	negCount := make([]int, limit+1)
	posPresent := make([]bool, limit+1)
	negPresent := make([]bool, limit+1)
	minnm := n
	if m < n {
		minnm = m
	}
	ans := 0
	for L := 1; L <= minnm; L++ {
		for top := 0; top <= n-L; top++ {
			for i := 1; i <= limit; i++ {
				posCount[i] = 0
				negCount[i] = 0
				posPresent[i] = false
				negPresent[i] = false
			}
			rFirst, nFirst := 1, 1
			for i := 0; i < L; i++ {
				for j := 0; j < L; j++ {
					addCell(&cells[top+i][j], posCount, negCount, posPresent, negPresent, &rFirst, &nFirst, limit)
				}
			}
			teamSize := rFirst + nFirst - 1
			if teamSize >= t {
				ans++
			}
			for left := 1; left <= m-L; left++ {
				for i := 0; i < L; i++ {
					removeCell(&cells[top+i][left-1], posCount, negCount, posPresent, negPresent, &rFirst, &nFirst)
					addCell(&cells[top+i][left+L-1], posCount, negCount, posPresent, negPresent, &rFirst, &nFirst, limit)
				}
				if rFirst <= limit {
					for rFirst <= limit && posPresent[rFirst] {
						rFirst++
					}
				}
				if nFirst <= limit {
					for nFirst <= limit && negPresent[nFirst] {
						nFirst++
					}
				}
				teamSize = rFirst + nFirst - 1
				if teamSize >= t {
					ans++
				}
			}
		}
	}
	return fmt.Sprintf("%d", ans)
}

func generateTests() []test {
	rng := rand.New(rand.NewSource(47))
	var tests []test
	tests = append(tests, test{input: "1 1 0 1\n", expected: solveCase(1, 1, 0, 1, nil)})
	for len(tests) < 100 {
		n := rng.Intn(3) + 1
		m := rng.Intn(3) + 1
		k := rng.Intn(n*m + 1)
		tVal := rng.Intn(3) + 1
		studs := make([][3]int, k)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d %d %d\n", n, m, k, tVal))
		for i := 0; i < k; i++ {
			x := rng.Intn(n) + 1
			y := rng.Intn(m) + 1
			w := rng.Intn(5) - 2
			if w == 0 {
				w = 1
			}
			studs[i] = [3]int{x, y, w}
			sb.WriteString(fmt.Sprintf("%d %d %d\n", x, y, w))
		}
		tests = append(tests, test{input: sb.String(), expected: solveCase(n, m, k, tVal, studs)})
	}
	return tests
}

func runBinary(bin, input string) (string, error) {
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
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		out, err := runBinary(bin, t.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(t.expected) {
			fmt.Printf("Wrong answer on test %d\nInput:\n%sExpected:%s\nGot:%s\n", i+1, t.input, t.expected, out)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type dsu struct {
	parent []int
	xor    []int
}

func newDSU(n int) *dsu {
	p := make([]int, n)
	x := make([]int, n)
	for i := range p {
		p[i] = i
	}
	return &dsu{parent: p, xor: x}
}

func (d *dsu) find(x int) int {
	if d.parent[x] != x {
		root := d.find(d.parent[x])
		d.xor[x] ^= d.xor[d.parent[x]]
		d.parent[x] = root
	}
	return d.parent[x]
}

func (d *dsu) union(a, b, v int) {
	ra := d.find(a)
	rb := d.find(b)
	if ra == rb {
		return
	}
	d.parent[ra] = rb
	d.xor[ra] = d.xor[a] ^ d.xor[b] ^ v
}

type testCase struct {
	n int
	a [][]int
}

// solve embeds the logic from 1713E.go.
func solve(tc testCase) string {
	n := tc.n
	a := tc.a
	d := newDSU(n)
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if a[i][j] == a[j][i] {
				continue
			}
			want := 0
			if a[i][j] > a[j][i] {
				want = 1
			}
			if d.find(i) != d.find(j) {
				d.union(i, j, want)
			}
		}
	}
	orient := make([]int, n)
	for i := 0; i < n; i++ {
		d.find(i)
		orient[i] = d.xor[i]
	}
	var out strings.Builder
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			val := a[i][j]
			if i != j && (orient[i]^orient[j]) == 1 {
				val = a[j][i]
			}
			if j > 0 {
				out.WriteByte(' ')
			}
			out.WriteString(strconv.Itoa(val))
		}
		if i+1 < n {
			out.WriteByte('\n')
		}
	}
	return out.String()
}

// Embedded copy of testcasesE.txt.
const testcaseData = `
2 3 5 2 5
3 3 3 1 4 4 5 2 1 1
3 3 3 1 1 2 1 5 4 1
2 1 3 4 2
3 3 4 2 5 2 4 3 5 1
1 2
1 4
3 5 1 5 1 2 4 5 3 1
1 5
1 3
1 3
2 4 3 4 5
1 4
3 1 1 3 2 5 1 5 1 2
3 4 2 4 5 5 2 3 2 4
3 1 1 5 1 3 1 1 1 4
3 3 3 5 4 5 4 2 1 2
2 1 2 2 1
2 1 5 4 2
3 1 3 3 5 2 4 4 1 4
1 4
2 5 5 3 3
3 2 4 1 4 4 3 2 4 1
2 1 1 1 3
1 2
2 1 5 3 4
3 1 4 1 4 5 1 3 4 5
1 1
2 1 3 1 4
1 5
3 1 1 4 3 2 5 5 2 2
1 2
3 1 1 5 4 4 3 2 2 2
3 2 2 3 3 2 2 4 4 4
2 5 2 4 5
1 2
3 2 1 4 1 5 3 1 1 1
1 1
1 5
3 2 4 2 4 3 5 4 4 2
1 3
3 4 1 2 4 3 4 5 4 1
2 2 4 4 4
1 4
1 4
1 5
3 1 1 3 3 1 5 4 2 1
3 1 4 3 5 2 5 4 1 4
3 1 4 4 4 1 2 4 3 3
2 2 4 5 3
3 1 2 3 3 1 1 5 4 2
2 3 5 3 3
2 1 3 5 4
1 4
2 2 4 3 2
1 1
3 2 1 4 5 4 4 1 1 1
2 2 5 5 2
1 2
1 2
1 5
1 4
2 5 5 3 1
2 3 3 3 5
2 5 4 4 5
2 3 4 4 4
2 5 5 4 4
1 2
1 1
3 2 1 5 2 1 5 5 1 3
3 1 1 2 1 1 4 4 3 4
2 1 4 4 5
3 2 1 4 5 5 2 2 5 2
1 2
3 1 4 4 5 3 2 5 2 2
2 4 5 1 4
2 4 5 3 1
1 5
2 4 5 1 1
2 4 5 4 5
2 3 2 2 5
2 3 5 2 2
2 1 1 4 4
3 4 1 5 4 2 1 3 3 2
2 1 1 2 2
3 1 1 2 4 4 5 5 1 1
2 5 5 3 1
3 2 2 1 3 1 2 4 3 4
2 1 5 3 1
2 3 4 2 3
2 2 2 3 2
3 1 1 4 3 2 3 4 4 3
3 2 3 4 3 3 3 4 3 3
1 4
2 1 2 4 3
1 2
2 4 4 2 4
2 4 5 4 5
3 3 2 4 4 2 5 5 1 1
1 4
`

var expectedOutputs = []string{
	`3 2
5 5`,
	`3 3 1
4 4 5
2 1 1`,
	`3 1 1
3 2 4
5 1 1`,
	`1 3
4 2`,
	`3 4 2
5 2 4
3 5 1`,
	`2`,
	`4`,
	`1 5 1
5 4 2
3 5 1`,
	`5`,
	`3`,
	`3`,
	`3 3
5 4`,
	`4`,
	`2 1 3
2 5 1
5 1 2`,
	`4 4 2
5 5 2
3 2 4`,
	`1 5 1
5 3 1
1 1 4`,
	`3 3 3
5 4 5
4 1 2`,
	`1 2
2 1`,
	`5 1
2 4`,
	`1 1 3
3 5 2
4 1 4`,
	`4`,
	`5 3
5 3`,
	`4 1 2
4 4 3
2 4 1`,
	`1 1
1 3`,
	`2`,
	`5 1
3 4`,
	`4 1 1
4 5 1
3 4 5`,
	`1`,
	`3 1
4 1`,
	`5`,
	`3 1 1
4 3 2
5 2 2`,
	`2`,
	`2 1 1
5 4 4
3 2 2`,
	`3 2 2
3 3 2
2 4 4`,
	`5 2
4 5`,
	`2`,
	`1 2 2
4 1 5
3 1 1`,
	`1`,
	`5`,
	`2 1 4
2 4 3
5 4 2`,
	`3`,
	`2 1 4
5 2 3
5 4 1`,
	`4 2
4 4`,
	`4`,
	`4`,
	`5 1
3 3`,
	`3 5 1
5 5 3
2 4 1`,
	`3 4 1
5 1 2
4 1 4`,
	`4 4 1
4 5 2
2 4 3`,
	`4 2
5 3`,
	`1 2 3
3 3 1
5 4 2`,
	`2 3
5 3`,
	`1 3
5 4`,
	`4`,
	`4 2
3 2`,
	`1`,
	`3 2 1
4 5 4
4 1 1`,
	`5 2
5 2`,
	`2`,
	`2`,
	`5`,
	`4`,
	`5 3
3 1`,
	`3 2
3 5`,
	`5 2
4 5`,
	`2 3
4 4`,
	`5 3
4 4`,
	`2`,
	`1`,
	`5 2 1
5 5 1
5 5 3`,
	`5 1 1
2 1 1
4 3 4`,
	`1 3
4 5`,
	`3 2 1
4 5 5
2 5 2`,
	`2`,
	`4 4 1
5 4 3
2 5 2`,
	`2 4
5 1`,
	`1 4
5 3`,
	`5`,
	`1 4
5 1`,
	`5 4
5 4`,
	`3 2
2 5`,
	`2 3
5 2`,
	`1 1
4 4`,
	`4 3 1
5 4 2
1 3 2`,
	`1 1
2 2`,
	`5 1 1
2 4 4
5 5 1`,
	`5 3
3 1`,
	`1 2 2
4 3 1
2 3 4`,
	`5 1
3 1`,
	`5 2
4 3`,
	`2 2
3 2`,
	`4 3 1
4 5 4
3 2 3`,
	`4 2 3
3 4 3
3 4 3`,
	`4`,
	`1 2
4 3`,
	`2`,
	`4 1
2 4`,
	`4 4
5 5`,
	`2 3 3
4 4 2
5 5 1`,
	`4`,
}

func loadTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcaseData), "\n")
	tests := make([]testCase, 0, len(lines))
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 1 {
			return nil, fmt.Errorf("line %d: empty", i+1)
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("line %d: bad n: %v", i+1, err)
		}
		if len(fields) != 1+n*n {
			return nil, fmt.Errorf("line %d: expected %d matrix entries, got %d", i+1, n*n, len(fields)-1)
		}
		mat := make([][]int, n)
		for r := 0; r < n; r++ {
			row := make([]int, n)
			for c := 0; c < n; c++ {
				val, err := strconv.Atoi(fields[1+r*n+c])
				if err != nil {
					return nil, fmt.Errorf("line %d: bad value at %d,%d: %v", i+1, r, c, err)
				}
				row[c] = val
			}
			mat[r] = row
		}
		tests = append(tests, testCase{n: n, a: mat})
	}
	return tests, nil
}

func runCase(bin string, tc testCase, expected string) error {
	var input strings.Builder
	input.WriteString("1\n")
	input.WriteString(strconv.Itoa(tc.n))
	input.WriteByte('\n')
	for i := 0; i < tc.n; i++ {
		for j := 0; j < tc.n; j++ {
			if j > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(strconv.Itoa(tc.a[i][j]))
		}
		input.WriteByte('\n')
	}

	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input.String())
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	got := strings.TrimSpace(out.String())
	if got != strings.TrimSpace(expected) {
		return fmt.Errorf("expected:\n%s\ngot:\n%s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests, err := loadTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to load testcases:", err)
		os.Exit(1)
	}
	if len(tests) != len(expectedOutputs) {
		fmt.Fprintf(os.Stderr, "testcase/expected mismatch: %d vs %d\n", len(tests), len(expectedOutputs))
		os.Exit(1)
	}

	for i, tc := range tests {
		exp := solve(tc)
		if err := runCase(bin, tc, exp); err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

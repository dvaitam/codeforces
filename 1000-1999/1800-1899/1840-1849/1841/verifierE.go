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

type DSU struct {
	parent []int
	size   []int
}

func NewDSU(n int) *DSU {
	d := &DSU{make([]int, n), make([]int, n)}
	for i := range d.parent {
		d.parent[i] = i
	}
	return d
}

func (d *DSU) find(x int) int {
	if d.parent[x] != x {
		d.parent[x] = d.find(d.parent[x])
	}
	return d.parent[x]
}

func (d *DSU) union(x, y int, cur *int64) {
	rx := d.find(x)
	ry := d.find(y)
	if rx == ry {
		return
	}
	if d.size[rx] < d.size[ry] {
		rx, ry = ry, rx
	}
	*cur -= int64(d.size[rx]/2 + d.size[ry]/2)
	d.size[rx] += d.size[ry]
	d.parent[ry] = rx
	*cur += int64(d.size[rx] / 2)
}

type testCaseE struct {
	input    string
	expected string
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
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func solveE(n int, a []int, m int64) int64 {
	events := make([][]int, n+2)
	for i, v := range a {
		if v < n {
			events[v+1] = append(events[v+1], i)
		}
	}

	d := NewDSU(n)
	active := make([]bool, n)
	var curPairs int64
	var totalPairs int64

	for r := 1; r <= n; r++ {
		for _, c := range events[r] {
			active[c] = true
			d.parent[c] = c
			d.size[c] = 1
			if c > 0 && active[c-1] {
				d.union(c, c-1, &curPairs)
			}
			if c+1 < n && active[c+1] {
				d.union(c, c+1, &curPairs)
			}
		}
		totalPairs += curPairs
	}

	if totalPairs > m/2 {
		totalPairs = m / 2
	}
	return totalPairs
}

func generateCaseE(rng *rand.Rand) testCaseE {
	t := rng.Intn(3) + 1
	var in strings.Builder
	var out strings.Builder
	in.WriteString(fmt.Sprintf("%d\n", t))
	for j := 0; j < t; j++ {
		n := rng.Intn(5) + 1
		a := make([]int, n)
		in.WriteString(fmt.Sprintf("%d\n", n))
		for i := 0; i < n; i++ {
			if i > 0 {
				in.WriteByte(' ')
			}
			a[i] = rng.Intn(n + 1)
			in.WriteString(fmt.Sprintf("%d", a[i]))
		}
		in.WriteByte('\n')
		white := 0
		for _, v := range a {
			white += n - v
		}
		var mVal int64
		if white == 0 {
			mVal = 0
		} else {
			mVal = int64(rng.Intn(white + 1))
		}
		in.WriteString(fmt.Sprintf("%d\n", mVal))
		out.WriteString(fmt.Sprintf("%d\n", solveE(n, a, mVal)))
	}
	return testCaseE{input: in.String(), expected: out.String()}
}

func runCaseE(bin string, tc testCaseE) error {
	got, err := runCandidate(bin, tc.input)
	if err != nil {
		return err
	}
	got = strings.TrimSpace(got)
	exp := strings.TrimSpace(tc.expected)
	if got != exp {
		return fmt.Errorf("expected:\n%s\ngot:\n%s", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := []testCaseE{generateCaseE(rng)}
	for i := 0; i < 100; i++ {
		cases = append(cases, generateCaseE(rng))
	}
	for i, tc := range cases {
		if err := runCaseE(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

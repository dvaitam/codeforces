package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func solve(input string) string {
	reader := bufio.NewReader(strings.NewReader(input))
	var xs, ys, n int
	fmt.Fscan(reader, &xs, &ys)
	fmt.Fscan(reader, &n)
	x := make([]int, n)
	y := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &x[i], &y[i])
	}
	d := make([]int, n)
	for i := 0; i < n; i++ {
		dx := x[i] - xs
		dy := y[i] - ys
		d[i] = dx*dx + dy*dy
	}
	a := make([][]int, n)
	for i := range a {
		a[i] = make([]int, n)
	}
	for i := 0; i < n; i++ {
		for j := 0; j < i; j++ {
			dx := x[i] - x[j]
			dy := y[i] - y[j]
			a[i][j] = dx*dx + dy*dy
			a[j][i] = a[i][j]
		}
	}
	size := 1 << n
	const INF = 1 << 60
	f := make([]int, size)
	g := make([]int, size)
	for i := range f {
		f[i] = INF
	}
	f[0] = 0
	for mask := 0; mask < size; mask++ {
		if f[mask] == INF {
			continue
		}
		for i := 0; i < n; i++ {
			bit := 1 << i
			if mask&bit == 0 {
				m1 := mask | bit
				cost1 := f[mask] + 2*d[i]
				if cost1 < f[m1] {
					f[m1] = cost1
					g[m1] = mask
				}
				for j := i + 1; j < n; j++ {
					bitj := 1 << j
					if mask&bitj == 0 {
						m2 := mask | bit | bitj
						cost2 := f[mask] + d[i] + d[j] + a[i][j]
						if cost2 < f[m2] {
							f[m2] = cost2
							g[m2] = mask
						}
					}
				}
				break
			}
		}
	}
	full := (1 << n) - 1
	var buf bytes.Buffer
	fmt.Fprintln(&buf, f[full])
	mask := full
	var path []int
	for mask > 0 {
		path = append(path, 0)
		prev := g[mask]
		diff := mask ^ prev
		for i := 0; i < n; i++ {
			if diff&(1<<i) != 0 {
				path = append(path, i+1)
			}
		}
		mask = prev
	}
	path = append(path, 0)
	for _, v := range path {
		fmt.Fprint(&buf, v, " ")
	}
	fmt.Fprintln(&buf)
	return buf.String()
}

type test struct{ input, expected string }

func generateTests() []test {
	rand.Seed(7)
	var tests []test
	for len(tests) < 100 {
		xs := rand.Intn(11) - 5
		ys := rand.Intn(11) - 5
		n := rand.Intn(4) + 1
		used := map[[2]int]bool{}
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d\n%d\n", xs, ys, n)
		for i := 0; i < n; i++ {
			for {
				x := rand.Intn(11) - 5
				y := rand.Intn(11) - 5
				if !used[[2]int{x, y}] && !(x == xs && y == ys) {
					used[[2]int{x, y}] = true
					fmt.Fprintf(&sb, "%d %d\n", x, y)
					break
				}
			}
		}
		inp := sb.String()
		tests = append(tests, test{inp, solve(inp)})
	}
	return tests
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		got, err := runBinary(bin, t.input)
		if err != nil {
			fmt.Printf("Runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != strings.TrimSpace(t.expected) {
			fmt.Printf("Wrong answer on test %d\nInput:\n%sExpected:\n%s\nGot:\n%s\n", i+1, t.input, strings.TrimSpace(t.expected), got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

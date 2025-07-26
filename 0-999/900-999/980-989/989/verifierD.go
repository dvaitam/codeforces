package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
)

type testCase struct {
	in  string
	out string
}

func expected(n int, l, wmax float64, xs []float64, vs []int) string {
	idx := make([]int, n)
	for i := range idx {
		idx[i] = i
	}
	sort.Slice(idx, func(i, j int) bool { return xs[idx[i]] < xs[idx[j]] })
	ans := 0
	for ii := 0; ii < n; ii++ {
		i := idx[ii]
		for jj := ii + 1; jj < n; jj++ {
			j := idx[jj]
			if vs[i] == vs[j] {
				continue
			}
			vi := float64(vs[i])
			vj := float64(vs[j])
			denom := xs[i] - xs[j]
			if denom == 0 {
				continue
			}
			w := ((xs[j]+l/2)*vi - (xs[i]+l/2)*vj) / denom
			if math.Abs(w) > wmax {
				continue
			}
			di := vi + w
			dj := vj + w
			if di == 0 || dj == 0 {
				continue
			}
			t := -(xs[i] + l/2) / di
			if t <= 0 {
				continue
			}
			ans++
		}
	}
	return fmt.Sprintf("%d", ans)
}

func generate() []testCase {
	const T = 100
	rand.Seed(4)
	cases := make([]testCase, T)
	for t := 0; t < T; t++ {
		n := rand.Intn(6) + 2
		l := float64(rand.Intn(10) + 1)
		wmax := float64(rand.Intn(5) + 1)
		xs := make([]float64, n)
		vs := make([]int, n)
		pos := rand.Intn(10) - 5
		xs[0] = float64(pos)
		for i := 1; i < n; i++ {
			pos += int(l) + rand.Intn(5) + 1
			xs[i] = float64(pos)
		}
		for i := 0; i < n; i++ {
			if rand.Intn(2) == 0 {
				vs[i] = -1
			} else {
				vs[i] = 1
			}
		}
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %.0f %.0f\n", n, l, wmax)
		for i := 0; i < n; i++ {
			fmt.Fprintf(&sb, "%.0f %d\n", xs[i], vs[i])
		}
		cases[t] = testCase{
			in:  sb.String(),
			out: expected(n, l, wmax, xs, vs),
		}
	}
	return cases
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := generate()
	for idx, tc := range cases {
		got, err := run(bin, tc.in)
		if err != nil {
			fmt.Printf("case %d: runtime error: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != strings.TrimSpace(tc.out) {
			fmt.Printf("case %d failed\ninput:\n%s\nexpected:\n%s\n\ngot:\n%s\n", idx+1, tc.in, tc.out, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}

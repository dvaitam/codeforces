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

type Test struct {
	in  string
	out string
}

func compute(a []int64) string {
	n := len(a)
	if n == 1 {
		return "-1"
	}
	sort.Slice(a, func(i, j int) bool { return a[i] < a[j] })
	if n == 2 {
		x, y := a[0], a[1]
		d := y - x
		if d == 0 {
			return fmt.Sprintf("1\n%d", x)
		} else if d%2 == 0 {
			mid := x + d/2
			res := []int64{x - d, mid, y + d}
			sort.Slice(res, func(i, j int) bool { return res[i] < res[j] })
			return fmt.Sprintf("%d\n%d %d %d", len(res), res[0], res[1], res[2])
		}
		return fmt.Sprintf("2\n%d %d", x-d, y+d)
	}
	diffs := make([]int64, n-1)
	for i := 1; i < n; i++ {
		diffs[i-1] = a[i] - a[i-1]
	}
	dmin := diffs[0]
	for _, d := range diffs {
		if d < dmin {
			dmin = d
		}
	}
	cntMin, cntOther := 0, 0
	var otherD int64
	pos := -1
	for i, d := range diffs {
		if d == dmin {
			cntMin++
		} else {
			cntOther++
			otherD = d
			pos = i
		}
	}
	if cntOther == 0 {
		if dmin == 0 {
			return fmt.Sprintf("1\n%d", a[0])
		}
		return fmt.Sprintf("2\n%d %d", a[0]-dmin, a[n-1]+dmin)
	}
	if cntOther == 1 && otherD == 2*dmin {
		x := a[pos] + dmin
		return fmt.Sprintf("1\n%d", x)
	}
	return "0"
}

func genCase(r *rand.Rand) Test {
	n := r.Intn(6) + 2
	arr := make([]int64, n)
	for i := range arr {
		arr[i] = int64(r.Intn(20))
	}
	var sb strings.Builder
	fmt.Fprintln(&sb, n)
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprint(&sb, v)
	}
	sb.WriteByte('\n')
	out := compute(append([]int64(nil), arr...))
	return Test{sb.String(), out}
}

func runCase(bin string, t Test) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(t.in)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	expect := strings.TrimSpace(t.out)
	if got != expect {
		return fmt.Errorf("expected:\n%s\ngot:\n%s", expect, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(3))
	for i := 0; i < 25; i++ {
		tc := genCase(rng)
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

func abs64(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func minDiff(a, b []int64) int64 {
	sa, sb := int64(0), int64(0)
	for _, v := range a {
		sa += v
	}
	for _, v := range b {
		sb += v
	}
	best := abs64(sa - sb)
	base := sa - sb
	for _, x := range a {
		for _, y := range b {
			diff := abs64(base - 2*x + 2*y)
			if diff < best {
				best = diff
			}
		}
	}
	if len(a) >= 2 && len(b) >= 2 {
		asum := make([]int64, 0, len(a)*(len(a)-1)/2)
		for i := 0; i < len(a); i++ {
			for j := i + 1; j < len(a); j++ {
				asum = append(asum, a[i]+a[j])
			}
		}
		bsum := make([]int64, 0, len(b)*(len(b)-1)/2)
		for i := 0; i < len(b); i++ {
			for j := i + 1; j < len(b); j++ {
				bsum = append(bsum, b[i]+b[j])
			}
		}
		sort.Slice(asum, func(i, j int) bool { return asum[i] < asum[j] })
		sort.Slice(bsum, func(i, j int) bool { return bsum[i] < bsum[j] })
		ptr := 0
		for _, x := range asum {
			for ptr < len(bsum)-1 && abs64(base-2*x+2*bsum[ptr]) >= abs64(base-2*x+2*bsum[ptr+1]) {
				ptr++
			}
			diff := abs64(base - 2*x + 2*bsum[ptr])
			if diff < best {
				best = diff
			}
		}
	}
	return best
}

func runCase(bin, input string) (string, error) {
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
		return out.String(), fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return out.String(), nil
}

func verify(a []int64, b []int64, out string) error {
	reader := bufio.NewReader(strings.NewReader(out))
	var d int64
	if _, err := fmt.Fscan(reader, &d); err != nil {
		return fmt.Errorf("failed to read diff: %v", err)
	}
	var k int
	if _, err := fmt.Fscan(reader, &k); err != nil {
		return fmt.Errorf("failed to read k: %v", err)
	}
	if k < 0 || k > 2 {
		return fmt.Errorf("invalid k")
	}
	aa := append([]int64(nil), a...)
	bb := append([]int64(nil), b...)
	for i := 0; i < k; i++ {
		var x, y int
		if _, err := fmt.Fscan(reader, &x, &y); err != nil {
			return fmt.Errorf("failed to read swap %d: %v", i+1, err)
		}
		if x < 1 || x > len(aa) || y < 1 || y > len(bb) {
			return fmt.Errorf("swap %d indices out of range", i+1)
		}
		aa[x-1], bb[y-1] = bb[y-1], aa[x-1]
	}
	if _, err := fmt.Fscan(reader); err == nil {
		return fmt.Errorf("extra output")
	}
	sa, sb := int64(0), int64(0)
	for _, v := range aa {
		sa += v
	}
	for _, v := range bb {
		sb += v
	}
	diff := abs64(sa - sb)
	if diff != d {
		return fmt.Errorf("reported diff %d but got %d", d, diff)
	}
	expect := minDiff(a, b)
	if d != expect {
		return fmt.Errorf("expected diff %d got %d", expect, d)
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

	type test struct{ a, b []int64 }
	tests := []test{
		{[]int64{1, 2}, []int64{3}},
		{[]int64{1}, []int64{1}},
	}
	for i := 0; i < 100; i++ {
		n := rng.Intn(6) + 1
		m := rng.Intn(6) + 1
		aa := make([]int64, n)
		bb := make([]int64, m)
		for j := 0; j < n; j++ {
			aa[j] = rng.Int63n(20) - 10
		}
		for j := 0; j < m; j++ {
			bb[j] = rng.Int63n(20) - 10
		}
		tests = append(tests, test{aa, bb})
	}

	for i, tc := range tests {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", len(tc.a)))
		for j, v := range tc.a {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		sb.WriteByte('\n')
		sb.WriteString(fmt.Sprintf("%d\n", len(tc.b)))
		for j, v := range tc.b {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		sb.WriteByte('\n')
		out, err := runCase(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, sb.String())
			os.Exit(1)
		}
		if err := verify(tc.a, tc.b, out); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%soutput:\n%s", i+1, err, sb.String(), out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

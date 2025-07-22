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
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func computeExpected(n, m int, a []int64) int64 {
	maxK := m / n
	check := func(k int) bool {
		if k == 0 {
			return true
		}
		kk := int64(k)
		base := int64(n) * kk * (kk - 1) / 2
		for j := 0; j < n; j++ {
			req := kk*int64(j+1) + base
			if req > a[j] {
				return false
			}
		}
		return true
	}
	lo, hi := 0, maxK+1
	for lo+1 < hi {
		mid := (lo + hi) >> 1
		if check(mid) {
			lo = mid
		} else {
			hi = mid
		}
	}
	k := lo
	if k == 0 {
		return 0
	}
	kk := int64(k)
	kn := int64(n) * kk
	minTotal := kn * (kn + 1) / 2
	base := int64(n) * kk * (kk - 1) / 2
	bLeft := make([]int64, n)
	for j := 0; j < n; j++ {
		minj := kk*int64(j+1) + base
		bLeft[j] = a[j] - minj
	}
	var extra int64
	nextVal := int64(m) + 1
	for i := kn; i >= 1; i-- {
		idx := int((i - 1) % int64(n))
		maxMon := nextVal - 1
		maxBud := i + bLeft[idx]
		if maxMon > int64(m) {
			maxMon = int64(m)
		}
		var si int64
		if maxBud < maxMon {
			si = maxBud
		} else {
			si = maxMon
		}
		if si < i {
			si = i
		}
		extra += si - i
		bLeft[idx] -= si - i
		nextVal = si
	}
	return minTotal + extra
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(3) + 1
	m := rng.Intn(20) + n
	a := make([]int64, n)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, m)
	for i := 0; i < n; i++ {
		a[i] = int64(rng.Intn(m*n + 1))
		fmt.Fprintf(&sb, "%d ", a[i])
	}
	sb.WriteByte('\n')
	ans := computeExpected(n, m, a)
	return sb.String(), fmt.Sprintf("%d", ans)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		out, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if out != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

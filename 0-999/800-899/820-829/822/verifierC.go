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

type voucher struct {
	l, r int
	c    int
	len  int
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
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func solveCase(segs []voucher, x int) int64 {
	n := len(segs)
	sort.Slice(segs, func(i, j int) bool { return segs[i].l < segs[j].l })
	byR := make([]voucher, n)
	copy(byR, segs)
	sort.Slice(byR, func(i, j int) bool { return byR[i].r < byR[j].r })
	const INF int64 = 1<<63 - 1
	best := make([]int64, x+1)
	for i := range best {
		best[i] = INF
	}
	ans := INF
	j := 0
	for _, s := range segs {
		for j < n && byR[j].r < s.l {
			l := byR[j].len
			if l <= x && int64(byR[j].c) < best[l] {
				best[l] = int64(byR[j].c)
			}
			j++
		}
		other := x - s.len
		if other >= 0 && best[other] != INF {
			cost := int64(s.c) + best[other]
			if cost < ans {
				ans = cost
			}
		}
	}
	if ans == INF {
		return -1
	}
	return ans
}

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(6) + 2
	x := rng.Intn(10) + 1
	segs := make([]voucher, n)
	for i := 0; i < n; i++ {
		l := rng.Intn(20) + 1
		r := l + rng.Intn(20)
		c := rng.Intn(20) + 1
		segs[i] = voucher{l, r, c, r - l + 1}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, x))
	for _, v := range segs {
		sb.WriteString(fmt.Sprintf("%d %d %d\n", v.l, v.r, v.c))
	}
	expect := solveCase(segs, x)
	return sb.String(), fmt.Sprintf("%d", expect)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := genCase(rng)
		out, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

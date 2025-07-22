package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type interval struct{ l, r int }

func solveA(n, x int, segs []interval) int {
	cur := 1
	watched := 0
	for _, seg := range segs {
		for cur+x <= seg.l {
			cur += x
		}
		if cur < seg.l {
			watched += seg.l - cur
			cur = seg.l
		}
		if cur <= seg.r {
			watched += seg.r - cur + 1
			cur = seg.r + 1
		}
	}
	return watched
}

func parseCase(input string) (int, int, []interval) {
	in := bufio.NewReader(strings.NewReader(input))
	var n, x int
	fmt.Fscan(in, &n, &x)
	segs := make([]interval, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &segs[i].l, &segs[i].r)
	}
	return n, x, segs
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(50) + 1
	x := rng.Intn(100000) + 1
	cur := 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, x)
	for i := 0; i < n; i++ {
		cur += rng.Intn(1000)
		l := cur
		r := l + rng.Intn(1000)
		if r > 100000 {
			r = 100000
		}
		if l > 100000 {
			l = 100000
		}
		if r < l {
			r = l
		}
		fmt.Fprintf(&sb, "%d %d\n", l, r)
		cur = r + 1
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var tests []string
	// some fixed edge cases
	tests = append(tests, "1 3\n4 6\n")
	tests = append(tests, "2 2\n1 1\n3 3\n")
	for i := 0; i < 98; i++ {
		tests = append(tests, generateCase(rng))
	}
	for i, tc := range tests {
		n, x, segs := parseCase(tc)
		expect := solveA(n, x, segs)
		out, err := runBinary(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: runtime error: %v\ninput:\n%s", i+1, err, tc)
			os.Exit(1)
		}
		got := strings.TrimSpace(out)
		if got != fmt.Sprint(expect) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %s\ninput:\n%s", i+1, expect, got, tc)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

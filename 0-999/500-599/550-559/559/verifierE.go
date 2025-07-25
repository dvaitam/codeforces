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

type interval struct{ l, r int }

type light struct{ a, l int }

type testCase struct{ lights []light }

func addInterval(ints []interval, iv interval) []interval {
	n := len(ints)
	pos := sort.Search(n, func(i int) bool { return ints[i].l >= iv.l })
	ints = append(ints, interval{})
	copy(ints[pos+1:], ints[pos:])
	ints[pos] = iv

	res := make([]interval, 0, len(ints))
	for _, it := range ints {
		if len(res) == 0 || res[len(res)-1].r < it.l {
			res = append(res, it)
		} else if res[len(res)-1].r < it.r {
			res[len(res)-1].r = it.r
		}
	}
	return res
}

func length(ints []interval) int {
	sum := 0
	for _, it := range ints {
		sum += it.r - it.l
	}
	return sum
}

func copyIntervals(ints []interval) []interval {
	res := make([]interval, len(ints))
	copy(res, ints)
	return res
}

func solveCase(lights []light) int {
	sort.Slice(lights, func(i, j int) bool { return lights[i].a < lights[j].a })
	intervals := []interval{}
	for _, lt := range lights {
		north := interval{lt.a, lt.a + lt.l}
		south := interval{lt.a - lt.l, lt.a}
		intsNorth := addInterval(copyIntervals(intervals), north)
		lenNorth := length(intsNorth)
		intsSouth := addInterval(copyIntervals(intervals), south)
		lenSouth := length(intsSouth)
		if lenNorth >= lenSouth {
			intervals = intsNorth
		} else {
			intervals = intsSouth
		}
	}
	return length(intervals)
}

func runCase(bin string, tc testCase) error {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(tc.lights))
	for _, lt := range tc.lights {
		fmt.Fprintf(&sb, "%d %d\n", lt.a, lt.l)
	}
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var got int
	if _, err := fmt.Fscan(strings.NewReader(out.String()), &got); err != nil {
		return fmt.Errorf("cannot parse output: %v", err)
	}
	exp := solveCase(tc.lights)
	if got != exp {
		return fmt.Errorf("expected %d got %d", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(5)
	cases := make([]testCase, 100)
	for i := range cases {
		n := rand.Intn(10) + 1
		ls := make([]light, n)
		for j := range ls {
			ls[j] = light{rand.Intn(100), rand.Intn(20) + 1}
		}
		cases[i] = testCase{ls}
	}
	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

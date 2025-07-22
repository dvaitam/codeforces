package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type testCaseD struct {
	n   int
	pts [][2]int
}

func genTestsD() []testCaseD {
	rand.Seed(4)
	tests := make([]testCaseD, 100)
	for i := range tests {
		n := rand.Intn(8) + 1 //1..8
		pts := make([][2]int, n)
		for j := range pts {
			pts[j][0] = rand.Intn(7)
			pts[j][1] = rand.Intn(7)
		}
		tests[i] = testCaseD{n, pts}
	}
	return tests
}

func solveD(tc testCaseD) int64 {
	pts := tc.pts
	xsAtY := map[int][]int{}
	ysAtX := map[int][]int{}
	pointSet := map[int64]struct{}{}
	for _, p := range pts {
		x, y := p[0], p[1]
		xsAtY[y] = append(xsAtY[y], x)
		ysAtX[x] = append(ysAtX[x], y)
		key := (int64(x) << 32) | int64(y)
		pointSet[key] = struct{}{}
	}
	for y := range xsAtY {
		sortInts(xsAtY[y])
	}
	for x := range ysAtX {
		sortInts(ysAtX[x])
	}
	var result int64
	for _, p := range pts {
		x, y := p[0], p[1]
		xs := xsAtY[y]
		ys := ysAtX[x]
		if len(xs) < len(ys) {
			for _, x2 := range xs {
				if x2 <= x {
					continue
				}
				d := x2 - x
				key1 := (int64(x) << 32) | int64(y+d)
				key2 := (int64(x2) << 32) | int64(y+d)
				if _, ok := pointSet[key1]; ok {
					if _, ok2 := pointSet[key2]; ok2 {
						result++
					}
				}
			}
		} else {
			for _, y2 := range ys {
				if y2 <= y {
					continue
				}
				d := y2 - y
				key1 := (int64(x+d) << 32) | int64(y)
				key2 := (int64(x+d) << 32) | int64(y2)
				if _, ok := pointSet[key1]; ok {
					if _, ok2 := pointSet[key2]; ok2 {
						result++
					}
				}
			}
		}
	}
	return result
}

func sortInts(a []int) {
	if len(a) < 2 {
		return
	}
	for i := 1; i < len(a); i++ {
		v := a[i]
		j := i - 1
		for j >= 0 && a[j] > v {
			a[j+1] = a[j]
			j--
		}
		a[j+1] = v
	}
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
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTestsD()
	for i, tc := range tests {
		var input bytes.Buffer
		fmt.Fprintln(&input, tc.n)
		for _, p := range tc.pts {
			fmt.Fprintf(&input, "%d %d\n", p[0], p[1])
		}
		expect := solveD(tc)
		out, err := run(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		val, err := strconv.ParseInt(strings.TrimSpace(out), 10, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: non-integer output %q\n", i+1, out)
			os.Exit(1)
		}
		if val != expect {
			fmt.Fprintf(os.Stderr, "test %d: expected %d got %d\n", i+1, expect, val)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

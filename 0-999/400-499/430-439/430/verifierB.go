package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

var rawTestcases = []string{
	"3 2 2 1 2 2",
	"8 5 2 1 4 1 4 4 5 1 4",
	"5 3 3 1 2 1 1 3",
	"1 5 2 4",
	"1 3 2 2",
	"4 4 2 2 4 3 1",
	"7 2 1 2 1 2 2 1 2 2",
	"8 5 5 1 4 2 4 4 2 3 5",
	"6 2 2 1 1 2 2 1 2",
	"1 4 4 2",
	"3 3 1 1 3 3",
	"4 5 5 3 5 3 4",
	"5 2 2 1 1 2 1 2",
	"6 3 3 2 2 1 3 3 2",
	"8 2 1 1 1 2 1 1 2 1 2",
	"4 4 1 2 3 3 1",
	"3 3 2 3 1 3",
	"5 4 4 3 4 4 1 1",
	"5 5 3 4 2 3 1 3",
	"4 5 1 2 1 4 2",
	"1 3 2 3",
	"7 3 3 3 3 2 1 3 3 1",
	"7 4 4 1 3 2 2 1 3 1",
	"2 4 3 2 4",
	"5 3 1 3 1 3 1 3",
	"8 3 3 3 3 1 2 1 2 1 1",
	"7 3 2 1 3 2 2 3 2 1",
	"6 5 3 1 2 2 3 5 2",
	"6 5 2 3 1 4 5 3 5",
	"8 3 1 3 1 1 3 1 2 2 3",
	"5 4 3 3 1 3 2 4",
	"3 2 2 1 2 1",
	"7 3 1 2 1 3 3 2 1 3",
	"4 2 2 2 2 1 2",
	"5 2 1 2 1 1 2 1",
	"1 3 1 3",
	"7 3 1 2 1 3 1 1 3 1",
	"7 5 5 3 5 3 4 3 1 2",
	"6 2 1 1 2 2 1 1 2",
	"8 2 2 1 2 2 1 1 2 1 1",
	"6 2 2 1 2 1 2 1 2",
	"5 2 2 1 2 2 1 2",
	"2 2 1 1 1",
	"4 5 1 3 5 1 1",
	"1 2 2 2",
	"8 5 2 1 5 3 1 5 2 2 3",
	"5 2 2 1 1 2 1 1",
	"5 5 5 2 1 2 3 1",
	"8 5 5 3 5 4 5 4 1 4 3",
	"3 4 4 1 4 1",
	"1 4 2 2",
	"3 4 3 4 4 2",
	"2 3 2 1 1",
	"6 5 2 2 3 4 4 2 4",
	"6 4 2 1 1 3 2 2 3",
	"5 4 3 2 4 1 1 4",
	"3 3 2 2 1 3",
	"1 5 4 3",
	"7 3 3 3 1 3 1 2 3 1",
	"5 2 1 1 2 1 2 2",
	"7 3 2 2 1 3 2 1 1 2",
	"7 2 2 2 1 2 1 1 2 1",
	"1 3 2 1",
	"3 4 2 2 3 3",
	"5 5 2 5 3 4 4 1",
	"4 5 2 3 1 1 5",
	"1 4 2 1",
	"6 4 4 3 3 1 1 4 4",
	"6 4 4 3 4 1 4 4 2",
	"1 4 2 4",
	"7 4 2 4 2 3 1 4 4 3",
	"2 5 2 3 1",
	"7 3 3 2 2 1 1 3 1 2",
	"5 5 5 3 2 4 3 4",
	"3 5 5 1 3 5",
	"2 5 1 3 1",
	"8 2 1 1 1 2 2 1 1 2 2",
	"2 2 2 2 1",
	"3 4 3 3 2 4",
	"7 3 2 2 3 2 3 1 2 3",
	"4 2 2 2 2 1 2",
	"4 2 1 2 1 2 2",
	"3 3 1 3 2 2",
	"5 5 2 1 2 3 1 1",
	"4 5 3 4 1 2 1",
	"1 2 1 1",
	"8 5 3 3 1 5 2 1 2 4 2",
	"8 5 4 2 2 3 4 5 5 4 2",
	"8 4 3 4 1 2 1 1 4 3 4",
	"5 3 2 1 3 1 1 2",
	"3 2 2 2 1 1",
	"8 4 1 1 1 2 1 3 1 4 1",
	"4 2 2 1 2 1 2",
	"7 4 3 3 2 2 1 2 3 4",
	"1 4 4 2",
	"7 2 2 1 2 1 1 2 1 1",
	"2 5 5 3 1",
	"6 2 1 1 2 1 2 2 1",
	"5 2 2 2 1 2 1 2",
	"1 4 3 2",
}

type testcase struct {
	n   int
	k   int
	x   int
	arr []int
	in  string
}

func parseCases() []testcase {
	var cases []testcase
	for _, line := range rawTestcases {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 3 {
			continue
		}
		n, _ := strconv.Atoi(fields[0])
		k, _ := strconv.Atoi(fields[1])
		x, _ := strconv.Atoi(fields[2])
		arr := make([]int, n)
		for i := 0; i < n && 3+i < len(fields); i++ {
			arr[i], _ = strconv.Atoi(fields[3+i])
		}
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d %d\n", n, k, x)
		for i, v := range arr {
			if i > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", v)
		}
		sb.WriteByte('\n')
		cases = append(cases, testcase{n: n, k: k, x: x, arr: arr, in: sb.String()})
	}
	return cases
}

// Ball structure from solution
type Ball struct {
	color int
	orig  bool
}

func simulate(c []int, x, pos int) int {
	n := len(c)
	balls := make([]Ball, 0, n+1)
	for i := 0; i < pos; i++ {
		balls = append(balls, Ball{color: c[i], orig: true})
	}
	balls = append(balls, Ball{color: x, orig: false})
	for i := pos; i < n; i++ {
		balls = append(balls, Ball{color: c[i], orig: true})
	}
	for {
		m := len(balls)
		removed := false
		newBalls := make([]Ball, 0, m)
		i := 0
		for i < m {
			j := i + 1
			for j < m && balls[j].color == balls[i].color {
				j++
			}
			if j-i >= 3 {
				removed = true
			} else {
				newBalls = append(newBalls, balls[i:j]...)
			}
			i = j
		}
		if !removed {
			break
		}
		balls = newBalls
	}
	rem := 0
	for _, b := range balls {
		if b.orig {
			rem++
		}
	}
	return len(c) - rem
}

func solve(c []int, x int) int {
	n := len(c)
	ans := 0
	for pos := 0; pos <= n; pos++ {
		d := simulate(c, x, pos)
		if d > ans {
			ans = d
		}
	}
	return ans
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: verifierB <solution-binary>")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := parseCases()
	for i, tc := range cases {
		expect := solve(tc.arr, tc.x)
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(tc.in)
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Test %d: runtime error: %v\n%s", i+1, err, out)
			os.Exit(1)
		}
		gotStr := strings.TrimSpace(string(out))
		var got int
		fmt.Sscan(gotStr, &got)
		if got != expect {
			fmt.Printf("Test %d failed.\nInput:\n%sExpected: %d\nGot: %s\n", i+1, tc.in, expect, gotStr)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}

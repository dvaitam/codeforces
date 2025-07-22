package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

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
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesB.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcases: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 4 {
			fmt.Printf("invalid testcase on line %d\n", idx+1)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(fields[0])
		k, _ := strconv.Atoi(fields[1]) // number of colors (unused in solution)
		x, _ := strconv.Atoi(fields[2])
		if len(fields) != 3+n {
			fmt.Printf("invalid testcase on line %d\n", idx+1)
			os.Exit(1)
		}
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			arr[i], _ = strconv.Atoi(fields[3+i])
		}
		expect := solve(arr, x)
		var buf bytes.Buffer
		fmt.Fprintf(&buf, "%d %d %d\n", n, k, x)
		for i, v := range arr {
			if i > 0 {
				buf.WriteByte(' ')
			}
			fmt.Fprintf(&buf, "%d", v)
		}
		buf.WriteByte('\n')
		cmd := exec.Command(bin)
		cmd.Stdin = bytes.NewReader(buf.Bytes())
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Test %d: runtime error: %v\n", idx+1, err)
			os.Exit(1)
		}
		gotStr := strings.TrimSpace(string(out))
		var got int
		fmt.Sscan(gotStr, &got)
		if got != expect {
			fmt.Printf("Test %d failed.\nInput:\n%sExpected: %d\nGot: %s\n", idx+1, buf.String(), expect, gotStr)
			os.Exit(1)
		}
		idx++
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}

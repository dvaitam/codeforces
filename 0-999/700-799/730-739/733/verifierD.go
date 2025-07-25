package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

type rect struct {
	a, b, c int
	id      int
}

func expectedCase(input string) string {
	reader := bufio.NewReader(strings.NewReader(input))
	var n int
	fmt.Fscan(reader, &n)
	rects := make([]rect, n)
	for i := 0; i < n; i++ {
		var a, b, c int
		fmt.Fscan(reader, &a, &b, &c)
		if a < b {
			a, b = b, a
		}
		if a < c {
			a, c = c, a
		}
		if b < c {
			b, c = c, b
		}
		rects[i] = rect{a: a, b: b, c: c, id: i + 1}
	}
	best := 0
	p1 := 1
	for _, r := range rects {
		if r.c > best {
			best = r.c
			p1 = r.id
		}
	}
	sort.Slice(rects, func(i, j int) bool {
		if rects[i].a != rects[j].a {
			return rects[i].a > rects[j].a
		}
		if rects[i].b != rects[j].b {
			return rects[i].b > rects[j].b
		}
		return rects[i].c > rects[j].c
	})
	flag := false
	p2 := 0
	for i := 1; i < n; i++ {
		if rects[i].a == rects[i-1].a && rects[i].b == rects[i-1].b {
			sumC := rects[i].c + rects[i-1].c
			cand := sumC
			if cand > rects[i].b {
				cand = rects[i].b
			}
			if cand > best {
				best = cand
				p1 = rects[i].id
				p2 = rects[i-1].id
				flag = true
			}
		}
	}
	if !flag {
		return fmt.Sprintf("1\n%d\n", p1)
	}
	return fmt.Sprintf("2\n%d %d\n", p1, p2)
}

func runCase(exe string, input, expected string) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	exp := strings.TrimSpace(expected)
	if got != exp {
		return fmt.Errorf("expected %q got %q", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	data, err := os.ReadFile("testcasesD.txt")
	if err != nil {
		fmt.Println("could not read testcasesD.txt:", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	for i := 0; i < t; i++ {
		if !scan.Scan() {
			fmt.Println("bad test file")
			os.Exit(1)
		}
		n, _ := strconv.Atoi(scan.Text())
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		rects := make([][3]int, n)
		for j := 0; j < n; j++ {
			scan.Scan()
			rects[j][0], _ = strconv.Atoi(scan.Text())
			scan.Scan()
			rects[j][1], _ = strconv.Atoi(scan.Text())
			scan.Scan()
			rects[j][2], _ = strconv.Atoi(scan.Text())
			sb.WriteString(fmt.Sprintf("%d %d %d\n", rects[j][0], rects[j][1], rects[j][2]))
		}
		exp := expectedCase(sb.String())
		if err := runCase(exe, sb.String(), exp); err != nil {
			fmt.Printf("case %d failed: %v\ninput:%s", i+1, err, sb.String())
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

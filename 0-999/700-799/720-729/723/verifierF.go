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

type Pair struct{ u, v int }

func expectedF(n, m int, edges []Pair, s, t, ds, dt int) (string, error) {
	parent := make([]int, n+1)
	for i := 1; i <= n; i++ {
		parent[i] = i
	}
	var find func(int) int
	find = func(x int) int {
		if parent[x] != x {
			parent[x] = find(parent[x])
		}
		return parent[x]
	}
	ans := make([]Pair, 0, n)
	for _, e := range edges {
		u, v := e.u, e.v
		if u != s && v != s && u != t && v != t {
			ru, rv := find(u), find(v)
			if ru != rv {
				parent[ru] = rv
				ans = append(ans, Pair{u, v})
			}
		}
	}
	a := make([]int, n+1)
	b := make([]int, n+1)
	flag := false
	for _, e := range edges {
		u, v := e.u, e.v
		if u == s && v == t {
			flag = true
		} else {
			if u == s && v != t {
				a[find(v)] = v
			}
			if v == s && u != t {
				a[find(u)] = u
			}
			if u == t && v != s {
				b[find(v)] = v
			}
			if v == t && u != s {
				b[find(u)] = u
			}
		}
	}
	x := make([]Pair, 0)
	y := make([]Pair, 0)
	z := make([]Pair, 0)
	for i := 1; i <= n; i++ {
		ri := find(i)
		if ri == i && i != s && i != t {
			av := a[i]
			bv := b[i]
			if av != 0 && bv != 0 {
				x = append(x, Pair{av, bv})
			} else if av != 0 {
				y = append(y, Pair{s, av})
			} else if bv != 0 {
				z = append(z, Pair{bv, t})
			} else {
				return "No\n", nil
			}
		}
	}
	ds -= len(y)
	dt -= len(z)
	ans = append(ans, y...)
	ans = append(ans, z...)
	if len(x) > 0 {
		ds--
		dt--
		last := x[len(x)-1]
		ans = append(ans, Pair{last.u, s})
		ans = append(ans, Pair{last.v, t})
		x = x[:len(x)-1]
		take := ds
		if take < 0 {
			take = 0
		}
		if take > len(x) {
			take = len(x)
		}
		ds -= take
		for i := 0; i < take; i++ {
			ans = append(ans, Pair{x[i].u, s})
		}
		rem := len(x) - take
		dt -= rem
		for i := take; i < len(x); i++ {
			ans = append(ans, Pair{x[i].v, t})
		}
	} else if flag {
		ans = append(ans, Pair{s, t})
		ds--
		dt--
	} else {
		return "No\n", nil
	}
	if ds < 0 || dt < 0 {
		return "No\n", nil
	}
	var sb strings.Builder
	sb.WriteString("Yes\n")
	for _, p := range ans {
		sb.WriteString(fmt.Sprintf("%d %d\n", p.u, p.v))
	}
	return sb.String(), nil
}

func runCase(exe, input, exp string) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	exp = strings.TrimSpace(exp)
	if got != exp {
		return fmt.Errorf("expected %q got %q", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	f, err := os.Open("testcasesF.txt")
	if err != nil {
		fmt.Println("could not read testcasesF.txt:", err)
		os.Exit(1)
	}
	defer f.Close()
	scan := bufio.NewScanner(f)
	scan.Split(bufio.ScanLines)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(strings.TrimSpace(scan.Text()))
	for caseNum := 0; caseNum < t; caseNum++ {
		if !scan.Scan() {
			fmt.Println("bad test file")
			os.Exit(1)
		}
		parts := strings.Fields(scan.Text())
		n, _ := strconv.Atoi(parts[0])
		m, _ := strconv.Atoi(parts[1])
		edges := make([]Pair, m)
		for i := 0; i < m; i++ {
			if !scan.Scan() {
				fmt.Println("bad test file")
				os.Exit(1)
			}
			p := strings.Fields(scan.Text())
			u, _ := strconv.Atoi(p[0])
			v, _ := strconv.Atoi(p[1])
			edges[i] = Pair{u, v}
		}
		if !scan.Scan() {
			fmt.Println("bad test file")
			os.Exit(1)
		}
		last := strings.Fields(scan.Text())
		s, _ := strconv.Atoi(last[0])
		tval, _ := strconv.Atoi(last[1])
		ds, _ := strconv.Atoi(last[2])
		dt, _ := strconv.Atoi(last[3])
		expStr, err := expectedF(n, m, edges, s, tval, ds, dt)
		if err != nil {
			fmt.Println("failed to compute expected:", err)
			os.Exit(1)
		}
		input := fmt.Sprintf("%d %d\n", n, m)
		for _, e := range edges {
			input += fmt.Sprintf("%d %d\n", e.u, e.v)
		}
		input += fmt.Sprintf("%d %d %d %d\n", s, tval, ds, dt)
		if err := runCase(exe, input, expStr); err != nil {
			fmt.Printf("case %d failed: %v\ninput:\n%s", caseNum+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Case struct {
	p         []int
	bestScore int
}

func computeQ(p []int) ([]int, int) {
	n := len(p)
	ev := make([][2]int, 0)
	od := make([][2]int, 0)
	pos1 := 0
	for i, v := range p {
		if v == 1 {
			pos1 = i
		}
		if i%2 == 0 {
			ev = append(ev, [2]int{v, i})
		} else {
			od = append(od, [2]int{v, i})
		}
	}
	sort.Slice(ev, func(i, j int) bool { return ev[i][0] < ev[j][0] })
	sort.Slice(od, func(i, j int) bool { return od[i][0] < od[j][0] })
	b := make([]int, n)
	cur := n
	if pos1%2 == 1 {
		for _, x := range ev {
			b[x[1]] = cur
			cur--
		}
		for _, x := range od {
			b[x[1]] = cur
			cur--
		}
	} else {
		for _, x := range od {
			b[x[1]] = cur
			cur--
		}
		for _, x := range ev {
			b[x[1]] = cur
			cur--
		}
	}
	return b, score(p, b)
}

func score(p, q []int) int {
	n := len(p)
	a := make([]int, n)
	for i := 0; i < n; i++ {
		a[i] = p[i] + q[i]
	}
	cnt := 0
	for i := 1; i < n-1; i++ {
		if a[i-1] < a[i] && a[i] > a[i+1] {
			cnt++
		}
	}
	return cnt
}

func genCases(n int) []Case {
	rand.Seed(time.Now().UnixNano())
	cs := make([]Case, n)
	for i := 0; i < n; i++ {
		size := rand.Intn(8) + 4
		if size%2 == 1 {
			size++
		}
		p := rand.Perm(size)
		for j := 0; j < size; j++ {
			p[j]++
		}
		_, sc := computeQ(append([]int(nil), p...))
		cs[i] = Case{p, sc}
	}
	return cs
}

func buildInput(cs []Case) string {
	var sb strings.Builder
	fmt.Fprintln(&sb, len(cs))
	for _, c := range cs {
		fmt.Fprintln(&sb, len(c.p))
		for j, v := range c.p {
			if j > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprint(&sb, v)
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cs := genCases(100)
	input := buildInput(cs)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println("runtime error:", err)
		os.Exit(1)
	}
	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(lines) != len(cs) {
		fmt.Printf("expected %d lines got %d\n", len(cs), len(lines))
		os.Exit(1)
	}
	for idx, line := range lines {
		parts := strings.Fields(line)
		if len(parts) != len(cs[idx].p) {
			fmt.Printf("case %d: expected %d numbers got %d\n", idx+1, len(cs[idx].p), len(parts))
			os.Exit(1)
		}
		q := make([]int, len(parts))
		used := make([]bool, len(parts)+1)
		for i, pv := range parts {
			v, err := strconv.Atoi(pv)
			if err != nil || v < 1 || v > len(parts) {
				fmt.Printf("case %d: invalid integer %s\n", idx+1, pv)
				os.Exit(1)
			}
			if used[v] {
				fmt.Printf("case %d: repeated value %d\n", idx+1, v)
				os.Exit(1)
			}
			used[v] = true
			q[i] = v
		}
		sc := score(cs[idx].p, q)
		if sc != cs[idx].bestScore {
			fmt.Printf("case %d: wrong score expected %d got %d\n", idx+1, cs[idx].bestScore, sc)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

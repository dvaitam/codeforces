package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

type Bus struct {
	s, f int
	t    int
	id   int
}

type Query struct {
	l, r, b   int
	id        int
	low, high int
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(46)
	var tests []string
	for t := 0; t < 100; t++ {
		n := rand.Intn(5) + 1
		m := rand.Intn(5) + 1
		buses := make([]Bus, n)
		usedT := make(map[int]bool)
		for i := 0; i < n; i++ {
			s := rand.Intn(19) + 1
			f := rand.Intn(20-s) + s + 1
			var ti int
			for {
				ti = rand.Intn(100) + 1
				if !usedT[ti] {
					usedT[ti] = true
					break
				}
			}
			buses[i] = Bus{s, f, ti, i + 1}
		}
		people := make([]Query, m)
		for i := 0; i < m; i++ {
			l := rand.Intn(19) + 1
			r := rand.Intn(20-l) + l + 1
			b := rand.Intn(100) + 1
			people[i] = Query{l, r, b, i, 0, n - 1}
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
		for _, b := range buses {
			sb.WriteString(fmt.Sprintf("%d %d %d\n", b.s, b.f, b.t))
		}
		for _, q := range people {
			sb.WriteString(fmt.Sprintf("%d %d %d\n", q.l, q.r, q.b))
		}
		tests = append(tests, sb.String())
	}
	for i, input := range tests {
		expect := solveE(strings.NewReader(input))
		out, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("test %d: execution failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(expect) {
			fmt.Printf("test %d failed: expected %q got %q\n", i+1, strings.TrimSpace(expect), strings.TrimSpace(out))
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return out.String(), err
}

// --- solution for problem E (brute force oracle) ---

func solveE(r io.Reader) string {
	in := bufio.NewReader(r)
	var n, m int
	fmt.Fscan(in, &n, &m)
	buses := make([]Bus, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &buses[i].s, &buses[i].f, &buses[i].t)
		buses[i].id = i + 1
	}
	sort.Slice(buses, func(i, j int) bool { return buses[i].t < buses[j].t })
	var buf strings.Builder
	for i := 0; i < m; i++ {
		var l, r, b int
		fmt.Fscan(in, &l, &r, &b)
		ans := -1
		for _, bus := range buses {
			if bus.t >= b && bus.s <= l && bus.f >= r {
				ans = bus.id
				break
			}
		}
		if i > 0 {
			buf.WriteByte(' ')
		}
		buf.WriteString(strconv.Itoa(ans))
	}
	buf.WriteByte('\n')
	return buf.String()
}

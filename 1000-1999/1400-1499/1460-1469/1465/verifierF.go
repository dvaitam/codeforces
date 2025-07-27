package main

import (
	"bytes"
	"container/list"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type testCaseF struct {
	n int
	p []int
}

func genTestsF() []testCaseF {
	rand.Seed(47)
	tests := make([]testCaseF, 100)
	for i := range tests {
		n := rand.Intn(4) + 3 // 3..6
		perm := rand.Perm(n)
		for j := range perm {
			perm[j]++
		}
		tests[i] = testCaseF{n, perm}
	}
	tests = append(tests, testCaseF{3, []int{2, 3, 1}})
	tests = append(tests, testCaseF{3, []int{2, 1, 3}})
	tests = append(tests, testCaseF{3, []int{1, 2, 3}})
	return tests
}

func prodCycles(p []int) int {
	n := len(p)
	vis := make([]bool, n)
	prod := 1
	for i := 0; i < n; i++ {
		if !vis[i] {
			l := 0
			j := i
			for !vis[j] {
				vis[j] = true
				j = p[j] - 1
				l++
			}
			prod *= l
		}
	}
	return prod
}

func solveF(tc testCaseF) (int, int) {
	start := make([]int, tc.n)
	copy(start, tc.p)
	type node struct {
		perm []int
		d    int
	}
	m := make(map[string]int)
	q := list.New()
	key := fmt.Sprint(start)
	m[key] = 0
	q.PushBack(node{start, 0})
	bestProd := prodCycles(start)
	bestDist := 0
	for q.Len() > 0 {
		e := q.Front()
		q.Remove(e)
		cur := e.Value.(node)
		p := cur.perm
		prod := prodCycles(p)
		if prod > bestProd {
			bestProd = prod
			bestDist = cur.d
		} else if prod == bestProd && cur.d < bestDist {
			bestDist = cur.d
		}
		for i := 0; i < tc.n; i++ {
			for j := i + 1; j < tc.n; j++ {
				np := make([]int, tc.n)
				copy(np, p)
				np[i], np[j] = np[j], np[i]
				k := fmt.Sprint(np)
				if _, ok := m[k]; !ok {
					m[k] = cur.d + 1
					q.PushBack(node{np, cur.d + 1})
				}
			}
		}
	}
	return bestProd % 1000000007, bestDist
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	cmd.Env = append(os.Environ(), "GOMAXPROCS=1")
	if err := cmd.Run(); err != nil {
		return out.String(), fmt.Errorf("runtime error: %v", err)
	}
	return out.String(), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTestsF()
	for i, tc := range tests {
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d\n", 1)
		fmt.Fprintf(&sb, "%d\n", tc.n)
		for j, v := range tc.p {
			if j > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", v)
		}
		sb.WriteByte('\n')
		input := sb.String()
		prod, dist := solveF(tc)
		exp := fmt.Sprintf("%d %d", prod, dist)
		got, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("test %d: %v\noutput:\n%s", i+1, err, got)
			os.Exit(1)
		}
		got = strings.TrimSpace(got)
		if got != exp {
			fmt.Printf("test %d failed: expected %q got %q\ninput:\n%s", i+1, exp, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

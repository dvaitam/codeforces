package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func solveE(n int, edges []struct {
	p int
	s string
}, t string) string {
	children := make([][]int, n+1)
	strs := make([][]byte, n+1)
	for i, e := range edges {
		v := i + 2
		children[e.p] = append(children[e.p], v)
		strs[v] = []byte(e.s)
	}
	tb := []byte(t)
	m := len(tb)
	pi := make([]int, m)
	for i := 1; i < m; i++ {
		j := pi[i-1]
		for j > 0 && tb[j] != tb[i] {
			j = pi[j-1]
		}
		if tb[j] == tb[i] {
			j++
		}
		pi[i] = j
	}
	var ans int64
	type frame struct {
		u                        int
		curBefore, curAfter, idx int
	}
	stack := []frame{{u: 1, idx: -1}}
	for len(stack) > 0 {
		fr := &stack[len(stack)-1]
		if fr.idx == -1 {
			cur := fr.curBefore
			if fr.u > 1 {
				for _, c := range strs[fr.u] {
					for cur > 0 && (cur == m || tb[cur] != c) {
						cur = pi[cur-1]
					}
					if cur < m && tb[cur] == c {
						cur++
					}
					if cur == m {
						ans++
					}
				}
			}
			fr.curAfter = cur
			fr.idx = 0
		} else if fr.idx < len(children[fr.u]) {
			v := children[fr.u][fr.idx]
			fr.idx++
			stack = append(stack, frame{u: v, curBefore: fr.curAfter, idx: -1})
		} else {
			stack = stack[:len(stack)-1]
		}
	}
	return fmt.Sprintf("%d\n", ans)
}

func randString(rng *rand.Rand, l int) string {
	letters := []byte("abc")
	b := make([]byte, l)
	for i := range b {
		b[i] = letters[rng.Intn(len(letters))]
	}
	return string(b)
}

func genCaseE(rng *rand.Rand) (string, string) {
	n := rng.Intn(6) + 1
	edges := make([]struct {
		p int
		s string
	}, n-1)
	for i := 0; i < n-1; i++ {
		edges[i].p = rng.Intn(i+1) + 1
		edges[i].s = randString(rng, rng.Intn(3)+1)
	}
	t := randString(rng, rng.Intn(3)+1)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for _, e := range edges {
		fmt.Fprintf(&sb, "%d %s\n", e.p, e.s)
	}
	fmt.Fprintf(&sb, "%s\n", t)
	input := sb.String()
	return input, solveE(n, edges, t)
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func main() {
	rand.Seed(time.Now().UnixNano())
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	for i := 0; i < 100; i++ {
		in, expect := genCaseE(rand.New(rand.NewSource(time.Now().UnixNano() + int64(i))))
		got, err := run(bin, in)
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\ninput:\n%soutput:\n%s", i+1, err, in, got)
			return
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Printf("test %d failed\ninput:\n%sexpected:\n%sbut got:\n%s", i+1, in, expect, got)
			return
		}
	}
	fmt.Println("All tests passed")
}

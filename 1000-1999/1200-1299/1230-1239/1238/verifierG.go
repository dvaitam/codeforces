package main

import (
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

// Embedded correct solver from the accepted solution

type Friend struct {
	t int64
	a int64
	b int64
}

type TNode struct {
	l, r int
	pr   uint32
	key  int64
	cnt  int64
	sz   int64
	sum  int64
}

var tnodes []TNode
var tseed uint64 = 88172645463393265

func trnd() uint32 {
	tseed ^= tseed << 7
	tseed ^= tseed >> 9
	return uint32(tseed)
}

func tnewNode(key, cnt int64) int {
	tnodes = append(tnodes, TNode{
		key: key,
		cnt: cnt,
		sz:  cnt,
		sum: key * cnt,
		pr:  trnd(),
	})
	return len(tnodes) - 1
}

func tpull(x int) {
	if x == 0 {
		return
	}
	n := &tnodes[x]
	n.sz = n.cnt + tnodes[n.l].sz + tnodes[n.r].sz
	n.sum = n.key*n.cnt + tnodes[n.l].sum + tnodes[n.r].sum
}

func trotateRight(x int) int {
	y := tnodes[x].l
	tnodes[x].l = tnodes[y].r
	tnodes[y].r = x
	tpull(x)
	tpull(y)
	return y
}

func trotateLeft(x int) int {
	y := tnodes[x].r
	tnodes[x].r = tnodes[y].l
	tnodes[y].l = x
	tpull(x)
	tpull(y)
	return y
}

func tinsert(x int, key, cnt int64) int {
	if x == 0 {
		return tnewNode(key, cnt)
	}
	if key == tnodes[x].key {
		tnodes[x].cnt += cnt
		tpull(x)
		return x
	}
	if key < tnodes[x].key {
		tnodes[x].l = tinsert(tnodes[x].l, key, cnt)
		if tnodes[x].l != 0 && tnodes[tnodes[x].l].pr > tnodes[x].pr {
			x = trotateRight(x)
		}
	} else {
		tnodes[x].r = tinsert(tnodes[x].r, key, cnt)
		if tnodes[x].r != 0 && tnodes[tnodes[x].r].pr > tnodes[x].pr {
			x = trotateLeft(x)
		}
	}
	tpull(x)
	return x
}

func tremoveSmallest(x int, k int64) (int, int64) {
	if x == 0 || k == 0 {
		return x, 0
	}
	l := tnodes[x].l
	if tnodes[l].sz >= k {
		nl, s := tremoveSmallest(l, k)
		tnodes[x].l = nl
		tpull(x)
		return x, s
	}
	s := tnodes[l].sum
	k -= tnodes[l].sz
	tnodes[x].l = 0
	if tnodes[x].cnt > k {
		s += tnodes[x].key * k
		tnodes[x].cnt -= k
		tpull(x)
		return x, s
	}
	s += tnodes[x].key * tnodes[x].cnt
	k -= tnodes[x].cnt
	nr, s2 := tremoveSmallest(tnodes[x].r, k)
	return nr, s + s2
}

func tremoveLargest(x int, k int64) (int, int64) {
	if x == 0 || k == 0 {
		return x, 0
	}
	r := tnodes[x].r
	if tnodes[r].sz >= k {
		nr, s := tremoveLargest(r, k)
		tnodes[x].r = nr
		tpull(x)
		return x, s
	}
	s := tnodes[r].sum
	k -= tnodes[r].sz
	tnodes[x].r = 0
	if tnodes[x].cnt > k {
		s += tnodes[x].key * k
		tnodes[x].cnt -= k
		tpull(x)
		return x, s
	}
	s += tnodes[x].key * tnodes[x].cnt
	k -= tnodes[x].cnt
	nl, s2 := tremoveLargest(tnodes[x].l, k)
	return nl, s + s2
}

func oracleSolve(input string) string {
	data := []byte(input)
	p := 0
	nextInt := func() int64 {
		for p < len(data) && (data[p] < '0' || data[p] > '9') && data[p] != '-' {
			p++
		}
		neg := false
		if p < len(data) && data[p] == '-' {
			neg = true
			p++
		}
		var v int64
		for p < len(data) && data[p] >= '0' && data[p] <= '9' {
			v = v*10 + int64(data[p]-'0')
			p++
		}
		if neg {
			v = -v
		}
		return v
	}

	q := int(nextInt())
	// Reset treap state
	tseed = 88172645463393265
	tnodes = make([]TNode, 1, 1200005)
	out := make([]byte, 0, q*4)

	for ; q > 0; q-- {
		n := int(nextInt())
		m := nextInt()
		c := nextInt()
		c0 := nextInt()

		friends := make([]Friend, n)
		for i := 0; i < n; i++ {
			t := nextInt()
			a := nextInt()
			b := nextInt()
			friends[i] = Friend{t: t, a: a, b: b}
		}

		if n > 1 {
			sort.Slice(friends, func(i, j int) bool {
				return friends[i].t < friends[j].t
			})
		}

		times := make([]int64, 0, n)
		bounds := make([]int, 0, n+1)
		for i := 0; i < n; {
			bounds = append(bounds, i)
			t := friends[i].t
			j := i + 1
			for j < n && friends[j].t == t {
				j++
			}
			times = append(times, t)
			i = j
		}
		bounds = append(bounds, n)

		k := len(times)
		g := make([]int64, k+1)
		prev := int64(0)
		for i := 0; i < k; i++ {
			g[i] = times[i] - prev
			prev = times[i]
		}
		g[k] = m - prev

		L := make([]int64, k+1)
		if g[k] < c {
			L[k] = g[k]
		} else {
			L[k] = c
		}
		for i := k - 1; i >= 0; i-- {
			x := g[i] + L[i+1]
			if x > c {
				x = c
			}
			L[i] = x
		}

		root := 0
		root = tinsert(root, 0, c0)
		if tnodes[root].sz > L[0] {
			root, _ = tremoveLargest(root, tnodes[root].sz-L[0])
		}

		ans := int64(0)
		feasible := true

		if tnodes[root].sz < g[0] {
			feasible = false
		} else {
			var s int64
			root, s = tremoveSmallest(root, g[0])
			ans += s
		}

		for i := 1; i <= k && feasible; i++ {
			for j := bounds[i-1]; j < bounds[i]; j++ {
				root = tinsert(root, friends[j].b, friends[j].a)
			}
			if tnodes[root].sz > L[i] {
				root, _ = tremoveLargest(root, tnodes[root].sz-L[i])
			}
			if tnodes[root].sz < g[i] {
				feasible = false
				break
			}
			var s int64
			root, s = tremoveSmallest(root, g[i])
			ans += s
		}

		if feasible {
			out = strconv.AppendInt(out, ans, 10)
		} else {
			out = append(out, '-', '1')
		}
		out = append(out, '\n')
	}

	return strings.TrimSpace(string(out))
}

// End of embedded solver

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

type TestEvent struct {
	t int64
	a int64
	b int64
}

func generateTests() []struct {
	n        int
	m, c, c0 int64
	events   []TestEvent
} {
	rng := rand.New(rand.NewSource(42))
	tests := make([]struct {
		n        int
		m, c, c0 int64
		events   []TestEvent
	}, 0, 100)
	for i := 0; i < 100; i++ {
		n := rng.Intn(4) + 1
		m := int64(rng.Intn(20) + 5)
		c := int64(rng.Intn(10) + 5)
		if c > m {
			c = m
		}
		c0 := int64(rng.Intn(int(c) + 1))
		events := make([]TestEvent, n)
		curT := int64(0)
		for j := 0; j < n; j++ {
			curT += int64(rng.Intn(5) + 1)
			a := int64(rng.Intn(5) + 1)
			b := int64(rng.Intn(5) + 1)
			events[j] = TestEvent{t: curT, a: a, b: b}
		}
		tests = append(tests, struct {
			n        int
			m, c, c0 int64
			events   []TestEvent
		}{n: n, m: m, c: c, c0: c0, events: events})
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		input := fmt.Sprintf("1\n%d %d %d %d\n", t.n, t.m, t.c, t.c0)
		for _, e := range t.events {
			input += fmt.Sprintf("%d %d %d\n", e.t, e.a, e.b)
		}
		want := oracleSolve(input)
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Printf("case %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != want {
			fmt.Printf("case %d failed: expected %s got %s\ninput:%s", i+1, want, got, input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))

	// suppress unused import warnings
	_ = io.Discard
}

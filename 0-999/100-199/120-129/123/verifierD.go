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

const ALPHA = 26

type state struct {
	next [ALPHA]int
	link int
	len  int
	occ  int
}

func solveD(s string) string {
	n := len(s)
	st := make([]state, 2*n)
	size, last := 1, 0
	st[0].link = -1
	var extend func(int)
	extend = func(c int) {
		cur := size
		size++
		st[cur].len = st[last].len + 1
		st[cur].occ = 1
		p := last
		for p != -1 && st[p].next[c] == 0 {
			st[p].next[c] = cur + 1
			p = st[p].link
		}
		if p == -1 {
			st[cur].link = 0
		} else {
			q := st[p].next[c] - 1
			if st[p].len+1 == st[q].len {
				st[cur].link = q
			} else {
				clone := size
				size++
				st[clone] = st[q]
				st[clone].len = st[p].len + 1
				st[clone].occ = 0
				for p != -1 && st[p].next[c] == q+1 {
					st[p].next[c] = clone + 1
					p = st[p].link
				}
				st[q].link = clone
				st[cur].link = clone
			}
		}
		last = cur
	}
	for i := 0; i < n; i++ {
		extend(int(s[i] - 'a'))
	}
	maxLen := n
	cntLen := make([]int, maxLen+1)
	for i := 0; i < size; i++ {
		cntLen[st[i].len]++
	}
	pos := make([]int, maxLen+1)
	for l := 1; l <= maxLen; l++ {
		pos[l] = pos[l-1] + cntLen[l-1]
	}
	order := make([]int, size)
	for i := 0; i < size; i++ {
		l := st[i].len
		order[pos[l]] = i
		pos[l]++
	}
	for i := size - 1; i >= 0; i-- {
		v := order[i]
		if st[v].link != -1 {
			st[st[v].link].occ += st[v].occ
		}
	}
	var ans uint64
	for i := 0; i < size; i++ {
		linkLen := 0
		if st[i].link != -1 {
			linkLen = st[st[i].link].len
		}
		delta := st[i].len - linkLen
		if delta > 0 {
			c := uint64(st[i].occ)
			ans += uint64(delta) * c * (c + 1) / 2
		}
	}
	return fmt.Sprintf("%d", ans)
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
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(20) + 1
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		b[i] = byte('a' + rng.Intn(3))
	}
	s := string(b)
	return s + "\n", s
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, s := genCase(rng)
		expect := solveD(s)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if got != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expect, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

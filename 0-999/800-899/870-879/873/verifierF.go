package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const alpha = 26

type state struct {
	next   [alpha]int
	link   int
	length int
	cnt    int64
}

var st []state
var size int
var last int
var pos []int

func initSAM(maxLen int) {
	st = make([]state, 2*maxLen+5)
	for i := range st {
		for j := 0; j < alpha; j++ {
			st[i].next[j] = -1
		}
		st[i].link = -1
	}
	size = 1
	last = 0
	st[0].link = -1
	st[0].length = 0
}

func saExtend(c int) {
	cur := size
	size++
	st[cur].length = st[last].length + 1
	for i := 0; i < alpha; i++ {
		st[cur].next[i] = -1
	}
	st[cur].cnt = 0
	p := last
	for p != -1 && st[p].next[c] == -1 {
		st[p].next[c] = cur
		p = st[p].link
	}
	if p == -1 {
		st[cur].link = 0
	} else {
		q := st[p].next[c]
		if st[p].length+1 == st[q].length {
			st[cur].link = q
		} else {
			clone := size
			size++
			st[clone] = st[q]
			st[clone].length = st[p].length + 1
			st[clone].cnt = 0
			for p != -1 && st[p].next[c] == q {
				st[p].next[c] = clone
				p = st[p].link
			}
			st[q].link = clone
			st[cur].link = clone
		}
	}
	last = cur
}

func expectedF(n int, s, t string) string {
	initSAM(n)
	pos = make([]int, n+1)
	for i := 0; i < n; i++ {
		saExtend(int(s[i] - 'a'))
		pos[i+1] = last
	}
	st = st[:size]
	maxLen := 0
	for i := 0; i < size; i++ {
		if st[i].length > maxLen {
			maxLen = st[i].length
		}
	}
	cntLen := make([]int, maxLen+1)
	for i := 0; i < size; i++ {
		cntLen[st[i].length]++
	}
	for i := 1; i <= maxLen; i++ {
		cntLen[i] += cntLen[i-1]
	}
	order := make([]int, size)
	for i := size - 1; i >= 0; i-- {
		l := st[i].length
		cntLen[l]--
		order[cntLen[l]] = i
	}
	for i := 1; i <= n; i++ {
		if t[i-1] == '0' {
			st[pos[i]].cnt++
		}
	}
	for i := size - 1; i > 0; i-- {
		v := order[i]
		p := st[v].link
		if p >= 0 {
			st[p].cnt += st[v].cnt
		}
	}
	var ans int64
	for i := 1; i < size; i++ {
		if st[i].cnt > 0 {
			val := int64(st[i].length) * st[i].cnt
			if val > ans {
				ans = val
			}
		}
	}
	return strconv.FormatInt(ans, 10)
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	err := cmd.Run()
	if err != nil {
		return errBuf.String(), fmt.Errorf("runtime error: %v", err)
	}
	return out.String(), nil
}

type TestF struct {
	n int
	s string
	t string
}

func (tc TestF) Input() string {
	return fmt.Sprintf("%d\n%s\n%s\n", tc.n, tc.s, tc.t)
}

func genTests() []TestF {
	rand.Seed(6)
	tests := make([]TestF, 0, 100)
	letters := []byte("abcdefghijklmnopqrstuvwxyz")
	for i := 0; i < 100; i++ {
		n := rand.Intn(20) + 1
		b := make([]byte, n)
		for j := range b {
			b[j] = letters[rand.Intn(26)]
		}
		forb := make([]byte, n)
		for j := range forb {
			if rand.Intn(2) == 0 {
				forb[j] = '0'
			} else {
				forb[j] = '1'
			}
		}
		tests = append(tests, TestF{n: n, s: string(b), t: string(forb)})
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierF.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	tests := genTests()
	for i, tc := range tests {
		exp := strings.TrimSpace(expectedF(tc.n, tc.s, tc.t))
		gotRaw, err := run(bin, tc.Input())
		if err != nil {
			fmt.Printf("runtime error on test %d: %v\n%s", i+1, err, gotRaw)
			os.Exit(1)
		}
		got := strings.TrimSpace(gotRaw)
		if got != exp {
			fmt.Printf("test %d failed\ninput:\n%sexpected: %s\ngot: %s\n", i+1, tc.Input(), exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

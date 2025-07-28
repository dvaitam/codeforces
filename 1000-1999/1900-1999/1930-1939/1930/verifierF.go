package main

import (
	"bytes"
	"fmt"
	"math/bits"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const maxBits = 22
const maxNodes = 1 << (maxBits + 1)

type trie struct {
	ch   [][2]int32
	mx   []int32
	size int32
}

func newTrie() *trie {
	t := &trie{
		ch:   make([][2]int32, maxNodes),
		mx:   make([]int32, maxNodes),
		size: 1,
	}
	for i := range t.mx {
		t.mx[i] = -1
	}
	return t
}

func (t *trie) reset() {
	t.size = 1
	t.ch[0][0], t.ch[0][1] = 0, 0
	t.mx[0] = -1
}

func (t *trie) add(v int) {
	node := int32(0)
	if int32(v) > t.mx[node] {
		t.mx[node] = int32(v)
	}
	for i := maxBits - 1; i >= 0; i-- {
		b := (v >> i) & 1
		if t.ch[node][b] == 0 {
			t.ch[node][b] = t.size
			t.ch[t.size][0], t.ch[t.size][1] = 0, 0
			t.mx[t.size] = -1
			t.size++
		}
		node = t.ch[node][b]
		if int32(v) > t.mx[node] {
			t.mx[node] = int32(v)
		}
	}
}

func (t *trie) query(mask int) int {
	if t.size == 1 {
		return 0
	}
	node := int32(0)
	res := 0
	for i := maxBits - 1; i >= 0; i-- {
		bit := (mask >> i) & 1
		c0 := t.ch[node][0]
		c1 := t.ch[node][1]
		if bit == 1 {
			if c1 != 0 {
				node = c1
				res |= 1 << i
			} else if c0 != 0 {
				node = c0
			} else {
				break
			}
		} else {
			if c0 == 0 && c1 == 0 {
				break
			} else if c0 == 0 {
				node = c1
			} else if c1 == 0 {
				node = c0
			} else {
				rem := mask & ((1 << i) - 1)
				v0 := int(t.mx[c0]) & rem
				v1 := int(t.mx[c1]) & rem
				if v0 >= v1 {
					node = c0
				} else {
					node = c1
				}
			}
		}
	}
	return res
}

func solveF(n, q int, e []int) []int {
	mask := (1 << bits.Len(uint(n-1))) - 1
	trieA := newTrie()
	trieB := newTrie()
	ans := 0
	last := 0
	res := make([]int, q)
	for i := 0; i < q; i++ {
		v := (e[i] + last) % n
		diff1 := trieA.query(mask ^ v)
		if diff1 > ans {
			ans = diff1
		}
		diff2 := trieB.query(v)
		if diff2 > ans {
			ans = diff2
		}
		trieA.add(v)
		trieB.add(mask ^ v)
		last = ans
		res[i] = ans
	}
	return res
}

func generateCaseF(rng *rand.Rand) (string, string) {
	n := rng.Intn(32) + 1
	q := rng.Intn(6) + 1
	e := make([]int, q)
	for i := range e {
		e[i] = rng.Intn(n)
	}
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d %d\n", n, q))
	for i, v := range e {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	input := sb.String()
	ans := solveF(n, q, e)
	var exp strings.Builder
	for i, v := range ans {
		if i > 0 {
			exp.WriteByte(' ')
		}
		exp.WriteString(strconv.Itoa(v))
	}
	exp.WriteByte('\n')
	return input, exp.String()
}

func runProg(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return out.String(), fmt.Errorf("%v", err)
	}
	return out.String(), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCaseF(rng)
		got, err := runProg(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: runtime error: %v\noutput:\n%s", i+1, err, got)
			os.Exit(1)
		}
		got = strings.TrimSpace(got)
		if got != strings.TrimSpace(exp) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %q got %q\ninput:\n%s", i+1, strings.TrimSpace(exp), got, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

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

type BIT struct {
	n    int
	tree []int
}

func NewBIT(n int) *BIT {
	return &BIT{n: n, tree: make([]int, n+1)}
}

func (b *BIT) Add(i, delta int) {
	for x := i; x <= b.n; x += x & -x {
		b.tree[x] += delta
	}
}

func (b *BIT) Sum(i int) int {
	s := 0
	for x := i; x > 0; x -= x & -x {
		s += b.tree[x]
	}
	return s
}

func (b *BIT) FindKth(k int) int {
	pos := 0
	bitMask := 1
	for bitMask<<1 <= b.n {
		bitMask <<= 1
	}
	for d := bitMask; d > 0; d >>= 1 {
		nxt := pos + d
		if nxt <= b.n && b.tree[nxt] < k {
			pos = nxt
			k -= b.tree[nxt]
		}
	}
	return pos + 1
}

func solveC(k int, s string, ops [][2]interface{}) string {
	m := len(s)
	total := k * m
	pos := make([][]int, 26)
	for i := 0; i < total; i++ {
		c := s[i%len(s)] - 'a'
		pos[c] = append(pos[c], i)
	}
	bits := make([]*BIT, 26)
	for c := 0; c < 26; c++ {
		bits[c] = NewBIT(len(pos[c]))
		for i := 1; i <= len(pos[c]); i++ {
			bits[c].Add(i, 1)
		}
	}
	deleted := make([]bool, total)
	for _, op := range ops {
		p := op[0].(int)
		c := op[1].(byte)
		idx := bits[c-'a'].FindKth(p) - 1
		g := pos[c-'a'][idx]
		deleted[g] = true
		bits[c-'a'].Add(idx+1, -1)
	}
	res := make([]byte, 0, total)
	for i := 0; i < total; i++ {
		if !deleted[i] {
			res = append(res, s[i%len(s)])
		}
	}
	return string(res)
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
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func randString(rng *rand.Rand, n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyz")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rng.Intn(len(letters))]
	}
	return string(b)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		k := rng.Intn(3) + 1
		s := randString(rng, rng.Intn(4)+1)
		total := k * len(s)
		nOps := rng.Intn(total/2 + 1)
		t := []byte(strings.Repeat(s, k))
		ops := make([][2]interface{}, nOps)
		for j := 0; j < nOps; j++ {
			c := t[rng.Intn(len(t))]
			// count occurrences
			cnt := 0
			for _, ch := range t {
				if ch == c {
					cnt++
				}
			}
			p := rng.Intn(cnt) + 1
			ops[j][0] = p
			ops[j][1] = c
			// delete
			occ := 0
			for idx := 0; idx < len(t); idx++ {
				if t[idx] == c {
					occ++
					if occ == p {
						t = append(t[:idx], t[idx+1:]...)
						break
					}
				}
			}
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n%s\n%d\n", k, s, nOps))
		for _, op := range ops {
			sb.WriteString(fmt.Sprintf("%d %c\n", op[0].(int), op[1].(byte)))
		}
		input := sb.String()
		expected := solveC(k, s, ops)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expected) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expected, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Bitset []uint64

func newBitset(words int) Bitset { return make(Bitset, words) }
func (b Bitset) clone() Bitset   { nb := make(Bitset, len(b)); copy(nb, b); return nb }
func (b Bitset) set(i int)       { b[i>>6] |= 1 << uint(i&63) }
func (b Bitset) get(i int) bool  { return (b[i>>6] & (1 << uint(i&63))) != 0 }
func (b Bitset) orShift(shift, words int) {
	w := shift >> 6
	off := uint(shift & 63)
	for i := words - 1; i >= 0; i-- {
		j := i - w
		if j < 0 {
			continue
		}
		v := b[j] << off
		if off != 0 && j-1 >= 0 {
			v |= b[j-1] >> (64 - off)
		}
		b[i] |= v
	}
}

func solveD(in string) string {
	reader := bufio.NewReader(strings.NewReader(in))
	var n, s int
	if _, err := fmt.Fscan(reader, &n, &s); err != nil {
		return ""
	}
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
	}
	p := 0
	for i := 1; i < n; i++ {
		if a[i] > a[p] {
			p = i
		}
	}
	if s < a[p] {
		return fmt.Sprintln(-1)
	}
	const S = 7
	m := (n + S - 1) / S
	bits := s + 1
	words := (bits + 63) >> 6
	f := make([]Bitset, m+1)
	f[0] = newBitset(words)
	f[0].set(0)
	flag := make([]bool, n)
	for i := 0; i < m; i++ {
		cur := f[i].clone()
		for j := 0; j < S; j++ {
			k := i*S + j
			if k >= n || k == p {
				continue
			}
			shift := a[k]
			if shift <= s {
				cur.orShift(shift, words)
			}
		}
		f[i+1] = cur
	}
	t := s - a[p]
	if t < 0 || t > s || !f[m].get(t) {
		return fmt.Sprintln(-1)
	}
	for i := m; i > 0; i-- {
		base := i * S
		for mask := 0; mask < 1<<S; mask++ {
			val := 0
			ok := true
			for j := 0; j < S; j++ {
				if mask>>j&1 == 1 {
					k := base - S + j
					if k < 0 || k >= n || k == p {
						ok = false
						break
					}
					val += a[k]
				}
			}
			if !ok || val > t {
				continue
			}
			if f[i-1].get(t - val) {
				for j := 0; j < S; j++ {
					if mask>>j&1 == 1 {
						k := base - S + j
						flag[k] = true
					}
				}
				t -= val
				break
			}
		}
	}
	flag[p] = true
	c := make([]int, n)
	adj := make([][]int, n)
	pool := make([]int, 0, n)
	for i := 0; i < n; i++ {
		if !flag[i] {
			pool = append(pool, i)
		} else {
			c[i] = a[i]
		}
	}
	sort.Slice(pool, func(i, j int) bool { return a[pool[i]] < a[pool[j]] })
	pool = append(pool, p)
	for i, x := range pool {
		c[x] = a[x]
		if i > 0 {
			par := pool[i-1]
			adj[x] = append(adj[x], par)
			c[x] -= a[par]
		}
	}
	sumC := 0
	for i := 0; i < n; i++ {
		sumC += c[i]
	}
	if sumC != s {
		return fmt.Sprintln(-1)
	}
	var sb strings.Builder
	for i := 0; i < n; i++ {
		if len(adj[i]) == 0 {
			fmt.Fprintf(&sb, "%d %d\n", c[i], 0)
		} else {
			fmt.Fprintf(&sb, "%d %d", c[i], len(adj[i]))
			for _, v := range adj[i] {
				fmt.Fprintf(&sb, " %d", v+1)
			}
			sb.WriteByte('\n')
		}
	}
	return sb.String()
}

func genTest(r *rand.Rand) string {
	n := r.Intn(5) + 1
	a := make([]int, n)
	sum := 0
	for i := 0; i < n; i++ {
		a[i] = r.Intn(5) + 1
		sum += a[i]
	}
	s := r.Intn(sum) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, s)
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(a[i]))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func runBinary(path, in string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, path)
	cmd.Stdin = strings.NewReader(in)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: verifierD <path-to-binary>")
		os.Exit(1)
	}
	bin := os.Args[1]
	r := rand.New(rand.NewSource(4))
	const tests = 100
	for i := 0; i < tests; i++ {
		in := genTest(r)
		expect := strings.TrimSpace(solveD(in))
		got, err := runBinary(bin, in)
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != expect {
			fmt.Printf("test %d failed\ninput:\n%sexpected:\n%s\ngot:\n%s\n", i+1, in, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

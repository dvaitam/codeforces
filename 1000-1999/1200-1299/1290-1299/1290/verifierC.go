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

type DSU struct {
	parent []int
	parity []int
	size   []int
	cnt    [][2]int
	forced []int
}

func NewDSU(n int) *DSU {
	d := &DSU{
		parent: make([]int, n),
		parity: make([]int, n),
		size:   make([]int, n),
		cnt:    make([][2]int, n),
		forced: make([]int, n),
	}
	for i := 0; i < n; i++ {
		d.parent[i] = i
		d.size[i] = 1
		d.forced[i] = -1
	}
	d.forced[0] = 0
	for i := 1; i < n; i++ {
		d.cnt[i][0] = 1
	}
	return d
}

func (d *DSU) find(x int) int {
	if d.parent[x] != x {
		px := d.parent[x]
		d.parent[x] = d.find(px)
		d.parity[x] ^= d.parity[px]
	}
	return d.parent[x]
}

func (d *DSU) cost(x int) int {
	if d.forced[x] == -1 {
		if d.cnt[x][0] < d.cnt[x][1] {
			return d.cnt[x][0]
		}
		return d.cnt[x][1]
	}
	if d.forced[x] == 0 {
		return d.cnt[x][1]
	}
	return d.cnt[x][0]
}

func (d *DSU) unite(a, b, w int, ans *int) {
	ra := d.find(a)
	rb := d.find(b)
	w ^= d.parity[a] ^ d.parity[b]
	if ra == rb {
		return
	}
	if d.size[ra] > d.size[rb] {
		ra, rb = rb, ra
	}
	*ans -= d.cost(ra)
	*ans -= d.cost(rb)
	d.parent[ra] = rb
	d.parity[ra] = w
	d.size[rb] += d.size[ra]
	c0 := d.cnt[ra][0]
	c1 := d.cnt[ra][1]
	if w == 1 {
		c0, c1 = c1, c0
	}
	d.cnt[rb][0] += c0
	d.cnt[rb][1] += c1
	if d.forced[ra] != -1 {
		val := d.forced[ra] ^ w
		if d.forced[rb] == -1 {
			d.forced[rb] = val
		}
	}
	*ans += d.cost(rb)
}

func solveC(n, k int, s string, pos [][]int) []int {
	dsu := NewDSU(k + 1)
	res := make([]int, n)
	cur := 0
	for i := 0; i < n; i++ {
		v := int('1' - s[i])
		if len(pos[i]) == 1 {
			dsu.unite(pos[i][0], 0, v, &cur)
		} else if len(pos[i]) == 2 {
			dsu.unite(pos[i][0], pos[i][1], v, &cur)
		}
		res[i] = cur
	}
	return res
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
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

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for t := 1; t <= 100; t++ {
		n := rng.Intn(8) + 1
		k := rng.Intn(4) + 1
		var sb strings.Builder
		sbytes := make([]byte, n)
		for i := 0; i < n; i++ {
			if rng.Intn(2) == 0 {
				sbytes[i] = '0'
			} else {
				sbytes[i] = '1'
			}
		}
		s := string(sbytes)
		pos := make([][]int, n)
		sets := make([][]int, k+1)
		for i := 0; i < n; i++ {
			tcnt := rng.Intn(3)
			used := map[int]bool{}
			for j := 0; j < tcnt; j++ {
				idx := rng.Intn(k) + 1
				if used[idx] {
					continue
				}
				used[idx] = true
				pos[i] = append(pos[i], idx)
				sets[idx] = append(sets[idx], i+1)
			}
		}
		fmt.Fprintf(&sb, "%d %d\n%s\n", n, k, s)
		for i := 1; i <= k; i++ {
			if len(sets[i]) == 0 {
				sets[i] = append(sets[i], rng.Intn(n)+1)
			}
			fmt.Fprintf(&sb, "%d\n", len(sets[i]))
			for j, v := range sets[i] {
				if j > 0 {
					sb.WriteByte(' ')
				}
				fmt.Fprintf(&sb, "%d", v)
			}
			sb.WriteByte('\n')
		}
		input := sb.String()
		expected := solveC(n, k, s, pos)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", t, err, input)
			os.Exit(1)
		}
		lines := strings.Split(strings.TrimSpace(out), "\n")
		if len(lines) != n {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d lines got %d\ninput:\n%s", t, n, len(lines), input)
			os.Exit(1)
		}
		for i := 0; i < n; i++ {
			got := strings.TrimSpace(lines[i])
			want := fmt.Sprintf("%d", expected[i])
			if got != want {
				fmt.Fprintf(os.Stderr, "case %d line %d: expected %s got %s\ninput:\n%s", t, i+1, want, got, input)
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed")
}

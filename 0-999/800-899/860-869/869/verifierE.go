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

type Fenwick2D struct {
	n, m int
	tree [][]uint64
}

func NewFenwick2D(n, m int) *Fenwick2D {
	tree := make([][]uint64, n+2)
	for i := range tree {
		tree[i] = make([]uint64, m+2)
	}
	return &Fenwick2D{n: n + 1, m: m + 1, tree: tree}
}

func (f *Fenwick2D) add(x, y int, val uint64) {
	for i := x; i <= f.n; i += i & -i {
		row := f.tree[i]
		for j := y; j <= f.m; j += j & -j {
			row[j] ^= val
		}
	}
}

func (f *Fenwick2D) RangeAdd(x1, y1, x2, y2 int, val uint64) {
	f.add(x1, y1, val)
	f.add(x1, y2+1, val)
	f.add(x2+1, y1, val)
	f.add(x2+1, y2+1, val)
}

func (f *Fenwick2D) Query(x, y int) uint64 {
	var res uint64
	for i := x; i > 0; i -= i & -i {
		row := f.tree[i]
		for j := y; j > 0; j -= j & -j {
			res ^= row[j]
		}
	}
	return res
}

type Rect struct{ r1, c1, r2, c2 int }

func solveE(n, m int, ops []struct {
	t, r1, c1, r2, c2 int
}) []string {
	fw := NewFenwick2D(n+2, m+2)
	rectID := make(map[Rect]uint64)
	rand.Seed(time.Now().UnixNano())
	var res []string
	for _, op := range ops {
		switch op.t {
		case 1:
			id := rand.Uint64()
			rectID[Rect{op.r1, op.c1, op.r2, op.c2}] = id
			fw.RangeAdd(op.r1, op.c1, op.r2, op.c2, id)
		case 2:
			rect := Rect{op.r1, op.c1, op.r2, op.c2}
			if id, ok := rectID[rect]; ok {
				fw.RangeAdd(op.r1, op.c1, op.r2, op.c2, id)
				delete(rectID, rect)
			}
		case 3:
			v1 := fw.Query(op.r1, op.c1)
			v2 := fw.Query(op.r2, op.c2)
			if v1 == v2 {
				res = append(res, "Yes")
			} else {
				res = append(res, "No")
			}
		}
	}
	return res
}

func run(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n := rng.Intn(10) + 4
		m := rng.Intn(10) + 4
		q := rng.Intn(20) + 1
		var ops []struct {
			t, r1, c1, r2, c2 int
		}
		active := make(map[Rect]struct{})
		for j := 0; j < q; j++ {
			t := rng.Intn(3) + 1
			if t == 1 || (t == 2 && len(active) == 0) {
				t = 1
			}
			if t == 1 {
				r1 := rng.Intn(n-2) + 2
				r2 := rng.Intn(n-r1) + r1
				c1 := rng.Intn(m-2) + 2
				c2 := rng.Intn(m-c1) + c1
				ops = append(ops, struct{ t, r1, c1, r2, c2 int }{1, r1, c1, r2, c2})
				active[Rect{r1, c1, r2, c2}] = struct{}{}
			} else if t == 2 {
				idx := rng.Intn(len(active))
				k := 0
				var rect Rect
				for r := range active {
					if k == idx {
						rect = r
						break
					}
					k++
				}
				ops = append(ops, struct{ t, r1, c1, r2, c2 int }{2, rect.r1, rect.c1, rect.r2, rect.c2})
				delete(active, rect)
			} else {
				r1 := rng.Intn(n) + 1
				c1 := rng.Intn(m) + 1
				r2 := rng.Intn(n) + 1
				c2 := rng.Intn(m) + 1
				ops = append(ops, struct{ t, r1, c1, r2, c2 int }{3, r1, c1, r2, c2})
			}
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d %d\n", n, m, q))
		for _, op := range ops {
			sb.WriteString(fmt.Sprintf("%d %d %d %d %d\n", op.t, op.r1, op.c1, op.r2, op.c2))
		}
		input := sb.String()
		answers := solveE(n, m, ops)
		expected := strings.Join(answers, "\n")
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expected) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %q got %q\ninput:%s", i+1, expected, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

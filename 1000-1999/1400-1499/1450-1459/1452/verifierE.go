package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type testE struct {
	n, m, k int
	l, r    []int
}

func solveE(t testE) int {
	n, m, k := t.n, t.m, t.k
	l, r := t.l, t.r
	w := n - k + 1
	f := make([][]int, m)
	for p := 0; p < m; p++ {
		f[p] = make([]int, w+1)
		for i := 1; i <= w; i++ {
			li, ri := l[p], r[p]
			wi, wj := i, i+k-1
			a := li
			if wi > a {
				a = wi
			}
			b := ri
			if wj < b {
				b = wj
			}
			if a <= b {
				f[p][i] = b - a + 1
			} else {
				f[p][i] = 0
			}
		}
	}
	base := make([]int, w+2)
	for i := 1; i <= w; i++ {
		s := 0
		for p := 0; p < m; p++ {
			s += f[p][i]
		}
		base[i] = s
	}
	bestLeft := make([]int, w+2)
	for i := 1; i <= w; i++ {
		if i-k >= 1 {
			if base[i-k] > bestLeft[i-1] {
				bestLeft[i] = base[i-k]
			} else {
				bestLeft[i] = bestLeft[i-1]
			}
		} else {
			bestLeft[i] = bestLeft[i-1]
		}
	}
	bestRight := make([]int, w+2)
	for i := w; i >= 1; i-- {
		if i+k <= w {
			if base[i+k] > bestRight[i+1] {
				bestRight[i] = base[i+k]
			} else {
				bestRight[i] = bestRight[i+1]
			}
		} else {
			bestRight[i] = bestRight[i+1]
		}
	}
	best := 0
	for i := 1; i <= w; i++ {
		if bestLeft[i] > 0 {
			tmp := base[i] + bestLeft[i]
			if tmp > best {
				best = tmp
			}
		}
		if bestRight[i] > 0 {
			tmp := base[i] + bestRight[i]
			if tmp > best {
				best = tmp
			}
		}
		if base[i] > best {
			best = base[i]
		}
		lo := i - k + 1
		if lo < 1 {
			lo = 1
		}
		hi := i + k - 1
		if hi > w {
			hi = w
		}
		for j := lo; j <= hi; j++ {
			mn := 0
			for p := 0; p < m; p++ {
				x := f[p][i]
				y := f[p][j]
				if x < y {
					mn += x
				} else {
					mn += y
				}
			}
			tmp := base[i] + base[j] - mn
			if tmp > best {
				best = tmp
			}
		}
	}
	return best
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return out.String(), fmt.Errorf("run error: %v\n%s", err, out.String())
	}
	return out.String(), nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		return
	}
	binary := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	const tnum = 100
	for ti := 0; ti < tnum; ti++ {
		n := rand.Intn(8) + 2
		m := rand.Intn(3) + 1
		k := rand.Intn(n) + 1
		l := make([]int, m)
		r := make([]int, m)
		for i := 0; i < m; i++ {
			a := rand.Intn(n) + 1
			b := rand.Intn(n) + 1
			if a > b {
				a, b = b, a
			}
			l[i] = a
			r[i] = b
		}
		test := testE{n: n, m: m, k: k, l: l, r: r}
		exp := solveE(test)
		var in strings.Builder
		fmt.Fprintf(&in, "%d %d %d\n", n, m, k)
		for i := 0; i < m; i++ {
			fmt.Fprintf(&in, "%d %d\n", l[i], r[i])
		}
		out, err := runBinary(binary, in.String())
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		vStr := strings.TrimSpace(out)
		v, err := strconv.Atoi(vStr)
		if err != nil || v != exp {
			fmt.Printf("test %d failed: expected=%d got=%s\n", ti+1, exp, vStr)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}

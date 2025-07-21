package main

import (
	"bytes"
	"fmt"
	"math"
	"os"
	"os/exec"
	"sort"
	"strings"
)

const M = 200010

var vis [M]bool

type pair struct{ first, second int }

var p [M]pair

func init() {
	for i := 0; i < 310; i++ {
		for j := 0; j < 310; j++ {
			s := i*i + j*j
			if s < M {
				vis[s] = true
				p[s] = pair{i, j}
			}
		}
	}
}

func getSign(arr []int) {
	type pi struct{ val, idx int }
	n := len(arr)
	arrp := make([]pi, n)
	for i := 0; i < n; i++ {
		arrp[i] = pi{arr[i], i}
	}
	sort.Slice(arrp, func(i, j int) bool { return arrp[i].val > arrp[j].val })
	d := 0
	for i := 0; i+1 < n; i += 2 {
		if arrp[i].val == arrp[i+1].val {
			idx := arrp[i].idx
			arr[idx] = -arr[idx]
		} else {
			diff := arrp[i].val - arrp[i+1].val
			if d <= 0 {
				idx := arrp[i+1].idx
				arr[idx] = -arr[idx]
				d += diff
			} else {
				idx := arrp[i].idx
				arr[idx] = -arr[idx]
				d -= diff
			}
		}
	}
	if n%2 == 1 {
		idx := arrp[n-1].idx
		arr[idx] = -arr[idx]
	}
}

func polygon(n int) [][2]int {
	x := make([]int, n)
	y := make([]int, n)
	cnt := 0
	sum := 0
	i := 1
	for cnt < n-1 {
		if i < M && vis[i] {
			x[cnt] = p[i].first
			y[cnt] = p[i].second
			sum += i
			cnt++
		}
		i++
	}
	var a, b, tcount int
	for j := i; tcount < 2; j++ {
		if j < M && vis[j] {
			if tcount == 0 {
				a = j
			} else {
				b = j
			}
			tcount++
		}
	}
	ok := false
	if (sum+a)%2 == 0 {
		x[cnt] = p[a].first
		y[cnt] = p[a].second
	} else if (sum+b)%2 == 0 {
		x[cnt] = p[b].first
		y[cnt] = p[b].second
	} else {
		x[cnt] = p[a].first
		y[cnt] = p[a].second
		if sum%2 != 0 {
			x[0] = p[b].first
			y[0] = p[b].second
		} else {
			if n >= 4 {
				x[3] = p[b].first
				y[3] = p[b].second
			} else {
				x[0] = p[b].first
				y[0] = p[b].second
			}
		}
		ok = true
	}
	cnt++
	sum = 0
	for k := 0; k < n; k++ {
		sum += x[k]
	}
	if sum%2 != 0 {
		if !ok {
			x[0], y[0] = y[0], x[0]
		} else {
			if n >= 4 {
				x[3], y[3] = y[3], x[3]
			} else {
				x[0], y[0] = y[0], x[0]
			}
		}
	}
	getSign(x)
	getSign(y)
	ang := make([]float64, n)
	for k := 0; k < n; k++ {
		ang[k] = math.Atan2(float64(y[k]), float64(x[k]))
	}
	ids := make([]int, n)
	for k := 0; k < n; k++ {
		ids[k] = k
	}
	sort.Slice(ids, func(i, j int) bool { return ang[ids[i]] < ang[ids[j]] })
	coords := make([][2]int, n)
	tx, ty := 0, 0
	for idx, id := range ids {
		tx += x[id]
		ty += y[id]
		coords[idx] = [2]int{tx, ty}
	}
	return coords
}

func runCase(exe string, n int) error {
	input := fmt.Sprintf("%d\n", n)
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	tokens := strings.Fields(out.String())
	if len(tokens) != 1+2*n {
		return fmt.Errorf("expected %d tokens got %d", 1+2*n, len(tokens))
	}
	if tokens[0] != "YES" {
		return fmt.Errorf("expected YES got %s", tokens[0])
	}
	exp := polygon(n)
	idx := 1
	for i := 0; i < n; i++ {
		var x, y int
		if _, err := fmt.Sscan(tokens[idx], &x); err != nil {
			return fmt.Errorf("bad output token %q", tokens[idx])
		}
		idx++
		if _, err := fmt.Sscan(tokens[idx], &y); err != nil {
			return fmt.Errorf("bad output token %q", tokens[idx])
		}
		idx++
		if x != exp[i][0] || y != exp[i][1] {
			return fmt.Errorf("vertex %d mismatch", i+1)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	for n := 3; n <= 102; n++ {
		if err := runCase(exe, n); err != nil {
			fmt.Fprintf(os.Stderr, "case n=%d failed: %v\n", n, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

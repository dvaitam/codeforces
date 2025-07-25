package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type rect struct {
	x1, y1, x2, y2 int64
}

type pair struct {
	val int64
	id  int
}

func solveC(recs []rect, k int) int64 {
	n := len(recs)
	centersX := make([]pair, n)
	centersY := make([]pair, n)
	for i := 0; i < n; i++ {
		cx := recs[i].x1 + recs[i].x2
		cy := recs[i].y1 + recs[i].y2
		centersX[i] = pair{val: cx, id: i}
		centersY[i] = pair{val: cy, id: i}
	}
	sort.Slice(centersX, func(i, j int) bool { return centersX[i].val < centersX[j].val })
	sort.Slice(centersY, func(i, j int) bool { return centersY[i].val < centersY[j].val })
	best := int64(1<<63 - 1)
	for a := 0; a <= k; a++ {
		for b := 0; b <= k; b++ {
			for c := 0; c <= k; c++ {
				for d := 0; d <= k; d++ {
					removed := make(map[int]struct{})
					for i := 0; i < a && i < n; i++ {
						removed[centersX[i].id] = struct{}{}
					}
					for i := 0; i < b && i < n; i++ {
						removed[centersX[n-1-i].id] = struct{}{}
					}
					for i := 0; i < c && i < n; i++ {
						removed[centersY[i].id] = struct{}{}
					}
					for i := 0; i < d && i < n; i++ {
						removed[centersY[n-1-i].id] = struct{}{}
					}
					if len(removed) > k {
						continue
					}
					li := 0
					for li < n {
						if _, ok := removed[centersX[li].id]; !ok {
							break
						}
						li++
					}
					ri := n - 1
					for ri >= 0 {
						if _, ok := removed[centersX[ri].id]; !ok {
							break
						}
						ri--
					}
					lj := 0
					for lj < n {
						if _, ok := removed[centersY[lj].id]; !ok {
							break
						}
						lj++
					}
					rj := n - 1
					for rj >= 0 {
						if _, ok := removed[centersY[rj].id]; !ok {
							break
						}
						rj--
					}
					if li > ri || lj > rj {
						continue
					}
					dx := centersX[ri].val - centersX[li].val
					dy := centersY[rj].val - centersY[lj].val
					w := dx / 2
					if dx%2 != 0 {
						w++
					}
					h := dy / 2
					if dy%2 != 0 {
						h++
					}
					if w < 1 {
						w = 1
					}
					if h < 1 {
						h = 1
					}
					area := w * h
					if area < best {
						best = area
					}
				}
			}
		}
	}
	if best == int64(1<<63-1) {
		best = 1
	}
	return best
}

func run(binary, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, binary)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(3)
	for t := 0; t < 100; t++ {
		n := rand.Intn(5) + 1
		k := rand.Intn(n + 1)
		recs := make([]rect, n)
		for i := 0; i < n; i++ {
			x1 := int64(rand.Intn(10))
			y1 := int64(rand.Intn(10))
			x2 := x1 + int64(rand.Intn(5)+1)
			y2 := y1 + int64(rand.Intn(5)+1)
			recs[i] = rect{x1, y1, x2, y2}
		}
		input := fmt.Sprintf("%d %d\n", n, k)
		for i, r := range recs {
			input += fmt.Sprintf("%d %d %d %d", r.x1, r.y1, r.x2, r.y2)
			if i+1 < n {
				input += "\n"
			} else {
				input += "\n"
			}
		}
		expected := fmt.Sprintf("%d", solveC(recs, k))
		got, err := run(bin, input)
		if err != nil {
			fmt.Println("test", t, "runtime error:", err)
			fmt.Println("output:", got)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expected) {
			fmt.Println("test", t, "failed")
			fmt.Println("input:\n" + input)
			fmt.Println("expected:", expected)
			fmt.Println("got:", got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

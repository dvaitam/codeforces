package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"time"
)

type Cake struct {
	idx int
	r2h float64
}

type BIT struct {
	n    int
	tree []float64
}

func NewBIT(n int) *BIT {
	return &BIT{n: n, tree: make([]float64, n+1)}
}

func (b *BIT) Update(i int, val float64) {
	for ; i <= b.n; i += i & -i {
		if b.tree[i] < val {
			b.tree[i] = val
		}
	}
}

func (b *BIT) Query(i int) float64 {
	var ret float64
	for ; i > 0; i -= i & -i {
		if ret < b.tree[i] {
			ret = b.tree[i]
		}
	}
	return ret
}

func solveD(n int, pairs [][2]int) float64 {
	cakes := make([]Cake, n)
	for i := 0; i < n; i++ {
		r := pairs[i][0]
		h := pairs[i][1]
		cakes[i].idx = n - i
		cakes[i].r2h = float64(r) * float64(r) * float64(h)
	}
	sort.Slice(cakes, func(i, j int) bool {
		return cakes[i].r2h > cakes[j].r2h
	})
	bit := NewBIT(n)
	v := make([]float64, n)
	var ans float64
	last := 0
	for i := 0; i < n; i++ {
		v[i] = cakes[i].r2h + bit.Query(cakes[i].idx-1)
		if i < n-1 && cakes[i].r2h == cakes[i+1].r2h {
		} else {
			for j := last; j <= i; j++ {
				bit.Update(cakes[j].idx, v[j])
			}
			last = i + 1
		}
		if v[i] > ans {
			ans = v[i]
		}
	}
	return ans * math.Pi
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	for t := 0; t < 100; t++ {
		n := rand.Intn(30) + 1
		pairs := make([][2]int, n)
		var input bytes.Buffer
		fmt.Fprintln(&input, n)
		for i := 0; i < n; i++ {
			r := rand.Intn(100) + 1
			h := rand.Intn(100) + 1
			pairs[i] = [2]int{r, h}
			fmt.Fprintf(&input, "%d %d\n", r, h)
		}
		expected := solveD(n, pairs)
		cmd := exec.Command(binary)
		cmd.Stdin = &input
		var out bytes.Buffer
		cmd.Stdout = &out
		if err := cmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "test %d: binary error: %v\n", t, err)
			os.Exit(1)
		}
		scanner := bufio.NewScanner(&out)
		if !scanner.Scan() {
			fmt.Fprintf(os.Stderr, "test %d: no output\n", t)
			os.Exit(1)
		}
		var got float64
		if _, err := fmt.Sscan(scanner.Text(), &got); err != nil {
			fmt.Fprintf(os.Stderr, "test %d: invalid output\n", t)
			os.Exit(1)
		}
		if math.Abs(got-expected) > 1e-6*math.Max(1, math.Abs(expected)) {
			fmt.Fprintf(os.Stderr, "test %d: expected %.6f got %.6f\n", t, expected, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed.")
}

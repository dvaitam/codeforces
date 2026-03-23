package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// Embedded correct solver for 150C (CF-accepted)

type Node struct {
	min   int64
	max   int64
	best  int64
	valid bool
}

func combine(a, b Node) Node {
	if !a.valid {
		return b
	}
	if !b.valid {
		return a
	}
	best := a.best
	if b.best > best {
		best = b.best
	}
	cross := b.max - a.min
	if cross > best {
		best = cross
	}
	minv := a.min
	if b.min < minv {
		minv = b.min
	}
	maxv := a.max
	if b.max > maxv {
		maxv = b.max
	}
	return Node{min: minv, max: maxv, best: best, valid: true}
}

func solveC(input string) string {
	data := []byte(input)
	ptr := 0
	nextInt := func() int {
		n := len(data)
		for ptr < n && (data[ptr] < '0' || data[ptr] > '9') {
			ptr++
		}
		val := 0
		for ptr < n && data[ptr] >= '0' && data[ptr] <= '9' {
			val = val*10 + int(data[ptr]-'0')
			ptr++
		}
		return val
	}

	n := nextInt()
	m := nextInt()
	c := int64(nextInt())

	x := make([]int64, n)
	for i := 0; i < n; i++ {
		x[i] = int64(nextInt())
	}

	pref := make([]int64, n)
	for i := 1; i < n; i++ {
		pref[i] = pref[i-1] + int64(nextInt())
	}

	t := make([]int64, n)
	for i := 0; i < n; i++ {
		t[i] = 100*x[i] - 2*c*pref[i]
	}

	size := 1
	for size < n {
		size <<= 1
	}
	tree := make([]Node, size<<1)
	for i := 0; i < n; i++ {
		tree[size+i] = Node{min: t[i], max: t[i], best: 0, valid: true}
	}
	for i := size - 1; i >= 1; i-- {
		tree[i] = combine(tree[i<<1], tree[i<<1|1])
	}

	query := func(l, r int) Node {
		l += size
		r += size
		var left, right Node
		for l < r {
			if l&1 == 1 {
				left = combine(left, tree[l])
				l++
			}
			if r&1 == 1 {
				r--
				right = combine(tree[r], right)
			}
			l >>= 1
			r >>= 1
		}
		return combine(left, right)
	}

	var sum int64
	for i := 0; i < m; i++ {
		a := nextInt()
		b := nextInt()
		res := query(a-1, b)
		sum += res.best
	}

	return fmt.Sprintf("%d.%03d", sum/200, (sum%200)*5)
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

func parseFloat(s string) float64 {
	v, _ := strconv.ParseFloat(s, 64)
	return v
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n := rng.Intn(8) + 2
		m := rng.Intn(5) + 1
		c := rng.Intn(10) + 1
		x := make([]int64, n)
		for j := 1; j < n; j++ {
			x[j] = x[j-1] + int64(rng.Intn(10)+1)
		}
		p := make([]int64, n)
		for j := 1; j < n; j++ {
			p[j] = int64(rng.Intn(101))
		}
		queries := make([][2]int, m)
		for j := 0; j < m; j++ {
			l := rng.Intn(n-1) + 1
			r := rng.Intn(n-l) + l + 1
			queries[j] = [2]int{l, r}
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d %d\n", n, m, c))
		for j := 0; j < n; j++ {
			sb.WriteString(fmt.Sprintf("%d ", x[j]))
		}
		sb.WriteString("\n")
		for j := 1; j < n; j++ {
			sb.WriteString(fmt.Sprintf("%d ", p[j]))
		}
		sb.WriteString("\n")
		for j := 0; j < m; j++ {
			sb.WriteString(fmt.Sprintf("%d %d\n", queries[j][0], queries[j][1]))
		}
		input := sb.String()
		expected := solveC(input)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		// Compare as floats with tolerance
		expVal := parseFloat(expected)
		gotVal := parseFloat(strings.TrimSpace(got))
		diff := math.Abs(expVal - gotVal)
		denom := math.Max(1.0, math.Abs(expVal))
		if diff/denom > 1e-6 {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expected, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

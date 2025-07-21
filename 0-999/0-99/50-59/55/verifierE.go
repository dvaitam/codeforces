package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type Vec struct {
	x, y int64
	half int
}

type Vecs []Vec

func (v Vecs) Len() int      { return len(v) }
func (v Vecs) Swap(i, j int) { v[i], v[j] = v[j], v[i] }
func (v Vecs) Less(i, j int) bool {
	if v[i].half != v[j].half {
		return v[i].half < v[j].half
	}
	return v[i].x*v[j].y-v[i].y*v[j].x > 0
}

func countTriangles(n int, px, py int64, xs, ys []int64) int64 {
	vecs := make(Vecs, n)
	for i := 0; i < n; i++ {
		dx := xs[i] - px
		dy := ys[i] - py
		half := 1
		if dy > 0 || (dy == 0 && dx > 0) {
			half = 0
		}
		vecs[i] = Vec{x: dx, y: dy, half: half}
	}
	sort.Sort(vecs)
	vecs = append(vecs, vecs...)
	var bad int64
	j := 0
	for i := 0; i < n; i++ {
		if j < i+1 {
			j = i + 1
		}
		for j < i+n && (vecs[i].x*vecs[j].y-vecs[i].y*vecs[j].x) > 0 {
			j++
		}
		k := int64(j - i - 1)
		if k >= 2 {
			bad += k * (k - 1) / 2
		}
	}
	total := int64(n)
	total = total * (total - 1) * (total - 2) / 6
	return total - bad
}

func generatePolygon(n int) ([]int64, []int64) {
	r := rand.Int63n(1000) + 10
	xs := make([]int64, n)
	ys := make([]int64, n)
	for i := 0; i < n; i++ {
		angle := -2 * math.Pi * float64(i) / float64(n)
		xs[i] = int64(math.Round(float64(r) * math.Cos(angle)))
		ys[i] = int64(math.Round(float64(r) * math.Sin(angle)))
	}
	return xs, ys
}

func expected(n int, xs, ys []int64, queries [][2]int64) string {
	var sb strings.Builder
	for i, q := range queries {
		ans := countTriangles(n, q[0], q[1], xs, ys)
		if i+1 == len(queries) {
			fmt.Fprintf(&sb, "%d", ans)
		} else {
			fmt.Fprintf(&sb, "%d\n", ans)
		}
	}
	return sb.String()
}

func generateCase() (string, string) {
	n := rand.Intn(4) + 3 // 3..6
	xs, ys := generatePolygon(n)
	t := rand.Intn(5) + 1
	queries := make([][2]int64, t)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for i := 0; i < n; i++ {
		fmt.Fprintf(&sb, "%d %d\n", xs[i], ys[i])
	}
	fmt.Fprintf(&sb, "%d\n", t)
	for i := 0; i < t; i++ {
		px := rand.Int63n(2000) - 1000
		py := rand.Int63n(2000) - 1000
		queries[i] = [2]int64{px, py}
		fmt.Fprintf(&sb, "%d %d\n", px, py)
	}
	return sb.String(), expected(n, xs, ys, queries)
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
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 100; i++ {
		input, exp := generateCase()
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

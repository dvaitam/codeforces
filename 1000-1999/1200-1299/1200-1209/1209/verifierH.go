package main

import (
	"bytes"
	"container/heap"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// ---- Embedded oracle solver (from correct CF-accepted solution) ----

type oracleItem struct {
	v   float64
	idx int
}

type oracleMaxHeap []oracleItem

func (h oracleMaxHeap) Len() int            { return len(h) }
func (h oracleMaxHeap) Less(i, j int) bool  { return h[i].v > h[j].v }
func (h oracleMaxHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *oracleMaxHeap) Push(x interface{}) { *h = append(*h, x.(oracleItem)) }
func (h *oracleMaxHeap) Pop() interface{} {
	old := *h
	n := len(old)
	it := old[n-1]
	*h = old[:n-1]
	return it
}
func (h *oracleMaxHeap) Peek() oracleItem { return (*h)[0] }

func oracleSolve(input string) string {
	words := strings.Fields(input)
	wi := 0
	nextInt := func() int {
		v, _ := strconv.Atoi(words[wi])
		wi++
		return v
	}
	nextFloat := func() float64 {
		v, _ := strconv.ParseFloat(words[wi], 64)
		wi++
		return v
	}

	n := nextInt()
	L := nextInt()

	x0 := make([]int, n)
	y0 := make([]int, n)
	v0 := make([]float64, n)
	for i := 0; i < n; i++ {
		x0[i] = nextInt()
		y0[i] = nextInt()
		v0[i] = nextFloat()
	}

	var l, v []float64
	if x0[0] > 0 {
		l = append(l, float64(x0[0]))
		v = append(v, 0)
	}
	for i := 0; i < n; i++ {
		l = append(l, float64(y0[i]-x0[i]))
		v = append(v, v0[i])
		if i+1 < n && y0[i] < x0[i+1] {
			l = append(l, float64(x0[i+1]-y0[i]))
			v = append(v, 0)
		}
		if i+1 == n && y0[i] < L {
			l = append(l, float64(L-y0[i]))
			v = append(v, 0)
		}
	}
	m := len(v)
	x := make([]float64, m)
	for i := range x {
		x[i] = 2.0
	}
	t := make([]float64, m)
	h := &oracleMaxHeap{}
	heap.Init(h)
	var bal, res float64
	const eps = 1e-10
	for i := 0; i < m; i++ {
		heap.Push(h, oracleItem{v: v[i], idx: i})
		t[i] = l[i] / (v[i] + x[i])
		bal += t[i] * (1 - x[i])
		res += t[i]
		for bal < -eps {
			if h.Len() == 0 {
				break
			}
			top := h.Peek()
			j := top.idx
			bal -= t[j] * (1 - x[j])
			res -= t[j]
			var newx float64
			if v[j] > eps && bal+l[j]/v[j] < eps {
				newx = 0
			} else {
				newx = (l[j] + bal*v[j]) / (l[j] - bal)
				if newx < 0 {
					newx = 0
				}
				if newx > 2 {
					newx = 2
				}
			}
			x[j] = newx
			if x[j] < eps {
				heap.Pop(h)
			}
			t[j] = l[j] / (v[j] + x[j])
			bal += t[j] * (1 - x[j])
			res += t[j]
		}
	}
	return fmt.Sprintf("%.15f", res)
}

// ---- Test generation ----

func genCase(rng *rand.Rand) string {
	n := rng.Intn(3) + 1
	L := rng.Intn(20) + 5
	type seg struct{ x, y int }
	segs := make([]seg, n)
	pos := 1
	for i := 0; i < n; i++ {
		x := pos + rng.Intn(2)
		if x >= L {
			x = L - 1
		}
		y := x + rng.Intn(3) + 1
		if y > L {
			y = L
		}
		if y <= x {
			y = x + 1
		}
		segs[i] = seg{x, y}
		pos = y + rng.Intn(2)
		if pos >= L {
			pos = L - 1
		}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, L))
	for i := 0; i < n; i++ {
		s := float64(rng.Intn(2)+1) + rng.Float64()
		sb.WriteString(fmt.Sprintf("%d %d %.2f\n", segs[i].x, segs[i].y, s))
	}
	return sb.String()
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func prepareBinary(path, tag string) (string, error) {
	if strings.HasSuffix(path, ".go") {
		bin := fmt.Sprintf("%s/%s_%d", os.TempDir(), tag, time.Now().UnixNano())
		cmd := exec.Command("go", "build", "-o", bin, path)
		if out, err := cmd.CombinedOutput(); err != nil {
			return "", fmt.Errorf("build %s: %v\n%s", path, err, out)
		}
		return bin, nil
	}
	return path, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierH.go /path/to/binary")
		os.Exit(1)
	}
	candPath, err := prepareBinary(os.Args[1], "candH")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input := genCase(rng)
		exp := oracleSolve(input)
		gotRaw, err := runBinary(candPath, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on case %d: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		got := strings.TrimSpace(gotRaw)

		// Compare as floats with tolerance
		expF, err1 := strconv.ParseFloat(strings.TrimSpace(exp), 64)
		gotF, err2 := strconv.ParseFloat(got, 64)
		if err1 != nil || err2 != nil {
			if strings.TrimSpace(exp) != got {
				fmt.Printf("case %d failed\ninput:\n%sexpected:\n%s\ngot:\n%s", i+1, input, exp, got)
				os.Exit(1)
			}
		} else if math.Abs(expF-gotF) > 1e-6*math.Max(1, math.Abs(expF)) {
			fmt.Printf("case %d failed\ninput:\n%sexpected:\n%s\ngot:\n%s", i+1, input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

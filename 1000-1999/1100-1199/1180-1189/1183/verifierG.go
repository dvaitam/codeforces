package main

import (
	"bytes"
	"container/heap"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"
)

func genCase(rng *rand.Rand) string {
	q := rng.Intn(5) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", q))
	for i := 0; i < q; i++ {
		n := rng.Intn(10) + 1
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for j := 0; j < n; j++ {
			a := rng.Intn(n) + 1
			f := rng.Intn(2)
			sb.WriteString(fmt.Sprintf("%d %d\n", a, f))
		}
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input := genCase(rng)
		out, err := runProg(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d runtime error: %v\n%s", i+1, err, out)
			os.Exit(1)
		}
		// Compute expected results per test programmatically and compare token-wise
		if err := validate1183G(input, out); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s\noutput:\n%s\n", i+1, err, input, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

func runProg(prog, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(prog, ".go") {
		cmd = exec.Command("go", "run", prog)
	} else {
		cmd = exec.Command(prog)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return out.String() + stderr.String(), fmt.Errorf("%v", err)
	}
	return strings.TrimSpace(out.String()), nil
}

type pair struct{ s, c1 int }

// computeExpected returns per query (total, goodMax)
func computeExpected(input string) ([][2]int64, error) {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	if len(lines) == 0 {
		return nil, fmt.Errorf("empty input")
	}
	q, err := strconv.Atoi(strings.TrimSpace(lines[0]))
	if err != nil {
		return nil, fmt.Errorf("bad q: %v", err)
	}
	res := make([][2]int64, 0, q)
	idx := 1
	for t := 0; t < q; t++ {
		if idx >= len(lines) {
			return nil, fmt.Errorf("unexpected EOF parsing n at case %d", t+1)
		}
		n, err := strconv.Atoi(strings.TrimSpace(lines[idx]))
		if err != nil {
			return nil, fmt.Errorf("bad n at case %d: %v", t+1, err)
		}
		idx++
		cnt := make(map[int]pair)
		for j := 0; j < n; j++ {
			if idx >= len(lines) {
				return nil, fmt.Errorf("unexpected EOF reading pairs at case %d", t+1)
			}
			f := strings.Fields(lines[idx])
			idx++
			if len(f) < 2 {
				return nil, fmt.Errorf("bad pair line at case %d", t+1)
			}
			a, _ := strconv.Atoi(f[0])
			ff, _ := strconv.Atoi(f[1])
			v := cnt[a]
			v.s++
			if ff == 1 {
				v.c1++
			}
			cnt[a] = v
		}
		types := make([]pair, 0, len(cnt))
		for _, v := range cnt {
			types = append(types, v)
		}
		sort.Slice(types, func(i, j int) bool { return types[i].s > types[j].s })
		// build ks with strictly decreasing upper bound
		ks := make([]int, 0, len(types))
		last := int(1e9)
		for _, v := range types {
			x := v.s
			if x > last-1 {
				x = last - 1
			}
			if x <= 0 {
				break
			}
			ks = append(ks, x)
			last = x
		}
		var total int64
		for _, k := range ks {
			total += int64(k)
		}
		// compute max good
		sort.Slice(types, func(i, j int) bool { return types[i].s > types[j].s })
		sort.Slice(ks, func(i, j int) bool { return ks[i] > ks[j] })
		// max-heap for c1
		h := &intMaxHeap{}
		heap.Init(h)
		p := 0
		var good int64
		for _, k := range ks {
			for p < len(types) && types[p].s >= k {
				heap.Push(h, types[p].c1)
				p++
			}
			if h.Len() == 0 {
				// should not happen
				continue
			}
			c1 := heap.Pop(h).(int)
			if c1 < k {
				good += int64(c1)
			} else {
				good += int64(k)
			}
		}
		res = append(res, [2]int64{total, good})
	}
	return res, nil
}

type intMaxHeap []int

func (h intMaxHeap) Len() int            { return len(h) }
func (h intMaxHeap) Less(i, j int) bool  { return h[i] > h[j] }
func (h intMaxHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *intMaxHeap) Push(x interface{}) { *h = append(*h, x.(int)) }
func (h *intMaxHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

func validate1183G(input, output string) error {
	exp, err := computeExpected(input)
	if err != nil {
		return err
	}
	toks := strings.Fields(strings.TrimSpace(output))
	if len(toks) != 2*len(exp) {
		return fmt.Errorf("wrong number of tokens: expected %d got %d", 2*len(exp), len(toks))
	}
	for i := 0; i < len(exp); i++ {
		a, err1 := strconv.ParseInt(toks[2*i], 10, 64)
		b, err2 := strconv.ParseInt(toks[2*i+1], 10, 64)
		if err1 != nil || err2 != nil {
			return fmt.Errorf("non-integer output at case %d", i+1)
		}
		if a != exp[i][0] || b != exp[i][1] {
			return fmt.Errorf("expected: %d %d\n got: %d %d\n", exp[i][0], exp[i][1], a, b)
		}
	}
	return nil
}

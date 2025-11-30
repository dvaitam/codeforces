package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"
)

type Event struct {
	x   int
	typ int
	y1  int
	y2  int
	val int64
	idx int
}

type BIT struct {
	n    int
	bit0 []int64
	bit1 []int64
}

func newBIT(n int) *BIT {
	return &BIT{n: n, bit0: make([]int64, n+5), bit1: make([]int64, n+5)}
}

func (b *BIT) upd(bit []int64, i int, v int64) {
	for ; i <= b.n; i += i & -i {
		bit[i] += v
	}
}

func (b *BIT) updateRange(l, r int, v int64) {
	if l > r {
		return
	}
	b.upd(b.bit0, l, v)
	b.upd(b.bit0, r+1, -v)
	b.upd(b.bit1, l, v*int64(l-1))
	b.upd(b.bit1, r+1, -v*int64(r))
}

func (b *BIT) prefixSum(i int) int64 {
	var s0, s1 int64
	for j := i; j > 0; j -= j & -j {
		s0 += b.bit0[j]
		s1 += b.bit1[j]
	}
	return s0*int64(i) - s1
}

func (b *BIT) queryRange(l, r int) int64 {
	if l > r {
		return 0
	}
	return b.prefixSum(r) - b.prefixSum(l-1)
}

// solve390E mirrors 390E.go using in-memory I/O.
func solve390E(input string) (string, error) {
	in := bufio.NewReader(strings.NewReader(input))
	var n, m, w int
	if _, err := fmt.Fscan(in, &n, &m, &w); err != nil {
		return "", err
	}
	events := make([]Event, 0, w*10)
	qcnt := 0
	for i := 0; i < w; i++ {
		var t int
		if _, err := fmt.Fscan(in, &t); err != nil {
			return "", err
		}
		if t == 0 {
			var x1, y1, x2, y2 int
			var v int64
			fmt.Fscan(in, &x1, &y1, &x2, &y2, &v)
			events = append(events, Event{x: x1, typ: 0, y1: y1, y2: y2, val: v})
			events = append(events, Event{x: x2 + 1, typ: 0, y1: y1, y2: y2, val: -v})
		} else {
			var x1, y1, x2, y2 int
			fmt.Fscan(in, &x1, &y1, &x2, &y2)
			id := qcnt
			qcnt++
			events = append(events, Event{x: x2, typ: 1, y1: y1, y2: y2, val: 1, idx: id})
			events = append(events, Event{x: x1 - 1, typ: 1, y1: y1, y2: y2, val: -1, idx: id})
			events = append(events, Event{x: x1 - 1, typ: 1, y1: 1, y2: y1 - 1, val: -1, idx: id})
			events = append(events, Event{x: 0, typ: 1, y1: 1, y2: y1 - 1, val: 1, idx: id})
			events = append(events, Event{x: x1 - 1, typ: 1, y1: y2 + 1, y2: m, val: -1, idx: id})
			events = append(events, Event{x: 0, typ: 1, y1: y2 + 1, y2: m, val: 1, idx: id})
			events = append(events, Event{x: n, typ: 1, y1: 1, y2: y1 - 1, val: -1, idx: id})
			events = append(events, Event{x: x2, typ: 1, y1: 1, y2: y1 - 1, val: 1, idx: id})
			events = append(events, Event{x: n, typ: 1, y1: y2 + 1, y2: m, val: -1, idx: id})
			events = append(events, Event{x: x2, typ: 1, y1: y2 + 1, y2: m, val: 1, idx: id})
		}
	}

	sort.Slice(events, func(i, j int) bool {
		if events[i].x != events[j].x {
			return events[i].x < events[j].x
		}
		return events[i].typ < events[j].typ
	})

	bit := newBIT(m + 2)
	ans := make([]int64, qcnt)
	for _, e := range events {
		if e.typ == 0 {
			bit.updateRange(e.y1, e.y2, e.val)
		} else {
			res := bit.queryRange(e.y1, e.y2)
			ans[e.idx] += e.val * res
		}
	}

	var out strings.Builder
	for i := 0; i < qcnt; i++ {
		if i > 0 {
			out.WriteByte('\n')
		}
		fmt.Fprint(&out, ans[i])
	}
	return out.String(), nil
}

var rawTestcases = []string{
	"6 5 1 1 2 2 5 2",
	"5 5 3 0 4 1 5 5 6 0 4 3 4 3 10 0 1 5 4 5 9",
	"4 6 3 0 3 1 3 1 4 0 2 6 3 6 1 1 3 4 3 5",
	"6 7 1 1 5 3 5 6",
	"3 7 4 1 2 4 2 4 0 3 2 3 2 6 1 2 4 2 6 0 2 6 2 6 3",
	"6 6 1 0 6 5 6 5 4",
	"5 6 1 0 2 4 2 5 4",
	"7 3 4 0 1 1 6 3 4 0 4 3 5 3 10 0 2 3 4 3 10 1 4 2 6 2",
	"3 7 4 1 2 4 2 5 0 2 3 2 3 1 1 2 3 2 6 1 2 3 2 6",
	"7 3 1 0 7 1 7 2 7",
	"7 4 2 0 1 1 7 4 7 1 3 2 6 3",
	"6 7 2 1 3 2 5 5 1 2 4 3 5",
	"4 7 2 0 1 1 2 3 6 1 3 3 3 6",
	"5 3 4 0 5 1 5 2 4 1 3 2 3 2 0 5 1 5 3 7 1 2 2 4 2",
	"4 4 3 0 3 2 3 4 9 0 4 1 4 3 3 0 2 4 4 4 7",
	"3 4 3 0 2 3 2 4 6 1 2 2 2 3 1 2 2 2 2",
	"4 4 4 0 3 3 3 4 1 0 4 3 4 3 9 1 2 3 2 3 0 2 3 2 4 3",
	"6 5 2 0 1 2 4 4 7 1 4 2 4 3",
	"5 7 2 0 4 7 5 7 6 1 3 4 3 4",
	"6 3 4 0 1 2 1 2 5 0 4 3 6 3 9 0 5 2 5 3 9 1 3 2 3 2",
	"3 7 2 1 2 6 2 6 1 2 4 2 4",
	"3 5 2 0 3 5 3 5 2 0 3 5 3 5 5",
	"5 7 1 0 1 3 4 7 6",
	"4 7 3 1 2 3 3 3 0 4 7 4 7 5 1 2 4 3 6",
	"4 5 3 1 3 2 3 3 1 3 4 3 4 1 2 2 3 2",
	"4 3 3 0 3 3 3 3 4 1 3 2 3 2 1 3 2 3 2",
	"6 3 2 0 2 3 6 3 2 1 2 2 3 2",
	"5 5 3 1 4 4 4 4 0 2 2 3 2 3 0 1 5 5 5 10",
	"7 6 4 1 5 4 5 4 0 5 5 6 6 7 0 4 4 5 5 9 1 2 4 4 5",
	"6 6 1 1 4 5 4 5",
	"3 3 3 0 3 3 3 3 1 0 2 1 3 3 1 1 2 2 2 2",
	"4 6 4 1 3 2 3 4 0 3 2 3 5 8 1 3 3 3 5 0 4 3 4 6 2",
	"3 4 1 1 2 3 2 3",
	"6 5 2 1 5 4 5 4 1 4 2 5 3",
	"3 6 2 1 2 5 2 5 1 2 4 2 4",
	"3 3 4 1 2 2 2 2 0 2 1 3 3 3 1 2 2 2 2 0 1 1 1 2 7",
	"7 7 3 0 3 5 6 6 5 1 3 6 3 6 1 6 5 6 6",
	"3 7 2 0 3 1 3 7 9 1 2 6 2 6",
	"7 6 4 0 4 2 6 2 6 1 4 4 6 4 1 2 2 2 3 1 6 2 6 2",
	"7 4 2 0 2 2 3 3 7 0 3 3 5 3 4",
	"7 3 3 1 4 2 4 2 0 4 2 5 3 3 0 7 3 7 3 9",
	"3 5 1 0 2 4 3 5 8",
	"4 6 4 1 2 3 3 4 0 3 1 4 3 6 1 3 5 3 5 0 2 4 4 6 7",
	"7 4 1 1 3 3 3 3",
	"7 3 1 1 4 2 6 2",
	"3 7 2 1 2 6 2 6 1 2 4 2 4",
	"4 5 1 1 2 2 3 4",
	"4 5 1 1 3 3 3 3",
	"5 7 1 1 4 4 4 4",
	"5 3 3 1 4 2 4 2 1 2 2 3 2 1 4 2 4 2",
	"7 4 1 1 5 3 6 3",
	"5 6 3 1 4 5 4 5 0 4 1 5 6 7 0 5 6 5 6 2",
	"5 5 2 0 2 2 4 3 6 0 5 1 5 1 3",
	"5 7 3 1 2 4 2 6 1 2 4 3 5 1 2 6 4 6",
	"3 5 4 1 2 2 2 3 1 2 4 2 4 0 1 2 1 5 6 0 2 5 3 5 9",
	"3 3 2 0 2 1 3 1 2 0 3 1 3 2 10",
	"3 6 2 1 2 3 2 4 0 3 2 3 5 7",
	"4 4 1 1 3 2 3 3",
	"5 5 2 1 4 4 4 4 0 3 1 4 1 2",
	"5 5 1 1 2 3 2 3",
	"3 6 4 0 2 2 2 3 3 0 3 4 3 5 9 0 2 2 2 6 2 1 2 2 2 5",
	"6 4 4 1 5 3 5 3 0 5 3 6 4 5 1 3 3 3 3 0 5 3 6 3 5",
	"6 4 4 0 3 1 5 3 1 1 2 2 2 2 1 4 3 5 3 1 2 3 5 3",
	"5 3 1 1 4 2 4 2",
	"7 7 1 0 5 4 7 6 2",
	"3 4 1 0 3 2 3 3 10",
	"5 6 2 0 5 1 5 2 3 1 4 2 4 3",
	"4 3 2 1 3 2 3 2 0 4 3 4 3 10",
	"4 3 3 1 3 2 3 2 0 1 3 3 3 6 1 2 2 3 2",
	"6 5 3 0 1 3 3 3 2 0 4 2 6 4 4 1 5 3 5 4",
	"6 3 4 1 2 2 3 2 1 2 2 5 2 0 5 1 5 2 7 1 5 2 5 2",
	"6 7 1 0 1 5 2 7 4",
	"4 5 4 0 3 2 4 2 5 0 2 2 2 2 1 1 3 4 3 4 1 3 2 3 2",
	"3 7 4 1 2 2 2 5 0 2 1 3 4 8 1 2 2 2 6 0 3 6 3 7 6",
	"5 4 1 1 4 2 4 2",
	"3 6 1 1 2 5 2 5",
	"3 7 2 0 2 3 2 7 9 0 1 1 1 5 6",
	"5 6 3 1 4 4 4 4 0 4 6 4 6 3 0 3 2 4 6 5",
	"5 5 3 0 5 4 5 4 10 1 3 4 3 4 1 4 3 4 3",
	"6 7 1 1 5 4 5 4",
	"3 7 4 1 2 3 2 6 1 2 3 2 5 0 2 5 2 6 5 1 2 3 2 5",
	"4 6 3 0 1 4 4 6 1 0 1 4 3 4 5 1 3 2 3 4",
	"4 3 4 1 3 2 3 2 0 4 1 4 2 3 0 4 1 4 1 4 1 3 2 3 2",
	"7 5 3 0 5 3 5 4 7 0 5 3 7 3 8 1 3 2 3 3",
	"6 4 3 1 2 2 2 3 0 4 4 4 4 9 1 3 2 3 3",
	"3 6 3 0 1 5 3 6 8 0 3 6 3 6 1 0 2 5 3 6 7",
	"7 3 4 0 4 3 6 3 9 1 5 2 5 2 1 5 2 6 2 1 4 2 4 2",
	"7 4 3 0 2 1 7 4 9 1 4 2 4 3 1 4 2 5 2",
	"3 4 4 0 1 4 2 4 2 1 2 3 2 3 0 1 2 2 3 3 0 1 2 3 4 1",
	"3 6 4 1 2 2 2 2 0 3 4 3 5 1 1 2 3 2 3 1 2 3 2 5",
	"7 5 2 0 2 1 4 1 5 1 5 4 6 4",
	"3 5 4 0 1 4 1 4 10 1 2 2 2 2 0 1 3 1 3 2 1 2 2 2 2",
	"6 4 3 1 2 2 3 3 1 4 2 5 2 1 4 3 5 3",
	"4 6 3 0 2 4 2 5 9 0 2 6 2 6 9 0 3 1 4 6 6",
	"4 6 2 1 2 3 3 4 1 3 5 3 5",
	"5 6 4 1 2 4 2 5 1 4 3 4 3 0 4 2 4 3 1 1 3 2 3 2",
	"5 4 2 1 3 3 4 3 1 3 2 4 3",
	"4 6 2 1 3 2 3 2 0 3 6 4 6 2",
	"5 5 1 0 4 1 4 4 8",
	"3 4 1 0 1 2 3 3 6",
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v, output: %s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: verifierE /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	for idx, tc := range rawTestcases {
		input := tc + "\n"
		expected, err := solve390E(input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d invalid: %v\n", idx+1, err)
			os.Exit(1)
		}
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != expected {
			fmt.Fprintf(os.Stderr, "case %d failed\nexpected: %s\n got: %s\n", idx+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(rawTestcases))
}

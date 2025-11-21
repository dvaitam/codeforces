package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type Box struct {
	items []int
	pos   map[int]int
}

type Move struct {
	from int
	to   int
	kind int
}

func (b *Box) remove(kind int) {
	idx := b.pos[kind]
	lastIdx := len(b.items) - 1
	if idx != lastIdx {
		lastKind := b.items[lastIdx]
		b.items[idx] = lastKind
		b.pos[lastKind] = idx
	}
	b.items = b.items[:lastIdx]
	delete(b.pos, kind)
}

func (b *Box) add(kind int) {
	b.pos[kind] = len(b.items)
	b.items = append(b.items, kind)
}

func findKind(boxes []Box, from, to int) int {
	dest := boxes[to].pos
	for _, kind := range boxes[from].items {
		if _, exists := dest[kind]; !exists {
			return kind
		}
	}
	return -1
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}

	boxes := make([]Box, n)
	total := 0
	for i := 0; i < n; i++ {
		var s int
		fmt.Fscan(in, &s)
		boxes[i].items = make([]int, s)
		boxes[i].pos = make(map[int]int, s)
		for j := 0; j < s; j++ {
			var val int
			fmt.Fscan(in, &val)
			val--
			boxes[i].items[j] = val
			boxes[i].pos[val] = j
		}
		total += s
	}

	base := total / n
	extra := total % n

	target := make([]int, n)
	for i := range target {
		target[i] = base
	}

	idx := make([]int, n)
	for i := range idx {
		idx[i] = i
	}

	sort.Slice(idx, func(i, j int) bool {
		if len(boxes[idx[i]].items) == len(boxes[idx[j]].items) {
			return idx[i] < idx[j]
		}
		return len(boxes[idx[i]].items) > len(boxes[idx[j]].items)
	})

	for i := 0; i < extra; i++ {
		target[idx[i]] = base + 1
	}

	type entry struct {
		id   int
		need int
	}

	surplus := make([]entry, 0)
	deficit := make([]entry, 0)

	for i := 0; i < n; i++ {
		cur := len(boxes[i].items)
		if cur > target[i] {
			surplus = append(surplus, entry{i, cur - target[i]})
		} else if cur < target[i] {
			deficit = append(deficit, entry{i, target[i] - cur})
		}
	}

	moves := make([]Move, 0)

	si, di := 0, 0
	for si < len(surplus) && di < len(deficit) {
		from := surplus[si].id
		to := deficit[di].id
		take := surplus[si].need
		if deficit[di].need < take {
			take = deficit[di].need
		}
		for t := 0; t < take; t++ {
			kind := findKind(boxes, from, to)
			if kind == -1 {
				panic("no valid move found")
			}
			boxes[from].remove(kind)
			boxes[to].add(kind)
			moves = append(moves, Move{from, to, kind})
		}
		surplus[si].need -= take
		deficit[di].need -= take
		if surplus[si].need == 0 {
			si++
		}
		if deficit[di].need == 0 {
			di++
		}
	}

	fmt.Fprintln(out, len(moves))
	for _, mv := range moves {
		fmt.Fprintf(out, "%d %d %d\n", mv.from+1, mv.to+1, mv.kind+1)
	}
}

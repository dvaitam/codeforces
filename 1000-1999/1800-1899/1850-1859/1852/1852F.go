package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type pair struct {
	s   int
	cnt int
}

type node struct {
	d   int
	arr []pair // sorted by s
}

var reds []node
var blues []node
var total int

// insert adds count c for (d,s) into the given slice of nodes.
func insert(list *[]node, d, s, c int) {
	nodes := *list
	idx := sort.Search(len(nodes), func(i int) bool { return nodes[i].d >= d })
	if idx == len(nodes) || nodes[idx].d != d {
		nodes = append(nodes, node{})
		copy(nodes[idx+1:], nodes[idx:])
		nodes[idx] = node{d: d, arr: []pair{{s, c}}}
	} else {
		arr := nodes[idx].arr
		j := sort.Search(len(arr), func(i int) bool { return arr[i].s >= s })
		if j < len(arr) && arr[j].s == s {
			arr[j].cnt += c
		} else {
			arr = append(arr, pair{})
			copy(arr[j+1:], arr[j:])
			arr[j] = pair{s, c}
		}
		nodes[idx].arr = arr
	}
	*list = nodes
}

// matchBlue tries to match as many pandas as possible for the blue node at idx.
func matchBlue(idx int) {
	if idx < 0 || idx >= len(blues) {
		return
	}
	b := &blues[idx]
	if len(b.arr) == 0 {
		blues = append(blues[:idx], blues[idx+1:]...)
		return
	}
	db := b.d
	sb := b.arr[0].s
	cnt := b.arr[0].cnt
	matched := 0
	for cnt > 0 {
		ridx := sort.Search(len(reds), func(i int) bool { return reds[i].d > db }) - 1
		found := false
		for ridx >= 0 {
			if ridx >= len(reds) {
				ridx = len(reds) - 1
			}
			rarr := reds[ridx].arr
			j := sort.Search(len(rarr), func(i int) bool { return rarr[i].s > sb }) - 1
			if j >= 0 && rarr[j].s <= sb {
				take := cnt
				if rarr[j].cnt < take {
					take = rarr[j].cnt
				}
				cnt -= take
				rarr[j].cnt -= take
				matched += take
				if rarr[j].cnt == 0 {
					rarr = append(rarr[:j], rarr[j+1:]...)
				}
				if len(rarr) == 0 {
					reds = append(reds[:ridx], reds[ridx+1:]...)
					ridx--
				} else {
					reds[ridx].arr = rarr
				}
				found = true
				if cnt == 0 {
					break
				}
			} else {
				ridx--
			}
		}
		if !found {
			break
		}
	}
	total += matched
	b.arr[0].cnt -= matched
	if b.arr[0].cnt == 0 {
		b.arr = b.arr[1:]
		if len(b.arr) == 0 {
			blues = append(blues[:idx], blues[idx+1:]...)
		}
	}
}

// tryMatch greedily pairs available red and blue pandas.
func tryMatch() {
	for len(blues) > 0 {
		idx := len(blues) - 1 // start from largest d
		old := total
		matchBlue(idx)
		if total == old {
			break
		}
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var q int
	fmt.Fscan(reader, &q)
	for i := 0; i < q; i++ {
		var x, t, c int
		fmt.Fscan(reader, &x, &t, &c)
		d := t - x
		s := t + x
		if c < 0 {
			insert(&reds, d, s, -c)
		} else {
			insert(&blues, d, s, c)
		}
		tryMatch()
		fmt.Fprintln(writer, total)
	}
}

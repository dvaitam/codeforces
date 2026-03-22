package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanWords)
	scanner.Buffer(make([]byte, 1024*1024), 10*1024*1024)

	scanInt := func() int {
		if !scanner.Scan() {
			return 0
		}
		s := scanner.Bytes()
		res := 0
		for _, b := range s {
			res = res*10 + int(b-'0')
		}
		return res
	}

	scanInt64 := func() int64 {
		if !scanner.Scan() {
			return 0
		}
		s := scanner.Bytes()
		res := int64(0)
		for _, b := range s {
			res = res*10 + int64(b-'0')
		}
		return res
	}

	k := scanInt()
	if k == 0 {
		return
	}
	n := scanInt()
	m := scanInt()

	a := make([]int64, k+1)
	for i := 1; i <= k; i++ {
		a[i] = scanInt64()
	}

	type Imp struct {
		id    int
		typ   int
		skill int
		val   int64
	}

	maxT1 := make([]Imp, k+1)
	t2List := make([][]Imp, k+1)
	var t3List []Imp

	for i := 1; i <= n; i++ {
		typ := scanInt()
		skill := scanInt()
		val := scanInt64()

		if typ == 1 {
			if val > maxT1[skill].val {
				maxT1[skill] = Imp{id: i, typ: 1, skill: skill, val: val}
			}
		} else if typ == 2 {
			t2List[skill] = append(t2List[skill], Imp{id: i, typ: 2, skill: skill, val: val})
		} else if typ == 3 {
			t3List = append(t3List, Imp{id: i, typ: 3, skill: skill, val: val})
		}
	}

	type Candidate struct {
		v         int64
		d         int64
		id        int
		orig_type int
	}

	var candidates []Candidate

	for i := 1; i <= k; i++ {
		if maxT1[i].val > a[i] {
			v := maxT1[i].val - a[i]
			t2List[i] = append(t2List[i], Imp{
				id:    maxT1[i].id,
				typ:   1,
				skill: i,
				val:   v,
			})
		}

		sort.Slice(t2List[i], func(x, y int) bool {
			return t2List[i][x].val > t2List[i][y].val
		})

		currentD := a[i]
		for _, imp := range t2List[i] {
			candidates = append(candidates, Candidate{
				v:         imp.val,
				d:         currentD,
				id:        imp.id,
				orig_type: imp.typ,
			})
			currentD += imp.val
		}
	}

	for _, imp := range t3List {
		candidates = append(candidates, Candidate{
			v:         imp.val - 1,
			d:         1,
			id:        imp.id,
			orig_type: 3,
		})
	}

	sort.Slice(candidates, func(i, j int) bool {
		v1 := candidates[i].v
		d1 := candidates[i].d
		v2 := candidates[j].v
		d2 := candidates[j].d
		return v1*d2 > v2*d1
	})

	take := m
	if take > len(candidates) {
		take = len(candidates)
	}

	var t1, t2, t3 []int
	realTake := 0

	for i := 0; i < take; i++ {
		c := candidates[i]
		if c.v == 0 {
			break
		}
		realTake++
		if c.orig_type == 1 {
			t1 = append(t1, c.id)
		} else if c.orig_type == 2 {
			t2 = append(t2, c.id)
		} else if c.orig_type == 3 {
			t3 = append(t3, c.id)
		}
	}

	outWriter := bufio.NewWriter(os.Stdout)
	defer outWriter.Flush()

	fmt.Fprintln(outWriter, realTake)
	if realTake > 0 {
		first := true
		for _, id := range t1 {
			if !first {
				fmt.Fprint(outWriter, " ")
			}
			fmt.Fprint(outWriter, id)
			first = false
		}
		for _, id := range t2 {
			if !first {
				fmt.Fprint(outWriter, " ")
			}
			fmt.Fprint(outWriter, id)
			first = false
		}
		for _, id := range t3 {
			if !first {
				fmt.Fprint(outWriter, " ")
			}
			fmt.Fprint(outWriter, id)
			first = false
		}
	}
	fmt.Fprintln(outWriter)
}

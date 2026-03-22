package main

import (
	"bytes"
	"fmt"
	"os"
	"sort"
)

type Portal struct {
	r, c int32
	s    int64
}

var (
	parent    []int32
	ends      [][2]int32
	ends_len  []int8
	Int       []int32
	rate      []int8
	last_time []int64
	E         []int8
)

func find(i int32) int32 {
	root := i
	for root != parent[root] {
		root = parent[root]
	}
	curr := i
	for curr != root {
		nxt := parent[curr]
		parent[curr] = root
		curr = nxt
	}
	return root
}

func merge_internal(fA, fB int32, S int64) {
	rA := find(fA)
	rB := find(fB)

	if rA == rB {
		for i := int8(0); i < ends_len[rA]; i++ {
			u := ends[rA][i]
			E[u] = int8((int64(E[u]) + int64(rate[u])*(S-last_time[u])) % 2)
			last_time[u] = S
		}
		ends_len[rA] = 0
		return
	}

	for i := int8(0); i < ends_len[rA]; i++ {
		u := ends[rA][i]
		E[u] = int8((int64(E[u]) + int64(rate[u])*(S-last_time[u])) % 2)
		last_time[u] = S
	}
	for i := int8(0); i < ends_len[rB]; i++ {
		u := ends[rB][i]
		E[u] = int8((int64(E[u]) + int64(rate[u])*(S-last_time[u])) % 2)
		last_time[u] = S
	}

	parent[rB] = rA
	Int[rA] += Int[rB] + 1

	new_len := int8(0)
	for i := int8(0); i < ends_len[rA]; i++ {
		u := ends[rA][i]
		if u != fA && u != fB {
			ends[rA][new_len] = u
			new_len++
		}
	}
	for i := int8(0); i < ends_len[rB]; i++ {
		u := ends[rB][i]
		if u != fA && u != fB {
			ends[rA][new_len] = u
			new_len++
		}
	}
	ends_len[rA] = new_len

	new_rate := int8(Int[rA] % 2)
	for i := int8(0); i < ends_len[rA]; i++ {
		u := ends[rA][i]
		rate[u] = new_rate
		last_time[u] = S
	}
}

func main() {
	content, _ := os.ReadFile("/dev/stdin")
	offset := 0
	nextInt := func() int {
		for offset < len(content) && (content[offset] < '0' || content[offset] > '9') {
			offset++
		}
		if offset >= len(content) {
			return 0
		}
		res := 0
		for offset < len(content) && content[offset] >= '0' && content[offset] <= '9' {
			res = res*10 + int(content[offset]-'0')
			offset++
		}
		return res
	}

	N := nextInt()
	if N == 0 {
		return
	}
	M := nextInt()

	totalFaces := int32(4 * N * M)
	parent = make([]int32, totalFaces)
	ends = make([][2]int32, totalFaces)
	ends_len = make([]int8, totalFaces)
	Int = make([]int32, totalFaces)
	rate = make([]int8, totalFaces)
	last_time = make([]int64, totalFaces)
	E = make([]int8, totalFaces)

	portals := make([]Portal, 0, N*M)

	for r := int32(0); r < int32(N); r++ {
		for c := int32(0); c < int32(M); c++ {
			s := int64(nextInt())
			portals = append(portals, Portal{r, c, s})
			for k := int32(0); k < 4; k++ {
				f := 4*(r*int32(M)+c) + k
				parent[f] = f
				ends[f][0] = f
				ends_len[f] = 1
				last_time[f] = 1
			}
		}
	}

	for r := int32(0); r < int32(N); r++ {
		for c := int32(0); c < int32(M); c++ {
			f := 4 * (r*int32(M) + c)
			if r > 0 {
				f2 := 4*((r-1)*int32(M)+c) + 2
				rA := f
				rB := f2
				parent[rB] = rA
				ends[rA][0] = f
				ends[rA][1] = f2
				ends_len[rA] = 2
			}
			if c > 0 {
				f3 := 4*(r*int32(M)+c) + 3
				f4 := 4*(r*int32(M)+c-1) + 1
				rA := f3
				rB := f4
				parent[rB] = rA
				ends[rA][0] = f3
				ends[rA][1] = f4
				ends_len[rA] = 2
			}
		}
	}

	sort.Slice(portals, func(i, j int) bool {
		return portals[i].s < portals[j].s
	})

	ans := make([][]byte, N)
	for i := 0; i < N; i++ {
		ans[i] = make([]byte, M)
	}

	i := 0
	for i < len(portals) {
		j := i
		for j < len(portals) && portals[j].s == portals[i].s {
			j++
		}

		group := portals[i:j]
		S := group[0].s

		for _, p := range group {
			for k := int32(0); k < 4; k++ {
				f := 4*(p.r*int32(M)+p.c) + k
				E[f] = int8((int64(E[f]) + int64(rate[f])*(S-last_time[f])) % 2)
				last_time[f] = S
			}
		}

		for _, p := range group {
			f0 := 4*(p.r*int32(M)+p.c) + 0
			f1 := 4*(p.r*int32(M)+p.c) + 1
			f2 := 4*(p.r*int32(M)+p.c) + 2
			f3 := 4*(p.r*int32(M)+p.c) + 3

			sum := E[f0] + E[f1] + E[f2] + E[f3]
			t := sum % 2
			ans[p.r][p.c] = byte('0' + t)

			if t == 0 {
				merge_internal(f0, f2, S)
				merge_internal(f1, f3, S)
			} else {
				merge_internal(f0, f3, S)
				merge_internal(f2, f1, S)
			}
		}

		i = j
	}

	var out bytes.Buffer
	for r := 0; r < N; r++ {
		out.Write(ans[r])
		out.WriteByte('\n')
	}
	fmt.Print(out.String())
}
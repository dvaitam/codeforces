package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

// We will search for the best clique among jarls that appear in special pairs.
// Other jarls act identically within their islands and we can always choose one
// of them if their island is not represented in the clique.

type Pair struct{ a, b int }

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var p, n int
	fmt.Fscan(in, &p, &n)
	island := make([]int, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(in, &island[i])
	}
	var k int
	fmt.Fscan(in, &k)

	pairs := make([]Pair, k)
	special := make([]bool, n+1)
	for i := 0; i < k; i++ {
		fmt.Fscan(in, &pairs[i].a, &pairs[i].b)
		special[pairs[i].a] = true
		special[pairs[i].b] = true
	}

	// collect all jarls that appear in at least one pair
	specialID := []int{}
	idToIdx := make(map[int]int)
	for i := 1; i <= n; i++ {
		if special[i] {
			idToIdx[i] = len(specialID)
			specialID = append(specialID, i)
		}
	}
	m := len(specialID)

	// prepare information about generic jarls per island
	genericRep := make([]int, p+1)
	for i := 1; i <= p; i++ {
		genericRep[i] = -1
	}
	for i := 1; i <= n; i++ {
		if !special[i] {
			if genericRep[island[i]] == -1 {
				genericRep[island[i]] = i
			}
		}
	}
	// map generic islands to bit indices
	genIndex := make(map[int]int)
	gi := 0
	for i := 1; i <= p; i++ {
		if genericRep[i] != -1 {
			genIndex[i] = gi
			gi++
		}
	}

	// for each special jarl store mask of its island if that island has generic jarls
	maskGeneric := make([]uint64, m)
	for idx, id := range specialID {
		if v, ok := genIndex[island[id]]; ok {
			maskGeneric[idx] = 1 << uint(v)
		}
	}

	// build base adjacency for special jarls
	adj := make([]uint64, m)
	for i := 0; i < m; i++ {
		for j := 0; j < m; j++ {
			if i == j {
				continue
			}
			if island[specialID[i]] != island[specialID[j]] {
				adj[i] |= 1 << uint(j)
			}
		}
	}
	// apply modifications
	for _, pr := range pairs {
		i := idToIdx[pr.a]
		j := idToIdx[pr.b]
		if island[pr.a] == island[pr.b] {
			adj[i] |= 1 << uint(j)
			adj[j] |= 1 << uint(i)
		} else {
			adj[i] &^= 1 << uint(j)
			adj[j] &^= 1 << uint(i)
		}
	}

	// Bron–Kerbosch with pivoting to find clique of maximal (|S| - usedGenericIslands)
	var bestMask uint64
	bestValue := -1
	var dfs func(uint64, uint64, uint64, uint64)
	dfs = func(R, P, X, genMask uint64) {
		if P == 0 && X == 0 {
			sz := bits.OnesCount64(R)
			val := int(sz) - bits.OnesCount64(genMask)
			if val > bestValue {
				bestValue = val
				bestMask = R
			}
			return
		}
		if P == 0 {
			return
		}
		// pruning by maximum possible size
		if int(bits.OnesCount64(R)+bits.OnesCount64(P)) <= bestValue {
			return
		}
		// choose pivot
		union := P | X
		var pivot int
		if union != 0 {
			pivot = bits.TrailingZeros64(union)
		} else {
			pivot = 0
		}
		candidates := P &^ adj[pivot]
		for candidates != 0 {
			v := bits.TrailingZeros64(candidates)
			candidates &^= 1 << uint(v)
			dfs(R|1<<uint(v), P&adj[v], X&adj[v], genMask|maskGeneric[v])
			P &^= 1 << uint(v)
			X |= 1 << uint(v)
		}
	}

	dfs(0, (uint64(1)<<uint(m))-1, 0, 0)

	// Determine which islands with generic jarls are used by the best clique
	usedGeneric := make([]bool, gi)
	for i := 0; i < m; i++ {
		if bestMask&(1<<uint(i)) != 0 {
			if maskGeneric[i] != 0 {
				idx := bits.TrailingZeros64(maskGeneric[i])
				usedGeneric[idx] = true
			}
		}
	}

	// Compose answer list
	ans := []int{}
	// add special jarls from best clique
	for i := 0; i < m; i++ {
		if bestMask&(1<<uint(i)) != 0 {
			ans = append(ans, specialID[i])
		}
	}
	// add one generic jarl from each island that has one and is not used by special jarl
	for isl, idx := range genIndex {
		if !usedGeneric[idx] {
			ans = append(ans, genericRep[isl])
		}
	}

	fmt.Fprintln(out, len(ans))
	for i, v := range ans {
		if i > 0 {
			fmt.Fprint(out, " ")
		}
		fmt.Fprint(out, v)
	}
	if len(ans) > 0 {
		fmt.Fprintln(out)
	}
}

package main

import (
	"bufio"
	"fmt"
	"io"
	"math/bits"
	"os"
)

type FastScanner struct {
	data []byte
	idx  int
}

func (fs *FastScanner) NextInt() int {
	n := len(fs.data)
	for fs.idx < n && (fs.data[fs.idx] < '0' || fs.data[fs.idx] > '9') {
		fs.idx++
	}
	val := 0
	for fs.idx < n && fs.data[fs.idx] >= '0' && fs.data[fs.idx] <= '9' {
		val = val*10 + int(fs.data[fs.idx]-'0')
		fs.idx++
	}
	return val
}

type Pair struct {
	a, b int
}

type IslandData struct {
	islandID   int
	normal     bool
	normalRep  int
	excVerts   []int
	globalIDs  []int
	bestVal    []int
	bestChoice []uint64
}

func computeLocalMask(isl *IslandData, s int, zPosByGlobal []int, neighborMaskZ []int) int {
	mask := 0
	for pos, gid := range isl.globalIDs {
		zp := zPosByGlobal[gid]
		if zp >= 0 {
			if s&(1<<uint(zp)) != 0 {
				mask |= 1 << uint(pos)
			}
		} else {
			if neighborMaskZ[gid]&s == 0 {
				mask |= 1 << uint(pos)
			}
		}
	}
	return mask
}

func main() {
	data, _ := io.ReadAll(os.Stdin)
	fs := FastScanner{data: data}

	p := fs.NextInt()
	n := fs.NextInt()

	islandOf := make([]int, n+1)
	islandSize := make([]int, p)
	for i := 1; i <= n; i++ {
		islandOf[i] = fs.NextInt() - 1
		islandSize[islandOf[i]]++
	}

	k := fs.NextInt()

	exceptional := make([]bool, n+1)
	intraPairsByIsland := make([][]Pair, p)
	crossPairs := make([]Pair, 0, k)

	for i := 0; i < k; i++ {
		a := fs.NextInt()
		b := fs.NextInt()
		exceptional[a] = true
		exceptional[b] = true
		if islandOf[a] == islandOf[b] {
			intraPairsByIsland[islandOf[a]] = append(intraPairsByIsland[islandOf[a]], Pair{a, b})
		} else {
			crossPairs = append(crossPairs, Pair{a, b})
		}
	}

	exceptionalCount := make([]int, p)
	normalRep := make([]int, p)
	for i := 0; i < p; i++ {
		normalRep[i] = -1
	}
	excVertsByIsland := make([][]int, p)

	for i := 1; i <= n; i++ {
		isl := islandOf[i]
		if exceptional[i] {
			exceptionalCount[isl]++
			excVertsByIsland[isl] = append(excVertsByIsland[isl], i)
		} else if normalRep[isl] == -1 {
			normalRep[isl] = i
		}
	}

	normalExists := make([]bool, p)
	base := 0
	for i := 0; i < p; i++ {
		if islandSize[i] > exceptionalCount[i] {
			normalExists[i] = true
			base++
		}
	}

	islands := make([]IslandData, 0)
	for isl := 0; isl < p; isl++ {
		if len(excVertsByIsland[isl]) > 0 {
			islands = append(islands, IslandData{
				islandID:  isl,
				normal:    normalExists[isl],
				normalRep: normalRep[isl],
				excVerts:  excVertsByIsland[isl],
			})
		}
	}

	excLocalPos := make([]int, n+1)
	for i := 0; i <= n; i++ {
		excLocalPos[i] = -1
	}
	for idx := range islands {
		for li, v := range islands[idx].excVerts {
			excLocalPos[v] = li
		}
	}

	globalFlag := make([]bool, n+1)
	for _, pr := range crossPairs {
		globalFlag[pr.a] = true
		globalFlag[pr.b] = true
	}

	globalIDByVertex := make([]int, n+1)
	for i := 0; i <= n; i++ {
		globalIDByVertex[i] = -1
	}
	g := 0
	for idx := range islands {
		for _, v := range islands[idx].excVerts {
			if globalFlag[v] {
				globalIDByVertex[v] = g
				g++
			}
		}
	}

	hedges := make([][2]int, len(crossPairs))
	for i, pr := range crossPairs {
		hedges[i] = [2]int{globalIDByVertex[pr.a], globalIDByVertex[pr.b]}
	}

	inZ := make([]bool, g)
	for _, e := range hedges {
		inZ[e[0]] = true
	}

	zPosByGlobal := make([]int, g)
	for i := 0; i < g; i++ {
		zPosByGlobal[i] = -1
	}
	z := 0
	for gid := 0; gid < g; gid++ {
		if inZ[gid] {
			zPosByGlobal[gid] = z
			z++
		}
	}

	adjZ := make([]int, z)
	neighborMaskZ := make([]int, g)
	for _, e := range hedges {
		u, v := e[0], e[1]
		zu, zv := zPosByGlobal[u], zPosByGlobal[v]
		if zu >= 0 && zv >= 0 {
			adjZ[zu] |= 1 << uint(zv)
			adjZ[zv] |= 1 << uint(zu)
		} else if zu >= 0 && zv < 0 {
			neighborMaskZ[v] |= 1 << uint(zu)
		} else if zu < 0 && zv >= 0 {
			neighborMaskZ[u] |= 1 << uint(zv)
		}
	}

	for idx := range islands {
		isl := &islands[idx]
		t := len(isl.excVerts)
		adj := make([]uint64, t)
		for _, pr := range intraPairsByIsland[isl.islandID] {
			a := excLocalPos[pr.a]
			b := excLocalPos[pr.b]
			adj[a] |= uint64(1) << uint(b)
			adj[b] |= uint64(1) << uint(a)
		}

		globalIDs := make([]int, 0)
		gBitByLocalExc := make([]int, t)
		for li, v := range isl.excVerts {
			gid := globalIDByVertex[v]
			if gid >= 0 {
				pos := len(globalIDs)
				globalIDs = append(globalIDs, gid)
				gBitByLocalExc[li] = 1 << uint(pos)
			}
		}

		size := 1 << uint(len(globalIDs))
		bestVal := make([]int, size)
		bestChoice := make([]uint64, size)

		var dfs func(uint64, uint64, int, int)
		dfs = func(cand uint64, currSet uint64, currGMask int, currSize int) {
			if currSize > bestVal[currGMask] {
				bestVal[currGMask] = currSize
				bestChoice[currGMask] = currSet
			}
			bitsLeft := cand
			for bitsLeft != 0 {
				v := bits.TrailingZeros64(bitsLeft)
				bitsLeft &= bitsLeft - 1
				dfs(bitsLeft&adj[v], currSet|(uint64(1)<<uint(v)), currGMask|gBitByLocalExc[v], currSize+1)
			}
		}

		allBits := (uint64(1) << uint(t)) - 1
		dfs(allBits, 0, 0, 0)

		gi := len(globalIDs)
		for b := 0; b < gi; b++ {
			bit := 1 << uint(b)
			for mask := 0; mask < size; mask++ {
				if mask&bit != 0 && bestVal[mask^bit] > bestVal[mask] {
					bestVal[mask] = bestVal[mask^bit]
					bestChoice[mask] = bestChoice[mask^bit]
				}
			}
		}

		isl.globalIDs = globalIDs
		isl.bestVal = bestVal
		isl.bestChoice = bestChoice
	}

	limit := 1 << uint(z)
	independent := make([]bool, limit)
	independent[0] = true
	for mask := 1; mask < limit; mask++ {
		v := bits.TrailingZeros(uint(mask))
		prev := mask & (mask - 1)
		independent[mask] = independent[prev] && (adjZ[v]&prev == 0)
	}

	bestTotal := -1
	bestS := 0

	for sMask := 0; sMask < limit; sMask++ {
		if !independent[sMask] {
			continue
		}
		total := base
		for idx := range islands {
			isl := &islands[idx]
			localMask := computeLocalMask(isl, sMask, zPosByGlobal, neighborMaskZ)
			size := isl.bestVal[localMask]
			if isl.normal {
				if size > 1 {
					total += size - 1
				}
			} else {
				total += size
			}
		}
		if total > bestTotal {
			bestTotal = total
			bestS = sMask
		}
	}

	chosen := make([]int, 0, bestTotal)
	excSelectedOnIsland := make([]bool, p)

	for idx := range islands {
		isl := &islands[idx]
		localMask := computeLocalMask(isl, bestS, zPosByGlobal, neighborMaskZ)
		size := isl.bestVal[localMask]
		if isl.normal && size <= 1 {
			continue
		}
		choice := isl.bestChoice[localMask]
		if choice != 0 {
			excSelectedOnIsland[isl.islandID] = true
		}
		for choice != 0 {
			v := bits.TrailingZeros64(choice)
			choice &= choice - 1
			chosen = append(chosen, isl.excVerts[v])
		}
	}

	for isl := 0; isl < p; isl++ {
		if normalExists[isl] && !excSelectedOnIsland[isl] {
			chosen = append(chosen, normalRep[isl])
		}
	}

	out := bufio.NewWriterSize(os.Stdout, 1<<20)
	fmt.Fprintln(out, len(chosen))
	for i, v := range chosen {
		if i > 0 {
			fmt.Fprint(out, " ")
		}
		fmt.Fprint(out, v)
	}
	fmt.Fprintln(out)
	out.Flush()
}

package main

import (
	"bufio"
	"fmt"
	"os"
)

const MAXN = 100005
const LOGN = 18
const BASE = 131

var (
	reader *bufio.Reader
	writer *bufio.Writer
)

type Edge struct {
	to   int
	char byte
}

var (
	adj1    = make([][]Edge, MAXN)
	adj2    = make([][]Edge, MAXN)
	ops     = make([]Operation, MAXN)
	map1    = make([]int, MAXN) // T1 node -> Trie node
	map2    = make([]int, MAXN) // T2 node -> Trie node

	// Trie
	trieNext = make([][26]int, MAXN)
	trieCnt  = 1
	// trieAdj is implicit via structure, but we need DFS order
	// so we can use trieNext for that.

	// Hashes and LCA
	powB      = make([]uint64, MAXN)
	trieHash  = make([]uint64, MAXN)
	hashMap   = make(map[uint64]int)
	
	up2      = make([][LOGN]int, MAXN)
	upHash2  = make([]uint64, MAXN)
	depth2   = make([]int, MAXN)
	
	// BITs
	bit1     = make([]int64, MAXN) // Subtree Sum K
	bit2     = make([]int64, MAXN) // Path Sum W (Range Update Point Query)
	tin      = make([]int, MAXN)
	tout     = make([]int, MAXN)
	timer    int
)

type Operation struct {
	t, v int
	c    byte
}

func main() {
	reader = bufio.NewReader(os.Stdin)
	writer = bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var numOps int
	fmt.Fscan(reader, &numOps)
	
	cnt1 := 1
	cnt2 := 1

	for i := 0; i < numOps; i++ {
		var t, v int
		var s string
		fmt.Fscan(reader, &t, &v, &s)
		c := s[0]
		ops[i] = Operation{t, v, c}
		if t == 1 {
			cnt1++
			adj1[v] = append(adj1[v], Edge{cnt1, c})
		} else {
			cnt2++
			adj2[v] = append(adj2[v], Edge{cnt2, c})
		}
	}

	// Build Trie from T1
	buildTrie(1, 1) // T1 root is 1, Trie root is 1
	
	// Precompute powers
	powB[0] = 1
	for i := 1; i < MAXN; i++ {
		powB[i] = powB[i-1] * BASE
	}

	// Trie DFS for hashing and BIT range
	dfsTrie(1, 0)

	// T2 processing
	dfsT2(1, 0, 0)
	
	// Precompute map2 for all T2 nodes
	for i := 1; i <= cnt2; i++ {
		map2[i] = binarySearchMatch(i)
	}

	// Initial state setup
	ans := int64(1)
	updateBIT2(tin[1], tout[1], 1) // Activate root 1 of T1
	updateBIT1(tin[1], 1)          // Activate root 1 of T2

	currM1 := 1
	currM2 := 1
	
	for i := 0; i < numOps; i++ {
		op := ops[i]
		if op.t == 1 {
			currM1++
			u := currM1
			tNode := map1[u]
			// Activate u in T1
			// Update BIT2: add 1 to subtree of tNode in Trie (effectively path sum increase)
			updateBIT2(tin[tNode], tout[tNode], 1)
			// Query BIT1: sum K in subtree of tNode
			term := queryBIT1(tout[tNode]) - queryBIT1(tin[tNode]-1)
			ans += term
		} else {
			currM2++
			v := currM2
			tNode := map2[v]
			// Activate v in T2
			// Update BIT1: add 1 to point tin[tNode]
			updateBIT1(tin[tNode], 1)
			// Query BIT2: value at tin[tNode]
			term := queryBIT2(tin[tNode])
			ans += term
		}
		fmt.Fprintln(writer, ans)
	}
}

// Build Trie from T1 by DFS
func buildTrie(u int, trieNode int) {
	map1[u] = trieNode
	for _, e := range adj1[u] {
		cIdx := int(e.char - 'a')
		if trieNext[trieNode][cIdx] == 0 {
			trieCnt++
			trieNext[trieNode][cIdx] = trieCnt
		}
		buildTrie(e.to, trieNext[trieNode][cIdx])
	}
}

func dfsTrie(u int, d int) {
	timer++
	tin[u] = timer
	hashMap[trieHash[u]] = u
	
	for c := 0; c < 26; c++ {
		v := trieNext[u][c]
		if v != 0 {
			trieHash[v] = trieHash[u] + uint64(c+'a')*powB[d]
			dfsTrie(v, d+1)
		}
	}
	tout[u] = timer
}

func dfsT2(u int, p int, d int) {
	depth2[u] = d
	up2[u][0] = p
	for j := 1; j < LOGN; j++ {
		up2[u][j] = up2[up2[u][j-1]][j-1]
	}
	for _, e := range adj2[u] {
		upHash2[e.to] = uint64(e.char) + upHash2[u] * BASE
		dfsT2(e.to, u, d+1)
	}
}

func binarySearchMatch(v int) int {
	low, high := 0, depth2[v]
	bestNode := 1 
	
	for low <= high {
		mid := (low + high) / 2
		if mid == 0 {
			low = 1
			continue
		}
		
		ancestor := getKthAncestor(v, mid)
		currentH := upHash2[v] - upHash2[ancestor] * powB[mid]
		
		if node, ok := hashMap[currentH]; ok {
			bestNode = node
			low = mid + 1
		} else {
			high = mid - 1
		}
	}
	return bestNode
}

func getKthAncestor(node, k int) int {
	for i := 0; i < LOGN; i++ {
		if (k>>i)&1 == 1 {
			node = up2[node][i]
		}
	}
	return node
}

func updateBIT1(idx int, val int64) {
	for ; idx <= trieCnt; idx += idx & -idx {
		bit1[idx] += val
	}
}

func queryBIT1(idx int) int64 {
	sum := int64(0)
	for ; idx > 0; idx -= idx & -idx {
		sum += bit1[idx]
	}
	return sum
}

func updateBIT2(l, r int, val int64) {
	internalUpdateBIT2(l, val)
	internalUpdateBIT2(r+1, -val)
}

func internalUpdateBIT2(idx int, val int64) {
	for ; idx <= trieCnt; idx += idx & -idx {
		bit2[idx] += val
	}
}

func queryBIT2(idx int) int64 {
	sum := int64(0)
	for ; idx > 0; idx -= idx & -idx {
		sum += bit2[idx]
	}
	return sum
}
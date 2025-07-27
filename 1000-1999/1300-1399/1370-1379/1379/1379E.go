package main

import (
	"bufio"
	"fmt"
	"os"
)

func isPowerOfTwo(x int) bool {
	return x > 0 && (x&(x-1)) == 0
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n, k int
	if _, err := fmt.Fscanf(reader, "%d %d", &n, &k); err != nil {
		return
	}
	if n%2 == 0 {
		fmt.Println("NO")
		return
	}
	t := (n - 1) / 2
	// maximum imbalanced achievable is t-1, minimum is 0
	if k < 0 || k > t-1 {
		if n == 1 && k == 0 {
			fmt.Println("YES")
			fmt.Println(0)
			return
		}
		fmt.Println("NO")
		return
	}
	// Special case k == t-1: simple chain
	if k == t-1 {
		s := make([]int, n+1)
		// leaves 1..t+1, internal nodes t+2..n
		curLeaf := 1
		nextNode := t + 2
		// first merge: leaves 1,2 -> node t+2
		s[1], s[2] = nextNode, nextNode
		last := nextNode
		nextNode++
		// remaining merges
		for i := 0; i < t-1; i++ {
			curLeaf++
			s[last] = nextNode
			s[curLeaf] = nextNode
			last = nextNode
			nextNode++
		}
		// last node has no child
		s[last] = 0
		fmt.Println("YES")
		for i := 1; i <= n; i++ {
			fmt.Printf("%d ", s[i])
		}
		fmt.Println()
		return
	}
	// Balanced perfect subtree for k==0 and t+1 is power of two leaves
	if k == 0 && isPowerOfTwo(t+1) {
		s := make([]int, n+1)
		nextNode := t + 2
		var build func(leaves []int) int
		build = func(leaves []int) int {
			if len(leaves) == 1 {
				return leaves[0]
			}
			var nextLeaves []int
			for i := 0; i < len(leaves); i += 2 {
				u, v := leaves[i], leaves[i+1]
				curr := nextNode
				nextNode++
				s[u], s[v] = curr, curr
				nextLeaves = append(nextLeaves, curr)
			}
			return build(nextLeaves)
		}
		leaves := make([]int, t+1)
		for i := 0; i < t+1; i++ {
			leaves[i] = i + 1
		}
		root := build(leaves)
		s[root] = 0
		fmt.Println("YES")
		for i := 1; i <= n; i++ {
			fmt.Printf("%d ", s[i])
		}
		fmt.Println()
		return
	}
	fmt.Println("NO")
}

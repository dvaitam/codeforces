package main

import (
	"bufio"
	"fmt"
	"math/big"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	f := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &f[i])
		f[i]--
	}

	state := make([]int, n)
	isCycle := make([]bool, n)
	cycleLengths := make([]int, 0)

	stack := make([]int, 0)
	for i := 0; i < n; i++ {
		if state[i] != 0 {
			continue
		}
		stack = stack[:0]
		cur := i
		for state[cur] == 0 {
			state[cur] = 1
			stack = append(stack, cur)
			cur = f[cur]
		}
		if state[cur] == 1 {
			idx := len(stack) - 1
			for idx >= 0 && stack[idx] != cur {
				idx--
			}
			if idx >= 0 {
				length := len(stack) - idx
				cycleLengths = append(cycleLengths, length)
				for j := idx; j < len(stack); j++ {
					isCycle[stack[j]] = true
				}
			}
		}
		for _, v := range stack {
			state[v] = 2
		}
	}

	depth := make([]int, n)
	for i := 0; i < n; i++ {
		depth[i] = -1
	}
	for i := 0; i < n; i++ {
		if isCycle[i] {
			depth[i] = 0
		}
	}

	var dfsDepth func(int) int
	dfsDepth = func(v int) int {
		if depth[v] != -1 {
			return depth[v]
		}
		depth[v] = dfsDepth(f[v]) + 1
		return depth[v]
	}

	maxDepth := 0
	for i := 0; i < n; i++ {
		if depth[i] == -1 {
			d := dfsDepth(i)
			if d > maxDepth {
				maxDepth = d
			}
		}
	}

	lcmVal := big.NewInt(1)
	for _, length := range cycleLengths {
		bigLen := big.NewInt(int64(length))
		gcdVal := new(big.Int).GCD(nil, nil, lcmVal, bigLen)
		tmp := new(big.Int).Div(new(big.Int).Set(lcmVal), gcdVal)
		tmp.Mul(tmp, bigLen)
		lcmVal = tmp
	}

	if len(cycleLengths) == 0 {
		lcmVal = big.NewInt(1)
	}

	result := new(big.Int).Set(lcmVal)
	bigMax := big.NewInt(int64(maxDepth))
	if result.Cmp(bigMax) < 0 {
		lcmInt := result.Int64()
		if lcmInt == 0 {
			lcmInt = 1
		}
		need := (int64(maxDepth) + lcmInt - 1) / lcmInt
		result.Mul(result, big.NewInt(need))
	}

	fmt.Fprintln(out, result)
}

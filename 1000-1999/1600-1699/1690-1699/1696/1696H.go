package main

import (
	"bufio"
	"fmt"
	"math/big"
	"os"
)

const mod = 1000000007

func maxBig(a, b *big.Int) *big.Int {
	if a == nil {
		if b == nil {
			return nil
		}
		c := new(big.Int).Set(b)
		return c
	}
	if b == nil {
		c := new(big.Int).Set(a)
		return c
	}
	if a.Cmp(b) >= 0 {
		c := new(big.Int).Set(a)
		return c
	}
	c := new(big.Int).Set(b)
	return c
}

func maxProduct(nums []int64, k int) *big.Int {
	dpPos := make([]*big.Int, k+1)
	dpNeg := make([]*big.Int, k+1)
	dpPos[0] = big.NewInt(1)
	for _, v := range nums {
		for j := k; j >= 1; j-- {
			val := big.NewInt(v)
			if v >= 0 {
				if dpPos[j-1] != nil {
					t := new(big.Int).Mul(dpPos[j-1], val)
					dpPos[j] = maxBig(dpPos[j], t)
				}
				if dpNeg[j-1] != nil {
					t := new(big.Int).Mul(dpNeg[j-1], val)
					dpNeg[j] = maxBig(dpNeg[j], t)
				}
			} else {
				if dpPos[j-1] != nil {
					t := new(big.Int).Mul(dpPos[j-1], val)
					dpNeg[j] = maxBig(dpNeg[j], t)
				}
				if dpNeg[j-1] != nil {
					t := new(big.Int).Mul(dpNeg[j-1], val)
					dpPos[j] = maxBig(dpPos[j], t)
				}
			}
		}
	}
	res := maxBig(dpPos[k], dpNeg[k])
	if res == nil {
		return big.NewInt(0)
	}
	return res
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n, k int
	fmt.Fscan(reader, &n, &k)
	a := make([]int64, n)
	for i := range a {
		fmt.Fscan(reader, &a[i])
	}

	var ans big.Int
	subset := make([]int64, 0, n)

	var dfs func(int)
	dfs = func(idx int) {
		if idx == n {
			if len(subset) >= k {
				prod := maxProduct(subset, k)
				ans.Add(&ans, prod)
			}
			return
		}
		dfs(idx + 1)
		subset = append(subset, a[idx])
		dfs(idx + 1)
		subset = subset[:len(subset)-1]
	}
	dfs(0)

	modBig := big.NewInt(mod)
	ans.Mod(&ans, modBig)
	fmt.Println(&ans)
}

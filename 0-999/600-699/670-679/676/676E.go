package main

import (
	"bufio"
	"fmt"
	"math/big"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	var k int
	if _, err := fmt.Fscan(in, &n, &k); err != nil {
		return
	}
	coeff := make([]int, n+1)
	unknown := make([]bool, n+1)
	knownCount := 0
	for i := 0; i <= n; i++ {
		var s string
		fmt.Fscan(in, &s)
		if s == "?" {
			unknown[i] = true
		} else {
			var val int
			fmt.Sscan(s, &val)
			coeff[i] = val
			knownCount++
		}
	}

	// Determine whose turn it is: computer starts first
	// After knownCount moves have been made
	computerTurn := knownCount%2 == 0
	movesLeft := n + 1 - knownCount

	// Special case k == 0
	if k == 0 {
		if !unknown[0] {
			if coeff[0] == 0 {
				fmt.Println("Yes")
			} else {
				fmt.Println("No")
			}
			return
		}
		if computerTurn {
			fmt.Println("No")
		} else {
			fmt.Println("Yes")
		}
		return
	}

	if movesLeft == 0 {
		// All coefficients are known, evaluate polynomial at k
		res := big.NewInt(0)
		kBig := big.NewInt(int64(k))
		tmp := new(big.Int)
		for i := n; i >= 0; i-- {
			res.Mul(res, kBig)
			tmp.SetInt64(int64(coeff[i]))
			res.Add(res, tmp)
		}
		if res.Sign() == 0 {
			fmt.Println("Yes")
		} else {
			fmt.Println("No")
		}
		return
	}

	// There is at least one unknown and k != 0
	// Last mover determines the result
	humanLast := false
	if computerTurn {
		humanLast = movesLeft%2 == 0
	} else {
		humanLast = movesLeft%2 == 1
	}
	if humanLast {
		fmt.Println("Yes")
	} else {
		fmt.Println("No")
	}
}

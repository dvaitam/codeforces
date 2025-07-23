package main

import (
	"bufio"
	"fmt"
	"math/big"
	"os"
)

// evaluate an expression defined by digits and operators using
// standard precedence (* before +).
func evalExpr(digits []*big.Int, ops []byte) *big.Int {
	if len(digits) == 0 {
		return big.NewInt(0)
	}
	cur := new(big.Int).Set(digits[0])
	res := big.NewInt(0)
	for i, op := range ops {
		if op == '+' {
			res.Add(res, cur)
			cur = new(big.Int).Set(digits[i+1])
		} else { // '*'
			cur.Mul(cur, digits[i+1])
		}
	}
	res.Add(res, cur)
	return res
}

// evaluate the whole expression after putting parentheses around
// digits[l:r+1]. The digits slice contains big.Int numbers for each
// original digit, ops contains operators between them.
func evalWithParen(digits []*big.Int, ops []byte, l, r int) *big.Int {
	// compute value of the sub-expression inside parentheses
	mid := evalExpr(digits[l:r+1], ops[l:r])

	// build new slices with the sub-expression replaced by 'mid'
	newDigits := make([]*big.Int, 0, len(digits)-(r-l))
	newDigits = append(newDigits, digits[:l]...)
	newDigits = append(newDigits, mid)
	if r+1 < len(digits) {
		newDigits = append(newDigits, digits[r+1:]...)
	}

	newOps := make([]byte, 0, len(ops)-(r-l))
	newOps = append(newOps, ops[:l]...)
	if r < len(ops) {
		newOps = append(newOps, ops[r])
		if r+1 < len(ops) {
			newOps = append(newOps, ops[r+1:]...)
		}
	}

	return evalExpr(newDigits, newOps)
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var s string
	if _, err := fmt.Fscan(reader, &s); err != nil {
		return
	}
	n := (len(s) + 1) / 2
	digits := make([]*big.Int, n)
	ops := make([]byte, 0, n-1)
	for i := 0; i < len(s); i++ {
		if i%2 == 0 {
			digits[i/2] = big.NewInt(int64(s[i] - '0'))
		} else {
			ops = append(ops, s[i])
		}
	}

	// collect candidate start and end positions based on '*' signs
	starts := []int{0}
	ends := []int{n - 1}
	for i, op := range ops {
		if op == '*' {
			starts = append(starts, i+1)
			ends = append(ends, i)
		}
	}

	maxVal := big.NewInt(0)
	for _, l := range starts {
		for _, r := range ends {
			if l > r {
				continue
			}
			val := evalWithParen(digits, ops, l, r)
			if val.Cmp(maxVal) > 0 {
				maxVal = val
			}
		}
	}
	fmt.Println(maxVal.String())
}

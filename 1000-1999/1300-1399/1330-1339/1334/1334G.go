package main

import (
	"bufio"
	"fmt"
	"math/big"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	p := make([]int, 26)
	for i := 0; i < 26; i++ {
		fmt.Fscan(reader, &p[i])
		p[i]--
	}
	var s, t string
	fmt.Fscan(reader, &s)
	fmt.Fscan(reader, &t)

	n := len(t)
	m := len(s)
	bitsets := make([]*big.Int, 26)
	for i := 0; i < 26; i++ {
		bitsets[i] = new(big.Int)
	}
	for i := 0; i < n; i++ {
		c := int(t[i] - 'a')
		bitsets[c].SetBit(bitsets[c], i, 1)
	}

	// initial result with allowed chars for first position
	idx0 := int(s[0] - 'a')
	res := new(big.Int).Or(new(big.Int).Set(bitsets[idx0]), bitsets[p[idx0]])

	for k := 1; k < m; k++ {
		idx := int(s[k] - 'a')
		tmp := new(big.Int).Or(new(big.Int).Set(bitsets[idx]), bitsets[p[idx]])
		tmp.Rsh(tmp, uint(k))
		res.And(res, tmp)
	}

	limit := n - m + 1
	for i := 0; i < limit; i++ {
		if res.Bit(i) == 1 {
			writer.WriteByte('1')
		} else {
			writer.WriteByte('0')
		}
	}
}

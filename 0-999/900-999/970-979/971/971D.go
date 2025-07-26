package main

import (
	"bufio"
	"fmt"
	"math/big"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var aStr, bStr string
	if _, err := fmt.Fscan(reader, &aStr, &bStr); err != nil {
		return
	}
	a := new(big.Int)
	b := new(big.Int)
	if _, ok := a.SetString(aStr, 10); !ok {
		return
	}
	if _, ok := b.SetString(bStr, 10); !ok {
		return
	}
	gcd := new(big.Int).GCD(nil, nil, a, b)
	fmt.Println(gcd)
}

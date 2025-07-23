package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var a, b, p, x int64
	if _, err := fmt.Fscan(reader, &a, &b, &p, &x); err != nil {
		return
	}

	P := int(p)
	powA := make([]int64, P)
	powA[0] = 1 % p
	for i := 1; i < P; i++ {
		powA[i] = powA[i-1] * a % p
	}

	L := p * (p - 1)
	ans := int64(0)
	for e := 0; e < P-1; e++ {
		r := (b * powA[P-1-e]) % p // b * inverse(a^e)
		k := int64(int64(e)-r) % (p - 1)
		if k < 0 {
			k += p - 1
		}
		n0 := r + p*k
		if n0 <= x {
			ans += 1 + (x-n0)/L
		}
	}

	fmt.Fprintln(writer, ans)
}

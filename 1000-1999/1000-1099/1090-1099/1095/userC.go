 package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, k int64
	fmt.Fscan(in, &n, &k)
	s := bits.OnesCount64(uint64(n))
	if k > n || int64(s) > k {
		fmt.Println("NO")
		return
	}
	fmt.Println("YES")
	var cnt [31]int
	temp := n
	for m := 0; temp > 0; m++ {
		if temp&1 == 1 {
			cnt[m] = 1
		}
		temp >>= 1
	}
	splits := int(k) - s
	for i := 0; i < splits; i++ {
		found := false
		for mm := 30; mm >= 1; mm-- {
			if cnt[mm] > 0 {
				cnt[mm]--
				cnt[mm-1] += 2
				found = true
				break
			}
		}
		if !found {
			// This should not happen
			panic("Unable to split further")
		}
	}
	// Now collect the list
	var res []int64
	for m := 0; m <= 30; m++ {
		pow := int64(1) << m
		for j := 0; j < cnt[m]; j++ {
			res = append(res, pow)
		}
	}
	// Output
	for i, v := range res {
		if i > 0 {
			fmt.Print(" ")
		}
		fmt.Print(v)
	}
	fmt.Println()
}
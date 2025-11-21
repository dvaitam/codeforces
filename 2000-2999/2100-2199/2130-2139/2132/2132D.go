package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

var pow10 [19]int64
var sumPow [19]int64

func init() {
	pow10[0] = 1
	for i := 1; i < len(pow10); i++ {
		pow10[i] = pow10[i-1] * 10
	}
	sumPow[0] = 0
	for i := 1; i < len(sumPow); i++ {
		sumPow[i] = int64(i) * 45 * pow10[i-1]
	}
}

func digitLen(n int64) int {
	for i := 1; i < len(pow10); i++ {
		if n < pow10[i] {
			return i
		}
	}
	return len(pow10) - 1
}

func sumDigits(n int64) int64 {
	if n <= 0 {
		return 0
	}
	if n < 10 {
		return n * (n + 1) / 2
	}
	length := digitLen(n)
	p := pow10[length-1]
	msd := n / p
	rest := n % p

	res := msd * sumPow[length-1]
	res += (msd * (msd - 1) / 2) * p
	res += msd * (rest + 1)
	res += sumDigits(rest)
	return res
}

func sumRange(l, r int64) int64 {
	if l > r {
		return 0
	}
	return sumDigits(r) - sumDigits(l-1)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var k int64
		fmt.Fscan(in, &k)
		var ans int64
		start := int64(1)
		for length := 1; ; length++ {
			count := 9 * start
			need := int64(length)
			if need > 0 && count <= k/need {
				blockDigits := count * need
				ans += sumRange(start, start+count-1)
				k -= blockDigits
				start *= 10
			} else {
				numFull := int64(0)
				if need > 0 {
					numFull = k / need
				}
				rem := int64(0)
				if need > 0 {
					rem = k % need
				}
				ans += sumRange(start, start+numFull-1)
				if rem > 0 {
					nextNum := start + numFull
					s := strconv.FormatInt(nextNum, 10)
					for i := int64(0); i < rem && i < int64(len(s)); i++ {
						ans += int64(s[i] - '0')
					}
				}
				break
			}
		}
		fmt.Fprintln(out, ans)
	}
}

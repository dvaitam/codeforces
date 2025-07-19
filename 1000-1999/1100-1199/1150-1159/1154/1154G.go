package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	const T = 10000000
	cnt := make([]int, T+1)
	ans := int64(math.MaxInt64)
	var lid, rid int

	for i := 1; i <= n; i++ {
		var v int
		fmt.Fscan(reader, &v)
		if cnt[v] != 0 {
			// duplicate value yields lcm = v
			if int64(v) < ans {
				ans = int64(v)
				lid = cnt[v]
				rid = i
			}
		}
		cnt[v] = i
	}

	// search for minimal lcm via common divisors
	for i := 1; i <= T; i++ {
		var l, r, lid1, rid1 int
		for j := i; j <= T; j += i {
			if idx := cnt[j]; idx != 0 {
				if l == 0 {
					l = j
					lid1 = idx
				} else {
					r = j
					rid1 = idx
					break
				}
			}
		}
		if r != 0 {
			res := int64(l) * int64(r) / int64(i)
			if res < ans {
				ans = res
				lid = lid1
				rid = rid1
			}
		}
	}

	if lid > rid {
		lid, rid = rid, lid
	}
	fmt.Fprintf(writer, "%d %d", lid, rid)
}

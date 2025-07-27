package main

import (
	"bufio"
	"fmt"
	"os"
)

// query sends '? v' and reads the gcd response.
func query(out *bufio.Writer, in *bufio.Reader, v int64) int64 {
	fmt.Fprintf(out, "? %d\n", v)
	out.Flush()
	var resp int64
	fmt.Fscan(in, &resp)
	return resp
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	if _, err := fmt.Fscan(in, &T); err != nil {
		return
	}

	groups := [][]int64{
		{2, 3, 5, 7, 11, 13},
		{17, 19, 23, 29, 31, 37},
		{41, 43, 47, 53, 59, 61},
	}

	for ; T > 0; T-- {
		primes := make([]int64, 0)
		for _, g := range groups {
			val := int64(1)
			for _, p := range g {
				val *= p
			}
			res := query(out, in, val)
			for _, p := range g {
				if res%p == 0 {
					primes = append(primes, p)
				}
			}
		}

		ans := int64(1)
		for _, p := range primes {
			pow := int64(1)
			for pow*p <= 1_000_000_000 {
				pow *= p
			}
			res := query(out, in, pow)
			cnt := int64(0)
			for res%p == 0 {
				res /= p
				cnt++
			}
			ans *= cnt + 1
		}

		if ans == 1 {
			ans = 2
		} else {
			ans *= 2
		}
		fmt.Fprintf(out, "! %d\n", ans)
	}
}

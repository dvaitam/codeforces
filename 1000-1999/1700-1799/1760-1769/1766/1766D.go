package main

import (
	"bufio"
	"fmt"
	"os"
)

const maxA = 10000000

var spf []int

func sieve(n int) {
	spf = make([]int, n+1)
	primes := make([]int, 0, 664000)
	for i := 2; i <= n; i++ {
		if spf[i] == 0 {
			spf[i] = i
			primes = append(primes, i)
		}
		for _, p := range primes {
			if p > spf[i] || p*i > n {
				break
			}
			spf[p*i] = p
		}
	}
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}

	sieve(maxA)

	first := true
	for i := 0; i < n; i++ {
		var x, y int
		fmt.Fscan(in, &x, &y)
		if !first {
			fmt.Fprint(out, " ")
		}
		first = false
		if gcd(x, y) != 1 {
			fmt.Fprint(out, 0)
			continue
		}
		d := y - x
		if d == 1 {
			fmt.Fprint(out, -1)
			continue
		}
		ans := int(1<<31 - 1)
		tmp := d
		for tmp > 1 {
			p := spf[tmp]
			if p == 0 {
				p = tmp
			}
			if v := p - x%p; v < ans {
				ans = v
			}
			for tmp%p == 0 {
				tmp /= p
			}
		}
		fmt.Fprint(out, ans)
	}
	fmt.Fprintln(out)
}

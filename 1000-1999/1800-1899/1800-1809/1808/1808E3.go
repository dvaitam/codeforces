package main

import (
	"fmt"
)

func abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func power(base, exp, mod int64) int64 {
	base %= mod
	if base < 0 {
		base += mod
	}
	if base == 0 {
		if exp == 0 {
			return 1
		}
		return 0
	}
	res := int64(1)
	for exp > 0 {
		if exp%2 == 1 {
			res = (res * base) % mod
		}
		base = (base * base) % mod
		exp /= 2
	}
	return res
}

func main() {
	var n, k, m int64
	if _, err := fmt.Scan(&n, &k, &m); err != nil {
		return
	}

	if k%2 == 1 {
		d := gcd(abs(n-2), k)
		term1 := power(k, n, m)
		term2 := power(k-1, n, m)
		ans := (term1 - term2 + m) % m
		diff := (d - 1) % m
		if n%2 == 0 {
			ans = (ans - diff + m) % m
		} else {
			ans = (ans + diff) % m
		}
		fmt.Println(ans)
	} else {
		d := gcd(abs(n-2), k/2)
		half := (m + 1) / 2
		term1 := (power(k, n, m) * half) % m
		term2 := (power(k-2, n, m) * half) % m
		ans := (term1 - term2 + m) % m
		term3 := (power(2, n-1, m) * ((d - 1) % m)) % m
		if n%2 == 0 {
			ans = (ans - term3 + m) % m
		} else {
			ans = (ans + term3) % m
		}
		fmt.Println(ans)
	}
}

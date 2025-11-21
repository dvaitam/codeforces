package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

func extendGCD(a, b int64) (int64, int64, int64) {
	if b == 0 {
		if a < 0 {
			return -a, -1, 0
		}
		return a, 1, 0
	}
	g, x1, y1 := extendGCD(b, a%b)
	return g, y1, x1-(a/b)*y1
}

func modInverse(a, mod int) int {
	g, x, _ := extendGCD(int64(a), int64(mod))
	if g != 1 {
		return 0
	}
	x %= int64(mod)
	if x < 0 {
		x += int64(mod)
	}
	return int(x)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var x int
		var m int64
		fmt.Fscan(in, &x, &m)

		if m == 0 {
			fmt.Fprintln(out, 0)
			continue
		}

		B := 1 << bits.Len(uint(x))
		Bmod := B % x
		g := gcd(Bmod, x)
		xPrime := x / g
		var inv int
		if xPrime > 1 {
			inv = modInverse(Bmod/g%xPrime, xPrime)
		}

		limitR := B - 1
		if m < int64(limitR) {
			limitR = int(m)
		}
		B64 := int64(B)
		ans := int64(0)

		for r := 0; r <= limitR; r++ {
			val := x ^ r
			valMod := val % x
			if valMod%g != 0 {
				continue
			}
			tMax := (m - int64(r)) / B64
			if tMax < 0 {
				continue
			}

			var cnt int64
			includeZero := false
			if xPrime == 1 {
				cnt = tMax + 1
				if r == 0 {
					includeZero = true
				}
			} else {
				valRed := valMod / g
				target := (-valRed) % xPrime
				if target < 0 {
					target += xPrime
				}
				t0 := (target * inv) % xPrime
				if int64(t0) > tMax {
					continue
				}
				cnt = (tMax-int64(t0))/int64(xPrime) + 1
				if r == 0 && t0 == 0 {
					includeZero = true
				}
			}

			if includeZero {
				cnt--
				if cnt < 0 {
					continue
				}
			}
			ans += cnt
		}

		upper := x
		if int64(upper) > m {
			upper = int(m)
		}
		for y := 1; y <= upper; y++ {
			z := x ^ y
			if z%y == 0 && z%x != 0 {
				ans++
			}
		}

		fmt.Fprintln(out, ans)
	}
}

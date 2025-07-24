package main

import (
	"bufio"
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"sort"
	"time"
)

func mulMod(a, b, mod uint64) uint64 {
	return new(big.Int).Mod(new(big.Int).Mul(new(big.Int).SetUint64(a), new(big.Int).SetUint64(b)), new(big.Int).SetUint64(mod)).Uint64()
}

func powMod(a, d, mod uint64) uint64 {
	res := uint64(1)
	a %= mod
	for d > 0 {
		if d&1 == 1 {
			res = mulMod(res, a, mod)
		}
		a = mulMod(a, a, mod)
		d >>= 1
	}
	return res
}

func gcd(a, b uint64) uint64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func isPrime(n uint64) bool {
	if n < 2 {
		return false
	}
	small := []uint64{2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37}
	for _, p := range small {
		if n%p == 0 {
			return n == p
		}
	}
	d := n - 1
	s := 0
	for d&1 == 0 {
		d >>= 1
		s++
	}
	bases := []uint64{2, 325, 9375, 28178, 450775, 9780504, 1795265022}
	for _, a := range bases {
		if a%n == 0 {
			continue
		}
		x := powMod(a, d, n)
		if x == 1 || x == n-1 {
			continue
		}
		composite := true
		for r := 1; r < s; r++ {
			x = mulMod(x, x, n)
			if x == n-1 {
				composite = false
				break
			}
		}
		if composite {
			return false
		}
	}
	return true
}

func pollardsRho(n uint64) uint64 {
	if n%2 == 0 {
		return 2
	}
	if n%3 == 0 {
		return 3
	}
	for {
		c := uint64(rand.Int63n(int64(n-1))) + 1
		x := uint64(rand.Int63n(int64(n)))
		y := x
		d := uint64(1)
		for d == 1 {
			x = (mulMod(x, x, n) + c) % n
			y = (mulMod(y, y, n) + c) % n
			y = (mulMod(y, y, n) + c) % n
			if x > y {
				d = gcd(x-y, n)
			} else {
				d = gcd(y-x, n)
			}
			if d == n {
				break
			}
		}
		if d > 1 && d < n {
			return d
		}
	}
}

func factor(n uint64, res *[]uint64) {
	if n == 1 {
		return
	}
	if isPrime(n) {
		*res = append(*res, n)
		return
	}
	d := pollardsRho(n)
	factor(d, res)
	factor(n/d, res)
}

func divisors(n uint64) []uint64 {
	var pf []uint64
	factor(n, &pf)
	mp := make(map[uint64]int)
	for _, p := range pf {
		mp[p]++
	}
	divs := []uint64{1}
	for p, e := range mp {
		sz := len(divs)
		mul := uint64(1)
		for i := 0; i < e; i++ {
			mul *= p
			for j := 0; j < sz; j++ {
				divs = append(divs, divs[j]*mul)
			}
		}
	}
	sort.Slice(divs, func(i, j int) bool { return divs[i] < divs[j] })
	return divs
}

type kData struct {
	present map[uint64]struct{}
	mex     uint64
}

func main() {
	rand.Seed(time.Now().UnixNano())
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var q int
	fmt.Fscan(reader, &q)

	numbers := make(map[uint64]struct{})
	numbers[0] = struct{}{}
	divisorCache := make(map[uint64][]uint64)
	kMap := make(map[uint64]*kData)

	for ; q > 0; q-- {
		var op string
		fmt.Fscan(reader, &op)
		if op == "+" {
			var x uint64
			fmt.Fscan(reader, &x)
			if _, ok := numbers[x]; ok {
				continue
			}
			numbers[x] = struct{}{}
			divs, ok := divisorCache[x]
			if !ok {
				divs = divisors(x)
				divisorCache[x] = divs
			}
			for _, d := range divs {
				if kd, ok := kMap[d]; ok {
					idx := x / d
					kd.present[idx] = struct{}{}
					if idx == kd.mex {
						for {
							if _, ex := kd.present[kd.mex]; ex {
								kd.mex++
							} else {
								break
							}
						}
					}
				}
			}
		} else if op == "-" {
			var x uint64
			fmt.Fscan(reader, &x)
			delete(numbers, x)
			divs := divisorCache[x]
			for _, d := range divs {
				if kd, ok := kMap[d]; ok {
					idx := x / d
					delete(kd.present, idx)
					if idx < kd.mex {
						kd.mex = idx
					}
				}
			}
		} else if op == "?" {
			var k uint64
			fmt.Fscan(reader, &k)
			kd, ok := kMap[k]
			if !ok {
				kd = &kData{present: make(map[uint64]struct{})}
				for num := range numbers {
					if num%k == 0 {
						kd.present[num/k] = struct{}{}
					}
				}
				kd.mex = 0
				for {
					if _, ex := kd.present[kd.mex]; ex {
						kd.mex++
					} else {
						break
					}
				}
				kMap[k] = kd
			} else {
				for {
					if _, ex := kd.present[kd.mex]; ex {
						kd.mex++
					} else {
						break
					}
				}
			}
			fmt.Fprintln(writer, kd.mex*k)
		}
	}
}

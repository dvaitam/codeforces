package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

const MOD = int64(999999893)

func power(base, exp int64) int64 {
	res := int64(1)
	base %= MOD
	for exp > 0 {
		if exp%2 == 1 {
			res = (res * base) % MOD
		}
		base = (base * base) % MOD
		exp /= 2
	}
	return res
}

func extGCD(a, b int64) (int64, int64, int64) {
	if a == 0 {
		return b, 0, 1
	}
	gcd, x1, y1 := extGCD(b%a, a)
	x := y1 - (b/a)*x1
	y := x1
	return gcd, x, y
}

func inverse(n int64) int64 {
	_, x, _ := extGCD(n, MOD)
	return (x%MOD + MOD) % MOD
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanWords)
	
	if !scanner.Scan() {
		return
	}
	tStr := scanner.Text()
	t, _ := strconv.Atoi(tStr)
	
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	
	for i := 0; i < t; i++ {
		scanner.Scan()
		nStr := scanner.Text()
		n, _ := strconv.ParseInt(nStr, 10, 64)
		
		var m int64 = n / 2
		var num, den int64
		
		if n%2 == 0 {
			powM := power(2, m)
			pow2M_minus_1 := power(2, 2*m-1)
			powM_plus_1 := power(2, m+1)
			
			num = (1 - powM) % MOD
			den = (pow2M_minus_1 - powM_plus_1 + 1) % MOD
		} else {
			powM := power(2, m)
			pow2M := power(2, 2*m)
			powM_plus_1 := power(2, m+1)
			
			num = (powM - 1) % MOD
			den = (pow2M + powM_plus_1 - 1) % MOD
		}
		
		num = (num%MOD + MOD) % MOD
		den = (den%MOD + MOD) % MOD
		
		ans := (num * inverse(den)) % MOD
		fmt.Fprintln(out, ans)
	}
}
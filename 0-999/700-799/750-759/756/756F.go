package main

import (
	"bufio"
	"fmt"
	"math/big"
	"os"
)

const MOD int64 = 1000000007
const PHI int64 = MOD - 1

var (
	input string
	idx   int
)

type Number struct {
	big    *big.Int
	modM   int64
	modPh  int64
	digits int
}

type Value struct {
	val int64 // modulo MOD
	len int64 // modulo PHI
}

func modPow(base, exp int64) int64 {
	base %= MOD
	if base < 0 {
		base += MOD
	}
	res := int64(1)
	for exp > 0 {
		if exp&1 == 1 {
			res = res * base % MOD
		}
		base = base * base % MOD
		exp >>= 1
	}
	return res
}

func modInv(a int64) int64 {
	return modPow(a, MOD-2)
}

func digitsBig(x *big.Int) int {
	return len(x.String())
}

func parseNumber() Number {
	start := idx
	for idx < len(input) && input[idx] >= '0' && input[idx] <= '9' {
		idx++
	}
	str := input[start:idx]
	b := new(big.Int)
	b.SetString(str, 10)
	modM := new(big.Int).Mod(b, big.NewInt(MOD)).Int64()
	modPh := new(big.Int).Mod(b, big.NewInt(PHI)).Int64()
	return Number{big: b, modM: modM, modPh: modPh, digits: len(str)}
}

func concat(a, b Value) Value {
	pow := modPow(10, b.len)
	return Value{
		val: (a.val*pow + b.val) % MOD,
		len: (a.len + b.len) % PHI,
	}
}

func repeat(val Value, count Number) Value {
	B := modPow(10, val.len)
	kphi := count.modPh
	km := count.modM
	var sum int64
	if B == 1 {
		sum = km % MOD
	} else {
		sum = (modPow(B, kphi) - 1 + MOD) % MOD
		sum = sum * modInv(B-1) % MOD
	}
	return Value{
		val: val.val * sum % MOD,
		len: val.len * kphi % PHI,
	}
}

func rangeSameDigits(l, r *big.Int, d int) (int64, int64) {
	n := new(big.Int).Sub(r, l)
	n.Add(n, big.NewInt(1))

	nModPhi := new(big.Int).Mod(n, big.NewInt(PHI)).Int64()
	nModM := new(big.Int).Mod(n, big.NewInt(MOD)).Int64()

	B := modPow(10, int64(d))
	BpowN := modPow(B, nModPhi)
	sum1 := (BpowN - 1 + MOD) % MOD * modInv(B-1) % MOD
	BpowNm1 := modPow(B, (nModPhi-1+PHI)%PHI)
	part := (1 - nModM*BpowNm1%MOD + ((nModM-1+MOD)%MOD)*BpowN%MOD) % MOD
	invDen := modInv(B - 1)
	sum2 := B % MOD * part % MOD * invDen % MOD * invDen % MOD

	rMod := new(big.Int).Mod(r, big.NewInt(MOD)).Int64()
	val := (rMod*sum1%MOD - sum2 + MOD) % MOD
	length := int64(d) % PHI * nModPhi % PHI
	return val, length
}

func rangeConcat(lNum, rNum Number) Value {
	l := new(big.Int).Set(lNum.big)
	r := rNum.big
	one := big.NewInt(1)
	ten := big.NewInt(10)
	res := Value{}
	for {
		d := digitsBig(l)
		end := new(big.Int).Exp(ten, big.NewInt(int64(d)), nil)
		end.Sub(end, one)
		if end.Cmp(r) > 0 {
			end.Set(r)
		}
		v, ln := rangeSameDigits(l, end, d)
		res = concat(res, Value{v, ln})
		if end.Cmp(r) == 0 {
			break
		}
		l = new(big.Int).Add(end, one)
	}
	return res
}

func parseTerm() Value {
	num := parseNumber()
	if idx < len(input) && input[idx] == '-' {
		idx++
		rnum := parseNumber()
		return rangeConcat(num, rnum)
	} else if idx < len(input) && input[idx] == '(' {
		idx++
		inner := parseExpression()
		if idx < len(input) && input[idx] == ')' {
			idx++
		}
		return repeat(inner, num)
	}
	return Value{val: num.modM, len: int64(num.digits) % PHI}
}

func parseExpression() Value {
	res := parseTerm()
	for idx < len(input) && input[idx] == '+' {
		idx++
		t := parseTerm()
		res = concat(res, t)
	}
	return res
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Fscan(reader, &input)
	idx = 0
	ans := parseExpression()
	fmt.Println(ans.val % MOD)
}

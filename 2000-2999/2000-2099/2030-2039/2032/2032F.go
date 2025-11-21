package main

import (
	"bufio"
	"os"
)

const mod = 998244353

type FastScanner struct {
	r *bufio.Reader
}

func NewFastScanner() *FastScanner {
	return &FastScanner{r: bufio.NewReader(os.Stdin)}
}

func (fs *FastScanner) NextInt() int {
	sign := 1
	val := 0
	b, _ := fs.r.ReadByte()
	for (b < '0' || b > '9') && b != '-' {
		b, _ = fs.r.ReadByte()
	}
	if b == '-' {
		sign = -1
		b, _ = fs.r.ReadByte()
	}
	for b >= '0' && b <= '9' {
		val = val*10 + int(b-'0')
		b, _ = fs.r.ReadByte()
	}
	return sign * val
}

func main() {
	fs := NewFastScanner()
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	t := fs.NextInt()
	for ; t > 0; t-- {
		n := fs.NextInt()
		prefix := 0
		totalEven, totalOdd := 1, 0 // dp[0][0]=1
		dpEven, dpOdd := 0, 0

		sumEven := make(map[int]int)
		sumOdd := make(map[int]int)
		sumEven[0] = 1

		for i := 0; i < n; i++ {
			x := fs.NextInt()
			prefix ^= x

			se := sumEven[prefix]
			so := sumOdd[prefix]

			deltaEven := totalOdd - so
			if deltaEven < 0 {
				deltaEven += mod
			}
			even := se + deltaEven
			if even >= mod {
				even -= mod
			}

			deltaOdd := totalEven - se
			if deltaOdd < 0 {
				deltaOdd += mod
			}
			odd := so + deltaOdd
			if odd >= mod {
				odd -= mod
			}

			dpEven = even
			dpOdd = odd

			totalEven += dpEven
			if totalEven >= mod {
				totalEven -= mod
			}
			totalOdd += dpOdd
			if totalOdd >= mod {
				totalOdd -= mod
			}

			se = (se + dpEven)
			if se >= mod {
				se -= mod
			}
			so = (so + dpOdd)
			if so >= mod {
				so -= mod
			}
			sumEven[prefix] = se
			if so != 0 {
				sumOdd[prefix] = so
			} else {
				delete(sumOdd, prefix)
			}
		}
		out.WriteString(intToString(dpOdd))
		out.WriteByte('\n')
	}
}

func intToString(x int) string {
	if x == 0 {
		return "0"
	}
	if x < 0 {
		x = (x%mod + mod) % mod
	}
	buf := make([]byte, 0, 12)
	for x > 0 {
		buf = append(buf, byte('0'+x%10))
		x /= 10
	}
	for i, j := 0, len(buf)-1; i < j; i, j = i+1, j-1 {
		buf[i], buf[j] = buf[j], buf[i]
	}
	return string(buf)
}


package main

import (
	"bufio"
	"fmt"
	"os"
)

var (
	reader = bufio.NewReader(os.Stdin)
	writer = bufio.NewWriter(os.Stdout)
)

func readInt64() int64 {
	sign := int64(1)
	val := int64(0)
	b, _ := reader.ReadByte()
	for (b < '0' || b > '9') && b != '-' {
		b, _ = reader.ReadByte()
	}
	if b == '-' {
		sign = -1
		b, _ = reader.ReadByte()
	}
	for b >= '0' && b <= '9' {
		val = val*10 + int64(b-'0')
		b, _ = reader.ReadByte()
	}
	return val * sign
}

func max(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func adjustLow(v int64, parity int64) int64 {
	if (v & 1) != parity {
		v++
	}
	return v
}

func adjustHigh(v int64, parity int64) int64 {
	if (v & 1) != parity {
		v--
	}
	return v
}

func feasible(points [][3]int64, d int64) (bool, int64, int64, int64) {
	const INF int64 = 3_000_000_000_000_000_000
	pl, ph := int64(-INF), int64(INF)
	ql, qh := int64(-INF), int64(INF)
	rl, rh := int64(-INF), int64(INF)
	sl, sh := int64(-INF), int64(INF)
	for _, p := range points {
		x, y, z := p[0], p[1], p[2]
		s1 := x + y + z
		pl = max(pl, s1-d)
		ph = min(ph, s1+d)
		s2 := x + y - z
		ql = max(ql, s2-d)
		qh = min(qh, s2+d)
		s3 := x - y + z
		rl = max(rl, s3-d)
		rh = min(rh, s3+d)
		s4 := -x + y + z
		sl = max(sl, s4-d)
		sh = min(sh, s4+d)
	}
	if pl > ph || ql > qh || rl > rh || sl > sh {
		return false, 0, 0, 0
	}
	for parity := int64(0); parity < 2; parity++ {
		pl1 := adjustLow(pl, parity)
		ph1 := adjustHigh(ph, parity)
		ql1 := adjustLow(ql, parity)
		qh1 := adjustHigh(qh, parity)
		rl1 := adjustLow(rl, parity)
		rh1 := adjustHigh(rh, parity)
		sl1 := adjustLow(sl, parity)
		sh1 := adjustHigh(sh, parity)
		if pl1 > ph1 || ql1 > qh1 || rl1 > rh1 || sl1 > sh1 {
			continue
		}
		tmin := ql1 + rl1 + sl1
		tmax := qh1 + rh1 + sh1
		low := max(pl1, tmin)
		high := min(ph1, tmax)
		if low > high {
			continue
		}
		p := low
		q := ql1
		r := rl1
		s := sl1
		rem := p - (q + r + s)
		delta := min(rem, qh1-q)
		q += delta
		rem -= delta
		delta = min(rem, rh1-r)
		r += delta
		rem -= delta
		delta = min(rem, sh1-s)
		s += delta
		rem -= delta
		if rem != 0 {
			continue
		}
		x := (q + r) / 2
		y := (p - r) / 2
		z := (p - q) / 2
		return true, x, y, z
	}
	return false, 0, 0, 0
}

func main() {
	defer writer.Flush()
	t := readInt64()
	for ; t > 0; t-- {
		n := readInt64()
		points := make([][3]int64, n)
		for i := int64(0); i < n; i++ {
			x := readInt64()
			y := readInt64()
			z := readInt64()
			points[i] = [3]int64{x, y, z}
		}
		low := int64(0)
		high := int64(3_000_000_000_000_000_000)
		var ansX, ansY, ansZ int64
		for low <= high {
			mid := (low + high) / 2
			ok, x, y, z := feasible(points, mid)
			if ok {
				ansX, ansY, ansZ = x, y, z
				high = mid - 1
			} else {
				low = mid + 1
			}
		}
		fmt.Fprintf(writer, "%d %d %d\n", ansX, ansY, ansZ)
	}
}

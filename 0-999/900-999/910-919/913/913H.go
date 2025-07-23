package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

const mod int64 = 998244353
const scale int64 = 1000000

func modPow(a, e int64) int64 {
	res := int64(1)
	for e > 0 {
		if e&1 == 1 {
			res = res * a % mod
		}
		a = a * a % mod
		e >>= 1
	}
	return res
}

func addPoly(a, b []int64) []int64 {
	n := len(a)
	if len(b) > n {
		n = len(b)
	}
	res := make([]int64, n)
	for i := 0; i < n; i++ {
		if i < len(a) {
			res[i] = (res[i] + a[i]) % mod
		}
		if i < len(b) {
			res[i] = (res[i] + b[i]) % mod
		}
	}
	return res
}

func subPoly(a, b []int64) []int64 {
	n := len(a)
	if len(b) > n {
		n = len(b)
	}
	res := make([]int64, n)
	for i := 0; i < n; i++ {
		if i < len(a) {
			res[i] = (res[i] + a[i]) % mod
		}
		if i < len(b) {
			res[i] = (res[i] - b[i]) % mod
		}
	}
	for i := 0; i < n; i++ {
		res[i] %= mod
		if res[i] < 0 {
			res[i] += mod
		}
	}
	return res
}

var binom [35][35]int64
var inv [35]int64

func shiftPoly(p []int64, shift int64) []int64 {
	s := shift % mod
	neg := (mod - s) % mod
	n := len(p) - 1
	res := make([]int64, len(p))
	for k := 0; k <= n; k++ {
		coeff := p[k]
		pow := int64(1)
		for j := 0; j <= k; j++ {
			val := coeff * binom[k][j] % mod * pow % mod
			res[j] = (res[j] + val) % mod
			pow = pow * neg % mod
		}
	}
	return res
}

func integratePoly(p []int64) []int64 {
	res := make([]int64, len(p)+1)
	for i := 0; i < len(p); i++ {
		res[i+1] = p[i] * inv[i+1] % mod
	}
	return res
}

func evalPoly(p []int64, x int64) int64 {
	res := int64(0)
	pow := int64(1)
	for i := 0; i < len(p); i++ {
		res = (res + p[i]*pow) % mod
		pow = pow * x % mod
	}
	return res
}

func toMod(x int64) int64 {
	return (x % mod) * invScale % mod
}

var invScale int64

type segment struct {
	l, r int64
	poly []int64
}

func getPoly(segs []segment, s, prevX, total int64) []int64 {
	if s < 0 {
		return []int64{0}
	}
	if s >= prevX {
		return []int64{total}
	}
	// binary search
	lo, hi := 0, len(segs)
	for lo < hi {
		mid := (lo + hi) / 2
		if segs[mid].r <= s {
			lo = mid + 1
		} else if segs[mid].l > s {
			hi = mid
		} else {
			return segs[mid].poly
		}
	}
	return []int64{total}
}

func step(segs []segment, prevX, xInt int64, total int64) ([]segment, int64) {
	if xInt < 0 {
		return []segment{}, 0
	}
	pointsMap := map[int64]struct{}{}
	pointsMap[0] = struct{}{}
	pointsMap[xInt] = struct{}{}
	if scale < xInt {
		pointsMap[scale] = struct{}{}
	}
	for _, seg := range segs {
		vals := []int64{seg.l, seg.r, seg.l + scale, seg.r + scale}
		for _, v := range vals {
			if v > 0 && v < xInt {
				pointsMap[v] = struct{}{}
			}
		}
	}
	points := make([]int64, 0, len(pointsMap))
	for v := range pointsMap {
		points = append(points, v)
	}
	sort.Slice(points, func(i, j int) bool { return points[i] < points[j] })
	newSegs := make([]segment, 0, len(points))
	curVal := int64(0)
	for i := 0; i < len(points)-1; i++ {
		L := points[i]
		R := points[i+1]
		if L >= xInt {
			break
		}
		poly1 := getPoly(segs, L, prevX, total)
		poly2 := []int64{0}
		if R <= scale {
			poly2 = []int64{0}
		} else if L >= scale {
			base := getPoly(segs, L-scale, prevX, total)
			poly2 = shiftPoly(base, scale)
		} else {
			// interval crosses the point scale, split automatically because points contain scale
			if R > scale {
				// first part [L, scale) should be processed earlier
				base := getPoly(segs, 0, prevX, total)
				_ = base
			}
		}
		gPoly := subPoly(poly1, poly2)
		intPoly := integratePoly(gPoly)
		valAtL := evalPoly(intPoly, toMod(L))
		adj := curVal - valAtL
		adj %= mod
		if adj < 0 {
			adj += mod
		}
		intPoly[0] = (intPoly[0] + adj) % mod
		newSeg := segment{l: L, r: R, poly: intPoly}
		newSegs = append(newSegs, newSeg)
		curVal = evalPoly(intPoly, toMod(R))
	}
	// ensure last segment ends at xInt
	if len(newSegs) > 0 {
		newSegs[len(newSegs)-1].r = xInt
	}
	return newSegs, curVal
}

func parseDecimal(s string) int64 {
	if strings.IndexByte(s, '.') == -1 {
		v, _ := strconv.ParseInt(s, 10, 64)
		return v * scale
	}
	parts := strings.SplitN(s, ".", 2)
	intPart, _ := strconv.ParseInt(parts[0], 10, 64)
	frac := parts[1]
	if len(frac) > 6 {
		frac = frac[:6]
	}
	for len(frac) < 6 {
		frac += "0"
	}
	fracPart, _ := strconv.ParseInt(frac, 10, 64)
	return intPart*scale + fracPart
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	xs := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		var s string
		fmt.Fscan(in, &s)
		xs[i] = parseDecimal(s)
	}

	invScale = modPow(scale%mod, mod-2)
	for i := 1; i < len(inv); i++ {
		inv[i] = modPow(int64(i), mod-2)
	}
	for i := 0; i < len(binom); i++ {
		binom[i][0] = 1
		binom[i][i] = 1
		for j := 1; j < i; j++ {
			binom[i][j] = (binom[i-1][j-1] + binom[i-1][j]) % mod
		}
	}

	big := int64(n+1) * scale
	segs := []segment{{l: 0, r: big, poly: []int64{1}}}
	prevX := big
	total := int64(1)
	for i := 1; i <= n; i++ {
		segs, total = step(segs, prevX, xs[i], total)
		prevX = xs[i]
	}
	ans := total % mod
	if ans < 0 {
		ans += mod
	}
	fmt.Fprintln(out, ans)
}

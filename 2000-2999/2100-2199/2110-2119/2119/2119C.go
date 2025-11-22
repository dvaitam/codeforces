package main

import (
	"bufio"
	"fmt"
	"os"
)

// We need two numbers u, v in [l, r] such that u&v = 0.
// For odd n, taking all elements equal to any value works, so we choose l.
// For even n >= 4, we set first n-2 elements to u (count even so XOR 0) and last two to v.
// XOR of all elements is 0, and AND of all elements is u & v, which we force to 0.
// This satisfies AND == XOR == 0. For n == 2 this would require u == v == 0 which is outside range (l >= 1), so impossible.
// Lexicographic minimality comes from minimizing u first (prefix of n-2 positions).

// DP to check existence of two numbers in [l, r] with bitwise AND 0, and to build lexicographically smallest pair.

type pairDP struct {
	l, r uint64
	mem  [61][2][2][2][2]int8
}

func (p *pairDP) dfs(pos int, ul, ur, vl, vr uint8) bool {
	if pos < 0 {
		return true
	}
	if p.mem[pos][ul][ur][vl][vr] != 0 {
		return p.mem[pos][ul][ur][vl][vr] == 1
	}
	lb := (p.l >> pos) & 1
	rb := (p.r >> pos) & 1
	for bu := uint64(0); bu <= 1; bu++ {
		for bv := uint64(0); bv <= 1; bv++ {
			if bu == 1 && bv == 1 {
				continue
			}
			if ul == 1 && bu < lb {
				continue
			}
			if ur == 1 && bu > rb {
				continue
			}
			if vl == 1 && bv < lb {
				continue
			}
			if vr == 1 && bv > rb {
				continue
			}
			nul := ul == 1 && bu == lb
			nur := ur == 1 && bu == rb
			nvl := vl == 1 && bv == lb
			nvr := vr == 1 && bv == rb
			if p.dfs(pos-1, b2u(nul), b2u(nur), b2u(nvl), b2u(nvr)) {
				p.mem[pos][ul][ur][vl][vr] = 1
				return true
			}
		}
	}
	p.mem[pos][ul][ur][vl][vr] = -1
	return false
}

func b2u(b bool) uint8 {
	if b {
		return 1
	}
	return 0
}

// build minimal u (and some consistent v flags)
func (p *pairDP) buildU() (uint64, bool) {
	if !p.dfs(60, 1, 1, 1, 1) {
		return 0, false
	}
	var u uint64
	ul, ur, vl, vr := uint8(1), uint8(1), uint8(1), uint8(1)
	for pos := 60; pos >= 0; pos-- {
		lb := (p.l >> pos) & 1
		rb := (p.r >> pos) & 1
		chosen := uint64(2)
		var nUL, nUR, nVL, nVR uint8
		for bu := uint64(0); bu <= 1; bu++ {
			okU := true
			if ul == 1 && bu < lb {
				okU = false
			}
			if ur == 1 && bu > rb {
				okU = false
			}
			if !okU {
				continue
			}
			for bv := uint64(0); bv <= 1; bv++ {
				if bu == 1 && bv == 1 {
					continue
				}
				okV := true
				if vl == 1 && bv < lb {
					okV = false
				}
				if vr == 1 && bv > rb {
					okV = false
				}
				if !okV {
					continue
				}
				nul := ul == 1 && bu == lb
				nur := ur == 1 && bu == rb
				nvl := vl == 1 && bv == lb
				nvr := vr == 1 && bv == rb
				if p.dfs(pos-1, b2u(nul), b2u(nur), b2u(nvl), b2u(nvr)) {
					chosen = bu
					nUL, nUR, nVL, nVR = b2u(nul), b2u(nur), b2u(nvl), b2u(nvr)
					break
				}
			}
			if chosen != 2 {
				break
			}
		}
		if chosen == 2 {
			return 0, false
		}
		u |= chosen << pos
		ul, ur, vl, vr = nUL, nUR, nVL, nVR
	}
	return u, true
}

// DP for building minimal v given fixed u with constraint v&u==0
type singleDP struct {
	l, r, u uint64
	mem     [61][2][2]int8
}

func (s *singleDP) dfs(pos int, tl, tr uint8) bool {
	if pos < 0 {
		return true
	}
	if s.mem[pos][tl][tr] != 0 {
		return s.mem[pos][tl][tr] == 1
	}
	lb := (s.l >> pos) & 1
	rb := (s.r >> pos) & 1
	allow1 := ((s.u >> pos) & 1) == 0
	for bv := uint64(0); bv <= 1; bv++ {
		if bv == 1 && !allow1 {
			continue
		}
		if tl == 1 && bv < lb {
			continue
		}
		if tr == 1 && bv > rb {
			continue
		}
		ntl := tl == 1 && bv == lb
		ntr := tr == 1 && bv == rb
		if s.dfs(pos-1, b2u(ntl), b2u(ntr)) {
			s.mem[pos][tl][tr] = 1
			return true
		}
	}
	s.mem[pos][tl][tr] = -1
	return false
}

func (s *singleDP) buildV() (uint64, bool) {
	if !s.dfs(60, 1, 1) {
		return 0, false
	}
	var v uint64
	tl, tr := uint8(1), uint8(1)
	for pos := 60; pos >= 0; pos-- {
		lb := (s.l >> pos) & 1
		rb := (s.r >> pos) & 1
		allow1 := ((s.u >> pos) & 1) == 0
		chosen := uint64(2)
		var nTL, nTR uint8
		for bv := uint64(0); bv <= 1; bv++ {
			if bv == 1 && !allow1 {
				continue
			}
			if tl == 1 && bv < lb {
				continue
			}
			if tr == 1 && bv > rb {
				continue
			}
			ntl := tl == 1 && bv == lb
			ntr := tr == 1 && bv == rb
			if s.dfs(pos-1, b2u(ntl), b2u(ntr)) {
				chosen = bv
				nTL, nTR = b2u(ntl), b2u(ntr)
				break
			}
		}
		if chosen == 2 {
			return 0, false
		}
		v |= chosen << pos
		tl, tr = nTL, nTR
	}
	return v, true
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n, l, r, k uint64
		fmt.Fscan(in, &n, &l, &r, &k)

		if n%2 == 1 { // all elements can be equal to l
			fmt.Fprintln(out, l)
			continue
		}
		if n == 2 { // would need both 0, impossible with l>=1
			fmt.Fprintln(out, -1)
			continue
		}

		pdp := pairDP{l: l, r: r}
		u, ok := pdp.buildU()
		if !ok {
			fmt.Fprintln(out, -1)
			continue
		}
		sdp := singleDP{l: l, r: r, u: u}
		v, ok := sdp.buildV()
		if !ok {
			fmt.Fprintln(out, -1)
			continue
		}

		if k <= n-2 {
			fmt.Fprintln(out, u)
		} else {
			fmt.Fprintln(out, v)
		}
	}
}

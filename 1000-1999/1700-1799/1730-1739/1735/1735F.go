package main

import (
	"bufio"
	"os"
	"strconv"
)

type FastScanner struct {
	data []byte
	idx  int
	n    int
}

func NewFastScanner() *FastScanner {
	data, _ := os.ReadFile("/dev/stdin")
	return &FastScanner{data: data, n: len(data)}
}

func (fs *FastScanner) NextInt64() int64 {
	for fs.idx < fs.n {
		c := fs.data[fs.idx]
		if c >= '0' && c <= '9' {
			break
		}
		fs.idx++
	}
	var v int64
	for fs.idx < fs.n {
		c := fs.data[fs.idx]
		if c < '0' || c > '9' {
			break
		}
		v = v*10 + int64(c-'0')
		fs.idx++
	}
	return v
}

type Node struct {
	l, r     *Node
	pr       uint64
	pk, qk   int64
	rate     float64
	len      int64
	selfDrop float64
	sumLen   int64
	sumDrop  float64
}

var seed uint64 = 88172645463393265

func nextRand() uint64 {
	seed += 0x9e3779b97f4a7c15
	z := seed
	z = (z ^ (z >> 30)) * 0xbf58476d1ce4e5b9
	z = (z ^ (z >> 27)) * 0x94d049bb133111eb
	return z ^ (z >> 31)
}

func getLen(t *Node) int64 {
	if t == nil {
		return 0
	}
	return t.sumLen
}

func getDrop(t *Node) float64 {
	if t == nil {
		return 0
	}
	return t.sumDrop
}

func upd(t *Node) {
	if t == nil {
		return
	}
	t.selfDrop = t.rate * float64(t.len)
	t.sumLen = t.len + getLen(t.l) + getLen(t.r)
	t.sumDrop = t.selfDrop + getDrop(t.l) + getDrop(t.r)
}

func cmpFrac(p1, q1, p2, q2 int64) int {
	lhs := q1 * p2
	rhs := q2 * p1
	if lhs < rhs {
		return -1
	}
	if lhs > rhs {
		return 1
	}
	return 0
}

func rotR(y *Node) *Node {
	x := y.l
	y.l = x.r
	x.r = y
	upd(y)
	upd(x)
	return x
}

func rotL(x *Node) *Node {
	y := x.r
	x.r = y.l
	y.l = x
	upd(x)
	upd(y)
	return y
}

func newNode(p, q, ln int64) *Node {
	t := &Node{
		pr:   nextRand(),
		pk:   p,
		qk:   q,
		rate: float64(q) / float64(p),
		len:  ln,
	}
	upd(t)
	return t
}

func insert(t *Node, p, q, ln int64) *Node {
	if t == nil {
		return newNode(p, q, ln)
	}
	c := cmpFrac(p, q, t.pk, t.qk)
	if c == 0 {
		t.len += ln
		upd(t)
		return t
	}
	if c < 0 {
		t.l = insert(t.l, p, q, ln)
		if t.l.pr > t.pr {
			t = rotR(t)
		}
	} else {
		t.r = insert(t.r, p, q, ln)
		if t.r.pr > t.pr {
			t = rotL(t)
		}
	}
	upd(t)
	return t
}

func removePrefix(t *Node, d int64) (*Node, float64) {
	if t == nil || d == 0 {
		return t, 0
	}
	leftLen := getLen(t.l)
	if d < leftLen {
		nl, rem := removePrefix(t.l, d)
		t.l = nl
		upd(t)
		return t, rem
	}
	if d == leftLen {
		rem := getDrop(t.l)
		t.l = nil
		upd(t)
		return t, rem
	}
	rem := getDrop(t.l)
	d -= leftLen
	t.l = nil
	if d < t.len {
		rem += t.rate * float64(d)
		t.len -= d
		upd(t)
		return t, rem
	}
	if d == t.len {
		rem += t.selfDrop
		return t.r, rem
	}
	rem += t.selfDrop
	nr, rem2 := removePrefix(t.r, d-t.len)
	return nr, rem + rem2
}

const eps = 1e-9

func lengthToDrop(t *Node, y float64) float64 {
	if y <= 0 {
		return 0
	}
	res := 0.0
	for t != nil {
		leftDrop := getDrop(t.l)
		if y < leftDrop-eps {
			t = t.l
			continue
		}
		leftLen := getLen(t.l)
		if y <= leftDrop+eps {
			res += float64(leftLen)
			return res
		}
		y -= leftDrop
		res += float64(leftLen)
		if y <= t.selfDrop+eps {
			return res + y/t.rate
		}
		y -= t.selfDrop
		res += float64(t.len)
		t = t.r
	}
	return res
}

func main() {
	fs := NewFastScanner()
	out := bufio.NewWriterSize(os.Stdout, 1<<20)
	defer out.Flush()

	t := int(fs.NextInt64())
	buf := make([]byte, 0, 1<<20)

	for ; t > 0; t-- {
		n := int(fs.NextInt64())
		a := fs.NextInt64()
		b := fs.NextInt64()

		p := make([]int64, n)
		q := make([]int64, n)
		for i := 0; i < n; i++ {
			p[i] = fs.NextInt64()
		}
		for i := 0; i < n; i++ {
			q[i] = fs.NextInt64()
		}

		var root *Node
		L := a
		Y := float64(b)

		for i := 0; i < n; i++ {
			pi := p[i]
			qi := q[i]

			Y += float64(qi)
			root = insert(root, pi, qi, 2*pi)

			var d int64
			if L > pi {
				L -= pi
			} else {
				d = pi - L
				L = 0
			}
			if d > 0 {
				var rem float64
				root, rem = removePrefix(root, d)
				Y -= rem
			}

			totalDrop := getDrop(root)
			totalLen := getLen(root)

			var ans float64
			if totalDrop <= Y+eps {
				ans = float64(L + totalLen)
			} else {
				yy := Y
				if yy < 0 {
					yy = 0
				}
				ans = float64(L) + lengthToDrop(root, yy)
			}

			buf = strconv.AppendFloat(buf, ans, 'f', 10, 64)
			if i+1 == n {
				buf = append(buf, '\n')
			} else {
				buf = append(buf, ' ')
			}
		}
	}

	out.Write(buf)
}
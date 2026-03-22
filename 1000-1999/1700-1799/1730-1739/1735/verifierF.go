package main

import (
	"bufio"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// ---------- embedded reference solver for 1735F ----------

type treapNode struct {
	l, r     *treapNode
	pr       uint64
	pk, qk   int64
	rate     float64
	length   int64
	selfDrop float64
	sumLen   int64
	sumDrop  float64
}

var treapSeed uint64 = 88172645463393265

func treapNextRand() uint64 {
	treapSeed += 0x9e3779b97f4a7c15
	z := treapSeed
	z = (z ^ (z >> 30)) * 0xbf58476d1ce4e5b9
	z = (z ^ (z >> 27)) * 0x94d049bb133111eb
	return z ^ (z >> 31)
}

func treapGetLen(t *treapNode) int64 {
	if t == nil {
		return 0
	}
	return t.sumLen
}

func treapGetDrop(t *treapNode) float64 {
	if t == nil {
		return 0
	}
	return t.sumDrop
}

func treapUpd(t *treapNode) {
	if t == nil {
		return
	}
	t.selfDrop = t.rate * float64(t.length)
	t.sumLen = t.length + treapGetLen(t.l) + treapGetLen(t.r)
	t.sumDrop = t.selfDrop + treapGetDrop(t.l) + treapGetDrop(t.r)
}

func treapCmpFrac(p1, q1, p2, q2 int64) int {
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

func treapRotR(y *treapNode) *treapNode {
	x := y.l
	y.l = x.r
	x.r = y
	treapUpd(y)
	treapUpd(x)
	return x
}

func treapRotL(x *treapNode) *treapNode {
	y := x.r
	x.r = y.l
	y.l = x
	treapUpd(x)
	treapUpd(y)
	return y
}

func treapNewNode(p, q, ln int64) *treapNode {
	t := &treapNode{
		pr:     treapNextRand(),
		pk:     p,
		qk:     q,
		rate:   float64(q) / float64(p),
		length: ln,
	}
	treapUpd(t)
	return t
}

func treapInsert(t *treapNode, p, q, ln int64) *treapNode {
	if t == nil {
		return treapNewNode(p, q, ln)
	}
	c := treapCmpFrac(p, q, t.pk, t.qk)
	if c == 0 {
		t.length += ln
		treapUpd(t)
		return t
	}
	if c < 0 {
		t.l = treapInsert(t.l, p, q, ln)
		if t.l.pr > t.pr {
			t = treapRotR(t)
		}
	} else {
		t.r = treapInsert(t.r, p, q, ln)
		if t.r.pr > t.pr {
			t = treapRotL(t)
		}
	}
	treapUpd(t)
	return t
}

func treapRemovePrefix(t *treapNode, d int64) (*treapNode, float64) {
	if t == nil || d == 0 {
		return t, 0
	}
	leftLen := treapGetLen(t.l)
	if d < leftLen {
		nl, rem := treapRemovePrefix(t.l, d)
		t.l = nl
		treapUpd(t)
		return t, rem
	}
	if d == leftLen {
		rem := treapGetDrop(t.l)
		t.l = nil
		treapUpd(t)
		return t, rem
	}
	rem := treapGetDrop(t.l)
	d -= leftLen
	t.l = nil
	if d < t.length {
		rem += t.rate * float64(d)
		t.length -= d
		treapUpd(t)
		return t, rem
	}
	if d == t.length {
		rem += t.selfDrop
		return t.r, rem
	}
	rem += t.selfDrop
	nr, rem2 := treapRemovePrefix(t.r, d-t.length)
	return nr, rem + rem2
}

const treapEps = 1e-9

func treapLengthToDrop(t *treapNode, y float64) float64 {
	if y <= 0 {
		return 0
	}
	res := 0.0
	for t != nil {
		leftDrop := treapGetDrop(t.l)
		if y < leftDrop-treapEps {
			t = t.l
			continue
		}
		leftLen := treapGetLen(t.l)
		if y <= leftDrop+treapEps {
			res += float64(leftLen)
			return res
		}
		y -= leftDrop
		res += float64(leftLen)
		if y <= t.selfDrop+treapEps {
			return res + y/t.rate
		}
		y -= t.selfDrop
		res += float64(t.length)
		t = t.r
	}
	return res
}

func refSolve(input string) string {
	// Reset seed for determinism per call
	treapSeed = 88172645463393265

	data := []byte(input)
	idx := 0
	nextInt64 := func() int64 {
		for idx < len(data) {
			c := data[idx]
			if c >= '0' && c <= '9' {
				break
			}
			idx++
		}
		var v int64
		for idx < len(data) {
			c := data[idx]
			if c < '0' || c > '9' {
				break
			}
			v = v*10 + int64(c-'0')
			idx++
		}
		return v
	}

	var sb strings.Builder

	t := int(nextInt64())
	for ; t > 0; t-- {
		n := int(nextInt64())
		a := nextInt64()
		b := nextInt64()

		p := make([]int64, n)
		q := make([]int64, n)
		for i := 0; i < n; i++ {
			p[i] = nextInt64()
		}
		for i := 0; i < n; i++ {
			q[i] = nextInt64()
		}

		var root *treapNode
		L := a
		Y := float64(b)

		for i := 0; i < n; i++ {
			pi := p[i]
			qi := q[i]

			Y += float64(qi)
			root = treapInsert(root, pi, qi, 2*pi)

			var d int64
			if L > pi {
				L -= pi
			} else {
				d = pi - L
				L = 0
			}
			if d > 0 {
				var rem float64
				root, rem = treapRemovePrefix(root, d)
				Y -= rem
			}

			totalDrop := treapGetDrop(root)
			totalLen := treapGetLen(root)

			var ans float64
			if totalDrop <= Y+treapEps {
				ans = float64(L + totalLen)
			} else {
				yy := Y
				if yy < 0 {
					yy = 0
				}
				ans = float64(L) + treapLengthToDrop(root, yy)
			}

			sb.WriteString(strconv.FormatFloat(ans, 'f', 10, 64))
			if i+1 == n {
				sb.WriteByte('\n')
			} else {
				sb.WriteByte(' ')
			}
		}
	}

	return strings.TrimSpace(sb.String())
}

// ---------- verifier harness ----------

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out strings.Builder
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func genCase(rng *rand.Rand) string {
	n := rng.Intn(5) + 1
	a := rng.Int63n(1000000000)
	b := rng.Int63n(1000000000)
	p := make([]int64, n)
	q := make([]int64, n)
	for i := 0; i < n; i++ {
		p[i] = rng.Int63n(1000000000) + 1
		q[i] = rng.Int63n(1000000000) + 1
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "1\n%d %d %d\n", n, a, b)
	for i, v := range p {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')
	for i, v := range q {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')
	return sb.String()
}

func genCases() []string {
	rng := rand.New(rand.NewSource(1735))
	cases := make([]string, 0, 100)
	cases = append(cases, "1\n1 0 0\n1\n1\n")
	for len(cases) < 100 {
		cases = append(cases, genCase(rng))
	}
	return cases
}

func floatsMatch(exp, got string) bool {
	expTokens := strings.Fields(exp)
	gotTokens := strings.Fields(got)
	if len(expTokens) != len(gotTokens) {
		return false
	}
	for i := range expTokens {
		a, errA := strconv.ParseFloat(expTokens[i], 64)
		b, errB := strconv.ParseFloat(gotTokens[i], 64)
		if errA != nil || errB != nil {
			if expTokens[i] != gotTokens[i] {
				return false
			}
			continue
		}
		diff := math.Abs(a - b)
		denom := math.Max(1.0, math.Abs(b))
		if diff/denom > 1e-6 {
			return false
		}
	}
	return true
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	_ = bufio.NewReader(os.Stdin)

	cases := genCases()
	for i, tc := range cases {
		exp := refSolve(tc)
		got, err := runBinary(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on case %d: %v\ninput:\n%s", i+1, err, tc)
			os.Exit(1)
		}
		if !floatsMatch(exp, got) {
			fmt.Fprintf(os.Stderr, "wrong answer on case %d\ninput:\n%sexpected:\n%s\ngot:\n%s\n", i+1, tc, exp, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed!\n", len(cases))
}

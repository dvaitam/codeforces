package main

import (
	"bufio"
	"fmt"
	"math"
	"math/rand"
	"os"
	"time"
)

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

func gcd3(a, b, c int64) int64 {
	return gcd(gcd(a, b), c)
}

type triple struct {
	x, y, z int64
}

func cross(a, b triple) triple {
	return triple{
		a.y*b.z - a.z*b.y,
		a.z*b.x - a.x*b.z,
		a.x*b.y - a.y*b.x,
	}
}

func dot(a, b triple) int64 {
	return a.x*b.x + a.y*b.y + a.z*b.z
}

func normalize(v triple) triple {
	g := gcd3(abs(v.x), abs(v.y), abs(v.z))
	if g == 0 {
		return triple{0, 0, 0}
	}
	return triple{v.x / g, v.y / g, v.z / g}
}

func abs(a int64) int64 {
	if a < 0 {
		return -a
	}
	return a
}

// treap implementation for multiset of int64 keys

type tnode struct {
	key   int64
	prio  uint32
	cnt   int
	left  *tnode
	right *tnode
}

func tsize(n *tnode) int {
	if n == nil {
		return 0
	}
	return n.cnt + tsize(n.left) + tsize(n.right)
}

func rotateRight(n *tnode) *tnode {
	l := n.left
	n.left = l.right
	l.right = n
	return l
}

func rotateLeft(n *tnode) *tnode {
	r := n.right
	n.right = r.left
	r.left = n
	return r
}

func tinsert(n *tnode, key int64) *tnode {
	if n == nil {
		return &tnode{key: key, prio: rand.Uint32(), cnt: 1}
	}
	if key == n.key {
		n.cnt++
	} else if key < n.key {
		n.left = tinsert(n.left, key)
		if n.left.prio < n.prio {
			n = rotateRight(n)
		}
	} else {
		n.right = tinsert(n.right, key)
		if n.right.prio < n.prio {
			n = rotateLeft(n)
		}
	}
	return n
}

func tdelete(n *tnode, key int64) *tnode {
	if n == nil {
		return nil
	}
	if key == n.key {
		if n.cnt > 1 {
			n.cnt--
			return n
		}
		if n.left == nil {
			return n.right
		}
		if n.right == nil {
			return n.left
		}
		if n.left.prio < n.right.prio {
			n = rotateRight(n)
			n.right = tdelete(n.right, key)
		} else {
			n = rotateLeft(n)
			n.left = tdelete(n.left, key)
		}
	} else if key < n.key {
		n.left = tdelete(n.left, key)
	} else {
		n.right = tdelete(n.right, key)
	}
	return n
}

func tpredecessor(n *tnode, key int64) *tnode {
	var res *tnode
	for n != nil {
		if key <= n.key {
			n = n.left
		} else {
			res = n
			n = n.right
		}
	}
	return res
}

func tsuccessor(n *tnode, key int64) *tnode {
	var res *tnode
	for n != nil {
		if key >= n.key {
			n = n.right
		} else {
			res = n
			n = n.left
		}
	}
	return res
}

func tmin(n *tnode) *tnode {
	for n != nil && n.left != nil {
		n = n.left
	}
	return n
}

func tmax(n *tnode) *tnode {
	for n != nil && n.right != nil {
		n = n.right
	}
	return n
}

// orientation data

type orientData struct {
	count    int
	posCount int
	angle    int64
}

var (
	target    triple
	base1     triple
	base2     triple
	scale     int64 = 1000000000 // 1e9
	fullAngle int64
	halfAngle int64

	angleRoot *tnode
	gapRoot   *tnode
	angleCnt  map[int64]int
	numAngles int

	gapCnt map[int64]int

	orientMap map[triple]*orientData
	pairExist map[triple]bool
	pairCount int

	positiveTotal int
	parallelCnt   int
)

func initBase(sf, pf, gf int64) {
	target = triple{sf, pf, gf}
	base1 = cross(target, triple{1, 0, 0})
	if base1.x == 0 && base1.y == 0 && base1.z == 0 {
		base1 = cross(target, triple{0, 1, 0})
	}
	base2 = cross(target, base1)
	fullAngle = int64(math.Round(2 * math.Pi * float64(scale)))
	halfAngle = int64(math.Round(math.Pi * float64(scale)))
	angleCnt = make(map[int64]int)
	gapCnt = make(map[int64]int)
	orientMap = make(map[triple]*orientData)
	pairExist = make(map[triple]bool)
}

func calcAngle(v triple) int64 {
	x := dot(v, base1)
	y := dot(v, base2)
	ang := math.Atan2(float64(y), float64(x))
	if ang < 0 {
		ang += 2 * math.Pi
	}
	return int64(math.Round(ang * float64(scale)))
}

func diffAngle(a, b int64) int64 {
	if b < a {
		b += fullAngle
	}
	return b - a
}

func insertGap(g int64) {
	gapRoot = tinsert(gapRoot, g)
	gapCnt[g]++
}

func removeGap(g int64) {
	if gapCnt[g] == 0 {
		return
	}
	gapRoot = tdelete(gapRoot, g)
	gapCnt[g]--
}

func maxGap() int64 {
	if gapRoot == nil {
		return fullAngle
	}
	return tmax(gapRoot).key
}

func insertAngle(a int64) {
	angleCnt[a]++
	if angleCnt[a] > 1 {
		return
	}
	if numAngles == 0 {
		angleRoot = tinsert(angleRoot, a)
		numAngles = 1
		insertGap(fullAngle)
		return
	}
	p := tpredecessor(angleRoot, a)
	if p == nil {
		p = tmax(angleRoot)
	}
	s := tsuccessor(angleRoot, a)
	if s == nil {
		s = tmin(angleRoot)
	}
	if numAngles == 1 {
		removeGap(fullAngle)
	} else {
		removeGap(diffAngle(p.key, s.key))
	}
	insertGap(diffAngle(p.key, a))
	insertGap(diffAngle(a, s.key))
	angleRoot = tinsert(angleRoot, a)
	numAngles++
}

func removeAngle(a int64) {
	if angleCnt[a] == 0 {
		return
	}
	angleCnt[a]--
	if angleCnt[a] > 0 {
		return
	}
	if numAngles == 1 {
		angleRoot = tdelete(angleRoot, a)
		numAngles = 0
		removeGap(fullAngle)
		return
	}
	p := tpredecessor(angleRoot, a)
	if p == nil {
		p = tmax(angleRoot)
	}
	s := tsuccessor(angleRoot, a)
	if s == nil {
		s = tmin(angleRoot)
	}
	removeGap(diffAngle(p.key, a))
	removeGap(diffAngle(a, s.key))
	if numAngles == 2 {
		insertGap(fullAngle)
	} else {
		insertGap(diffAngle(p.key, s.key))
	}
	angleRoot = tdelete(angleRoot, a)
	numAngles--
}

func less(a, b triple) bool {
	if a.x != b.x {
		return a.x < b.x
	}
	if a.y != b.y {
		return a.y < b.y
	}
	return a.z < b.z
}

func canonical(w triple) triple {
	opp := triple{-w.x, -w.y, -w.z}
	if less(opp, w) {
		return opp
	}
	return w
}

func pairCond(w triple) bool {
	ow, ok1 := orientMap[w]
	oppKey := triple{-w.x, -w.y, -w.z}
	ow2, ok2 := orientMap[oppKey]
	if !ok1 || ow.count == 0 || !ok2 || ow2.count == 0 {
		return false
	}
	return ow.posCount > 0 || ow2.posCount > 0
}

func updatePair(w triple) {
	key := canonical(w)
	old := pairExist[key]
	new := pairCond(key)
	if old != new {
		if new {
			pairCount++
		} else {
			pairCount--
		}
		pairExist[key] = new
	}
}

func addBottle(id int, b triple) {
	c := cross(b, target)
	if c.x == 0 && c.y == 0 && c.z == 0 {
		parallelCnt++
		positiveTotal++
		bottles[id] = bottleInfo{crossZero: true}
		return
	}
	ang := calcAngle(c)
	w := normalize(c)
	dotPos := dot(b, target) > 0
	data, ok := orientMap[w]
	if !ok {
		data = &orientData{angle: ang}
		orientMap[w] = data
		insertAngle(ang)
	}
	updatePair(w) // old state
	data.count++
	if dotPos {
		data.posCount++
		positiveTotal++
	}
	updatePair(w) // new state
	bottles[id] = bottleInfo{crossZero: false, orient: w, angle: ang, dotPos: dotPos}
}

func removeBottle(id int) {
	info := bottles[id]
	if info.crossZero {
		parallelCnt--
		positiveTotal--
		return
	}
	w := info.orient
	data := orientMap[w]
	if data == nil {
		return
	}
	updatePair(w)
	data.count--
	if info.dotPos {
		data.posCount--
		positiveTotal--
	}
	if data.count == 0 {
		removeAngle(info.angle)
		delete(orientMap, w)
	}
	updatePair(w)
}

type bottleInfo struct {
	crossZero bool
	orient    triple
	angle     int64
	dotPos    bool
}

var bottles []bottleInfo

func main() {
	rand.Seed(time.Now().UnixNano())
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var sf, pf, gf int64
	if _, err := fmt.Fscan(reader, &sf, &pf, &gf); err != nil {
		return
	}
	initBase(sf, pf, gf)

	var n int
	fmt.Fscan(reader, &n)
	bottles = make([]bottleInfo, n+1)
	for i := 1; i <= n; i++ {
		var typ string
		fmt.Fscan(reader, &typ)
		if typ == "A" {
			var s, p, g int64
			fmt.Fscan(reader, &s, &p, &g)
			addBottle(i, triple{s, p, g})
		} else {
			var r int
			fmt.Fscan(reader, &r)
			removeBottle(r)
		}
		var ans int
		if parallelCnt > 0 {
			ans = 1
		} else if pairCount > 0 {
			ans = 2
		} else if positiveTotal > 0 && maxGap() <= halfAngle {
			ans = 3
		} else {
			ans = 0
		}
		fmt.Fprintln(writer, ans)
	}
}

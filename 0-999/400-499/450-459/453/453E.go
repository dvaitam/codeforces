package main

import (
   "bufio"
   "fmt"
   "io"
   "os"
)

const INF = int64(4e18)

// Node represents a segment tree node
type Node struct {
   // static
   staticSr    int64 // sum of r_i for reset state
   sumM        int64 // sum of m_i
   staticUmin  int64 // min ceil(m_i/r_i) for reset
   staticUmax  int64 // max ceil(m_i/r_i) for reset
   // dynamic
   sm   int64 // current sum of mana
   sr   int64 // sum of r_i for regenable leaves
   mnD  int64 // min remaining time to full among regenable leaves
   mxD  int64 // max remaining time to full among regenable leaves
   // lazy
   lazyAdd   int64
   lazyReset bool
}

var (
   n, m    int
   poniesS []int64
   poniesM []int64
   poniesR []int64
   tree    []Node
   lastT   int64
)

func ceilDiv(a, b int64) int64 {
   if b == 0 {
       return INF
   }
   if a <= 0 {
       return 0
   }
   d := a / b
   if a % b != 0 {
       d++
   }
   return d
}

func build(idx, l, r int) {
   if idx >= len(tree) {
       return
   }
   if l == r {
       s := poniesS[l]
       mcap := poniesM[l]
       rrate := poniesR[l]
       tree[idx].staticSr = rrate
       tree[idx].sumM = mcap
       // reset state u_i = ceil(m_i / r_i)
       ui := ceilDiv(mcap, rrate)
       tree[idx].staticUmin = ui
       tree[idx].staticUmax = ui
       // dynamic state initial d_i = ceil((m_i - s_i)/r_i)
       di := ceilDiv(mcap - s, rrate)
       tree[idx].mnD = di
       tree[idx].mxD = di
       tree[idx].sr = rrate
       tree[idx].sm = s
       tree[idx].lazyAdd = 0
       tree[idx].lazyReset = false
       return
   }
   mid := (l + r) >> 1
   left, right := idx<<1, idx<<1|1
   build(left, l, mid)
   build(right, mid+1, r)
   pull(idx)
}

// pull merges children into node
func pull(idx int) {
   left, right := idx<<1, idx<<1|1
   a, b := &tree[left], &tree[right]
   cur := &tree[idx]
   // static
   cur.staticSr = a.staticSr + b.staticSr
   cur.sumM = a.sumM + b.sumM
   // static u bounds
   if a.staticUmin < b.staticUmin {
       cur.staticUmin = a.staticUmin
   } else {
       cur.staticUmin = b.staticUmin
   }
   if a.staticUmax > b.staticUmax {
       cur.staticUmax = a.staticUmax
   } else {
       cur.staticUmax = b.staticUmax
   }
   // dynamic
   cur.sm = a.sm + b.sm
   cur.sr = a.sr + b.sr
   if a.mnD < b.mnD {
       cur.mnD = a.mnD
   } else {
       cur.mnD = b.mnD
   }
   if a.mxD > b.mxD {
       cur.mxD = a.mxD
   } else {
       cur.mxD = b.mxD
   }
   cur.lazyAdd = 0
   cur.lazyReset = false
}

// applyApply applies time advance dt to node idx
func applyTime(idx int, dt int64) {
   if dt <= 0 {
       return
   }
   node := &tree[idx]
   if node.sr == 0 {
       return
   }
   // no one full in dt
   if dt < node.mnD {
       node.lazyAdd += dt
       node.sm += node.sr * dt
       node.mnD -= dt
       node.mxD -= dt
       return
   }
   // all full
   if dt >= node.mxD {
       node.lazyAdd = 0
       node.lazyReset = false
       node.sm = node.sumM
       node.sr = 0
       node.mnD = INF
       node.mxD = INF
       return
   }
   // partial
   push(idx)
   applyTime(idx<<1, dt)
   applyTime(idx<<1|1, dt)
   pull(idx)
}

// push down lazy tags
func push(idx int) {
   node := &tree[idx]
   left, right := idx<<1, idx<<1|1
   if node.lazyReset {
       // apply reset to children
       for _, c := range []int{left, right} {
           child := &tree[c]
           child.lazyReset = true
           child.lazyAdd = 0
           child.sm = 0
           child.sr = child.staticSr
           child.mnD = child.staticUmin
           child.mxD = child.staticUmax
       }
       node.lazyReset = false
   }
   if node.lazyAdd != 0 {
       dt := node.lazyAdd
       for _, c := range []int{left, right} {
           applyTime(c, dt)
       }
       node.lazyAdd = 0
   }
}

// rangeSum returns sum over [ql,qr]
func rangeSum(idx, l, r, ql, qr int) int64 {
   if ql > r || qr < l {
       return 0
   }
   if ql <= l && r <= qr {
       return tree[idx].sm
   }
   push(idx)
   mid := (l + r) >> 1
   return rangeSum(idx<<1, l, mid, ql, qr) + rangeSum(idx<<1|1, mid+1, r, ql, qr)
}

// rangeReset resets state over [ql,qr] at current time
func rangeReset(idx, l, r, ql, qr int) {
   if ql > r || qr < l {
       return
   }
   if ql <= l && r <= qr {
       node := &tree[idx]
       node.lazyReset = true
       node.lazyAdd = 0
       node.sm = 0
       node.sr = node.staticSr
       node.mnD = node.staticUmin
       node.mxD = node.staticUmax
       return
   }
   push(idx)
   mid := (l + r) >> 1
   rangeReset(idx<<1, l, mid, ql, qr)
   rangeReset(idx<<1|1, mid+1, r, ql, qr)
   pull(idx)
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   // input n
   if _, err := fmt.Fscan(in, &n); err == io.EOF {
       return
   }
   poniesS = make([]int64, n)
   poniesM = make([]int64, n)
   poniesR = make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &poniesS[i], &poniesM[i], &poniesR[i])
   }
   tree = make([]Node, 4*n+10)
   build(1, 0, n-1)
   lastT = 0
   fmt.Fscan(in, &m)
   for i := 0; i < m; i++ {
       var t, l, r int64
       fmt.Fscan(in, &t, &l, &r)
       // zero-based
       li := int(l - 1)
       ri := int(r - 1)
       dt := t - lastT
       if dt > 0 {
           applyTime(1, dt)
           lastT = t
       }
       res := rangeSum(1, 0, n-1, li, ri)
       fmt.Fprintln(out, res)
       rangeReset(1, 0, n-1, li, ri)
   }
}

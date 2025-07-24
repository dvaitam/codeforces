package main

import (
   "bufio"
   "fmt"
   "os"
)

type Node struct {
   len        int
   lVal, rVal int64
   preInc, preDec  int
   sufInc, sufDec  int
   maxHill    int
   add        int64
}

var (
   n int
   a []int64
   tree []Node
   reader *bufio.Reader
   writer *bufio.Writer
)

func max(a, b int) int {
   if a > b {
       return a
   }
   return b
}

func newNode(v int64) Node {
   return Node{len:1, lVal:v, rVal:v,
       preInc:1, preDec:1, sufInc:1, sufDec:1,
       maxHill:1, add:0}
}

func merge(left, right Node) Node {
   if left.len == 0 {
       return right
   }
   if right.len == 0 {
       return left
   }
   var res Node
   res.len = left.len + right.len
   res.lVal = left.lVal
   res.rVal = right.rVal
   // prefix increasing
   res.preInc = left.preInc
   if left.preInc == left.len && left.rVal < right.lVal {
       res.preInc = left.len + right.preInc
   }
   // prefix decreasing
   res.preDec = left.preDec
   if left.preDec == left.len && left.rVal > right.lVal {
       res.preDec = left.len + right.preDec
   }
   // suffix increasing
   res.sufInc = right.sufInc
   if right.sufInc == right.len && left.rVal < right.lVal {
       res.sufInc = right.len + left.sufInc
   }
   // suffix decreasing
   res.sufDec = right.sufDec
   if right.sufDec == right.len && left.rVal > right.lVal {
       res.sufDec = right.len + left.sufDec
   }
   // max hill
   res.maxHill = max(left.maxHill, right.maxHill)
   // pure increasing merge
   if left.rVal < right.lVal {
       res.maxHill = max(res.maxHill, left.sufInc + right.preInc)
   }
   // pure decreasing merge
   if left.rVal > right.lVal {
       res.maxHill = max(res.maxHill, left.sufDec + right.preDec)
       // hill with peak at boundary
       res.maxHill = max(res.maxHill, left.sufInc + right.preDec)
   }
   return res
}

func build(idx, l, r int) {
   if l == r {
       tree[idx] = newNode(a[l])
       return
   }
   mid := (l + r) >> 1
   lc, rc := idx<<1, idx<<1|1
   build(lc, l, mid)
   build(rc, mid+1, r)
   tree[idx] = merge(tree[lc], tree[rc])
}

func applyAdd(idx int, v int64) {
   node := &tree[idx]
   node.lVal += v
   node.rVal += v
   node.add += v
}

func push(idx int) {
   v := tree[idx].add
   if v != 0 {
       applyAdd(idx<<1, v)
       applyAdd(idx<<1|1, v)
       tree[idx].add = 0
   }
}

func update(idx, l, r, ql, qr int, v int64) {
   if ql <= l && r <= qr {
       applyAdd(idx, v)
       return
   }
   push(idx)
   mid := (l + r) >> 1
   if ql <= mid {
       update(idx<<1, l, mid, ql, qr, v)
   }
   if qr > mid {
       update(idx<<1|1, mid+1, r, ql, qr, v)
   }
   tree[idx] = merge(tree[idx<<1], tree[idx<<1|1])
}

func main() {
   reader = bufio.NewReader(os.Stdin)
   writer = bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   var m int
   fmt.Fscan(reader, &n)
   a = make([]int64, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   fmt.Fscan(reader, &m)
   tree = make([]Node, 4*(n+1))
   build(1, 1, n)
   for i := 0; i < m; i++ {
       var l, r int
       var d int64
       fmt.Fscan(reader, &l, &r, &d)
       update(1, 1, n, l, r, d)
       fmt.Fprintln(writer, tree[1].maxHill)
   }
}

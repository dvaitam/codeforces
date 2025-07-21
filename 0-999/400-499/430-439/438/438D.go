package main

import (
   "bufio"
   "fmt"
   "os"
)

// segment tree for range sum and range max, supporting range modulo update and point assign
var (
   n, m   int
   a       []int64
   treeSum []int64
   treeMax []int64
)

func build(node, l, r int) {
   if l == r {
       treeSum[node] = a[l]
       treeMax[node] = a[l]
       return
   }
   mid := (l + r) >> 1
   lc := node << 1
   rc := lc | 1
   build(lc, l, mid)
   build(rc, mid+1, r)
   treeSum[node] = treeSum[lc] + treeSum[rc]
   if treeMax[lc] > treeMax[rc] {
       treeMax[node] = treeMax[lc]
   } else {
       treeMax[node] = treeMax[rc]
   }
}

// apply modulo x on [ql, qr]
func updateMod(node, l, r, ql, qr int, x int64) {
   if l > qr || r < ql || treeMax[node] < x {
       return
   }
   if l == r {
       treeSum[node] %= x
       treeMax[node] = treeSum[node]
       return
   }
   mid := (l + r) >> 1
   lc := node << 1
   rc := lc | 1
   updateMod(lc, l, mid, ql, qr, x)
   updateMod(rc, mid+1, r, ql, qr, x)
   treeSum[node] = treeSum[lc] + treeSum[rc]
   if treeMax[lc] > treeMax[rc] {
       treeMax[node] = treeMax[lc]
   } else {
       treeMax[node] = treeMax[rc]
   }
}

// point assign a[idx] = x
func updateSet(node, l, r, idx int, x int64) {
   if l == r {
       treeSum[node] = x
       treeMax[node] = x
       return
   }
   mid := (l + r) >> 1
   lc := node << 1
   rc := lc | 1
   if idx <= mid {
       updateSet(lc, l, mid, idx, x)
   } else {
       updateSet(rc, mid+1, r, idx, x)
   }
   treeSum[node] = treeSum[lc] + treeSum[rc]
   if treeMax[lc] > treeMax[rc] {
       treeMax[node] = treeMax[lc]
   } else {
       treeMax[node] = treeMax[rc]
   }
}

// query sum on [ql, qr]
func querySum(node, l, r, ql, qr int) int64 {
   if l > qr || r < ql {
       return 0
   }
   if ql <= l && r <= qr {
       return treeSum[node]
   }
   mid := (l + r) >> 1
   return querySum(node<<1, l, mid, ql, qr) + querySum(node<<1|1, mid+1, r, ql, qr)
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fscan(reader, &n, &m)
   a = make([]int64, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   size := 4 * n
   treeSum = make([]int64, size)
   treeMax = make([]int64, size)
   build(1, 1, n)
   for i := 0; i < m; i++ {
       var typ int
       fmt.Fscan(reader, &typ)
       switch typ {
       case 1:
           var l, r int
           fmt.Fscan(reader, &l, &r)
           res := querySum(1, 1, n, l, r)
           fmt.Fprintln(writer, res)
       case 2:
           var l, r int
           var x int64
           fmt.Fscan(reader, &l, &r, &x)
           updateMod(1, 1, n, l, r, x)
       case 3:
           var k int
           var x int64
           fmt.Fscan(reader, &k, &x)
           updateSet(1, 1, n, k, x)
       }
   }
}

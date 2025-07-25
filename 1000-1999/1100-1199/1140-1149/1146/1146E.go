package main

import (
   "bufio"
   "fmt"
   "os"
)

const (
   Offset = 100000
   N      = 200001
)

type Node struct {
   mn, mx int
   lazy   bool
}

var tree []Node

func min(a, b int) int {
   if a < b {
       return a
   }
   return b
}

func max(a, b int) int {
   if a > b {
       return a
   }
   return b
}

// apply invert to node
func applyInv(idx int) {
   node := &tree[idx]
   node.lazy = !node.lazy
   mn, mx := node.mn, node.mx
   node.mn = -mx
   node.mx = -mn
}

func push(idx int) {
   if tree[idx].lazy {
       applyInv(idx*2)
       applyInv(idx*2+1)
       tree[idx].lazy = false
   }
}

func build(idx, l, r int) {
   if r-l == 1 {
       v := l - Offset
       tree[idx].mn = v
       tree[idx].mx = v
       return
   }
   m := (l + r) >> 1
   build(idx*2, l, m)
   build(idx*2+1, m, r)
   tree[idx].mn = min(tree[idx*2].mn, tree[idx*2+1].mn)
   tree[idx].mx = max(tree[idx*2].mx, tree[idx*2+1].mx)
}

// flip all f(i) > x
func updateGreater(idx, l, r, x int) {
   if tree[idx].mx <= x {
       return
   }
   if tree[idx].mn > x {
       applyInv(idx)
       return
   }
   push(idx)
   m := (l + r) >> 1
   updateGreater(idx*2, l, m, x)
   updateGreater(idx*2+1, m, r, x)
   tree[idx].mn = min(tree[idx*2].mn, tree[idx*2+1].mn)
   tree[idx].mx = max(tree[idx*2].mx, tree[idx*2+1].mx)
}

// flip all f(i) < x
func updateLess(idx, l, r, x int) {
   if tree[idx].mn >= x {
       return
   }
   if tree[idx].mx < x {
       applyInv(idx)
       return
   }
   push(idx)
   m := (l + r) >> 1
   updateLess(idx*2, l, m, x)
   updateLess(idx*2+1, m, r, x)
   tree[idx].mn = min(tree[idx*2].mn, tree[idx*2+1].mn)
   tree[idx].mx = max(tree[idx*2].mx, tree[idx*2+1].mx)
}

func getVal(idx, l, r, pos int) int {
   if r-l == 1 {
       return tree[idx].mn
   }
   push(idx)
   m := (l + r) >> 1
   if pos < m {
       return getVal(idx*2, l, m, pos)
   }
   return getVal(idx*2+1, m, r, pos)
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n, q int
   fmt.Fscan(in, &n, &q)
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &a[i])
   }

   // build segment tree
   tree = make([]Node, 4*N)
   build(1, 0, N)

   for i := 0; i < q; i++ {
       var s byte
       var x int
       fmt.Fscan(in, &s, &x)
       if s == '>' {
           updateGreater(1, 0, N, x)
       } else {
           updateLess(1, 0, N, x)
       }
   }

   for i, v := range a {
       pos := v + Offset
       res := getVal(1, 0, N, pos)
       if i > 0 {
           out.WriteByte(' ')
       }
       fmt.Fprint(out, res)
   }
   out.WriteByte('\n')
}

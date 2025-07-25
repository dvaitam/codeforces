package main

import (
   "bufio"
   "fmt"
   "os"
)

const MAXV = 1000005

type Node struct {
   sum, maxSuf int
}

var (
   tree []Node
)

func build(idx, l, r int, D []int) {
   if l == r {
       v := D[l]
       tree[idx] = Node{sum: v, maxSuf: max(v, 0)}
       return
   }
   mid := (l + r) >> 1
   build(idx<<1, l, mid, D)
   build(idx<<1|1, mid+1, r, D)
   left, right := tree[idx<<1], tree[idx<<1|1]
   sum := left.sum + right.sum
   // max suffix in this segment: max(right.maxSuf, right.sum + left.maxSuf)
   ms := max(right.maxSuf, right.sum+left.maxSuf)
   tree[idx] = Node{sum: sum, maxSuf: ms}
}

func update(idx, l, r, pos, delta int) {
   if l == r {
       tree[idx].sum += delta
       tree[idx].maxSuf = max(tree[idx].sum, 0)
       return
   }
   mid := (l + r) >> 1
   if pos <= mid {
       update(idx<<1, l, mid, pos, delta)
   } else {
       update(idx<<1|1, mid+1, r, pos, delta)
   }
   left, right := tree[idx<<1], tree[idx<<1|1]
   tree[idx].sum = left.sum + right.sum
   tree[idx].maxSuf = max(right.maxSuf, right.sum+left.maxSuf)
}

// find rightmost position c in [l..r] s.t. suffix sum D[c..MAXV) + suffAfter > 0
func findPos(idx, l, r, suffAfter int) int {
   if tree[idx].maxSuf+suffAfter <= 0 {
       return -1
   }
   if l == r {
       return l
   }
   mid := (l + r) >> 1
   // try right child first
   // right child covers [mid+1..r], suffAfter unchanged
   if res := findPos(idx<<1|1, mid+1, r, suffAfter); res != -1 {
       return res
   }
   // then left child, suffAfter includes right.sum
   return findPos(idx<<1, l, mid, suffAfter+tree[idx<<1|1].sum)
}

func max(a, b int) int {
   if a > b {
       return a
   }
   return b
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   var n, m, q int
   fmt.Fscan(in, &n, &m)
   a := make([]int, n)
   b := make([]int, m)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &a[i])
   }
   for i := 0; i < m; i++ {
       fmt.Fscan(in, &b[i])
   }
   fmt.Fscan(in, &q)
   DA := make([]int, MAXV)
   DB := make([]int, MAXV)
   for _, v := range a {
       DA[v]++
   }
   for _, v := range b {
       DB[v]++
   }
   // D[x] = DA[x] - DB[x]
   D := make([]int, MAXV)
   for i := 0; i < MAXV; i++ {
       D[i] = DA[i] - DB[i]
   }
   tree = make([]Node, MAXV*4)
   build(1, 0, MAXV-1, D)
   for i := 0; i < q; i++ {
       var tp, idx, x int
       fmt.Fscan(in, &tp, &idx, &x)
       if tp == 1 {
           // dish
           old := a[idx-1]
           // update DA[old]--, D[old]--
           update(1, 0, MAXV-1, old, -1)
           // update DA[x]++, D[x]++
           update(1, 0, MAXV-1, x, 1)
           a[idx-1] = x
       } else {
           // pupil
           old := b[idx-1]
           // DB[old]-- => D[old]++
           update(1, 0, MAXV-1, old, 1)
           // DB[x]++ => D[x]--
           update(1, 0, MAXV-1, x, -1)
           b[idx-1] = x
       }
       if tree[1].maxSuf <= 0 {
           fmt.Fprintln(out, -1)
       } else {
           ans := findPos(1, 0, MAXV-1, 0)
           fmt.Fprintln(out, ans)
       }
   }
}

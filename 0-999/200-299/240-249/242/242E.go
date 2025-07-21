package main

import (
   "bufio"
   "fmt"
   "os"
)

const BITS = 20

type Node struct {
   cnt  [BITS]int
   lazy int
}

var (
   n, m  int
   a      []int
   tree   []Node
)

func build(pos, l, r int) {
   if l == r {
       v := a[l]
       for b := 0; b < BITS; b++ {
           if (v>>b)&1 == 1 {
               tree[pos].cnt[b] = 1
           }
       }
       return
   }
   mid := (l + r) >> 1
   lc, rc := pos<<1, pos<<1|1
   build(lc, l, mid)
   build(rc, mid+1, r)
   for b := 0; b < BITS; b++ {
       tree[pos].cnt[b] = tree[lc].cnt[b] + tree[rc].cnt[b]
   }
}

// apply XOR x to node pos covering [l,r]
func applyXor(pos, l, r, x int) {
   length := r - l + 1
   for b := 0; b < BITS; b++ {
       if (x>>b)&1 == 1 {
           tree[pos].cnt[b] = length - tree[pos].cnt[b]
       }
   }
   tree[pos].lazy ^= x
}

func push(pos, l, r int) {
   if tree[pos].lazy != 0 && l < r {
       mid := (l + r) >> 1
       lc, rc := pos<<1, pos<<1|1
       applyXor(lc, l, mid, tree[pos].lazy)
       applyXor(rc, mid+1, r, tree[pos].lazy)
       tree[pos].lazy = 0
   }
}

func update(pos, l, r, ql, qr, x int) {
   if ql > r || qr < l {
       return
   }
   if ql <= l && r <= qr {
       applyXor(pos, l, r, x)
       return
   }
   push(pos, l, r)
   mid := (l + r) >> 1
   lc, rc := pos<<1, pos<<1|1
   if ql <= mid {
       update(lc, l, mid, ql, qr, x)
   }
   if qr > mid {
       update(rc, mid+1, r, ql, qr, x)
   }
   for b := 0; b < BITS; b++ {
       tree[pos].cnt[b] = tree[lc].cnt[b] + tree[rc].cnt[b]
   }
}

func query(pos, l, r, ql, qr int) int64 {
   if ql > r || qr < l {
       return 0
   }
   if ql <= l && r <= qr {
       var sum int64
       for b := 0; b < BITS; b++ {
           if tree[pos].cnt[b] > 0 {
               sum += int64(tree[pos].cnt[b]) << b
           }
       }
       return sum
   }
   push(pos, l, r)
   mid := (l + r) >> 1
   var res int64
   if ql <= mid {
       res += query(pos<<1, l, mid, ql, qr)
   }
   if qr > mid {
       res += query(pos<<1|1, mid+1, r, ql, qr)
   }
   return res
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   fmt.Fscan(reader, &n)
   a = make([]int, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   tree = make([]Node, 4*(n+1))
   build(1, 1, n)
   fmt.Fscan(reader, &m)
   for i := 0; i < m; i++ {
       var t, l, r, x int
       fmt.Fscan(reader, &t, &l, &r)
       if t == 1 {
           res := query(1, 1, n, l, r)
           fmt.Fprintln(writer, res)
       } else {
           fmt.Fscan(reader, &x)
           update(1, 1, n, l, r, x)
       }
   }
}

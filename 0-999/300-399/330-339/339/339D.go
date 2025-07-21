package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n, m int
   fmt.Fscan(in, &n, &m)
   size := 1 << n
   a := make([]int, size+1)
   for i := 1; i <= size; i++ {
       fmt.Fscan(in, &a[i])
   }

   tree := make([]int, 4*size)
   // rootOp indicates operation at the root: true for OR, false for XOR
   rootOp := n%2 == 1

   var build func(node, l, r int, op bool)
   build = func(node, l, r int, op bool) {
       if l == r {
           tree[node] = a[l]
           return
       }
       mid := (l + r) >> 1
       // children use opposite operation
       build(node<<1, l, mid, !op)
       build(node<<1|1, mid+1, r, !op)
       if op {
           tree[node] = tree[node<<1] | tree[node<<1|1]
       } else {
           tree[node] = tree[node<<1] ^ tree[node<<1|1]
       }
   }
   build(1, 1, size, rootOp)

   var update func(node, l, r, pos, val int, op bool)
   update = func(node, l, r, pos, val int, op bool) {
       if l == r {
           tree[node] = val
           return
       }
       mid := (l + r) >> 1
       if pos <= mid {
           update(node<<1, l, mid, pos, val, !op)
       } else {
           update(node<<1|1, mid+1, r, pos, val, !op)
       }
       if op {
           tree[node] = tree[node<<1] | tree[node<<1|1]
       } else {
           tree[node] = tree[node<<1] ^ tree[node<<1|1]
       }
   }

   for i := 0; i < m; i++ {
       var p, b int
       fmt.Fscan(in, &p, &b)
       update(1, 1, size, p, b, rootOp)
       fmt.Fprintln(out, tree[1])
   }
}

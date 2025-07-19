package main

import (
   "bufio"
   "fmt"
   "os"
)

type Node struct {
   m [3][3]float64
}

var (
   a      []int
   tree   []Node
   x, y   int
   negInf = -1e18
)

func combine(res *Node, left, right Node) {
   for i := 0; i < 3; i++ {
       for j := 0; j < 3; j++ {
           // transitions: (0,0),(1,0),(0,1),(1,1),(2,0),(0,2)
           maxv := left.m[i][0] + right.m[0][j]
           if v := left.m[i][1] + right.m[0][j]; v > maxv {
               maxv = v
           }
           if v := left.m[i][0] + right.m[1][j]; v > maxv {
               maxv = v
           }
           if v := left.m[i][1] + right.m[1][j]; v > maxv {
               maxv = v
           }
           if v := left.m[i][2] + right.m[0][j]; v > maxv {
               maxv = v
           }
           if v := left.m[i][0] + right.m[2][j]; v > maxv {
               maxv = v
           }
           res.m[i][j] = maxv
       }
   }
}

func build(idx, l, r int) {
   if r-l == 1 {
       // leaf
       for i := 0; i < 3; i++ {
           for j := 0; j < 3; j++ {
               tree[idx].m[i][j] = negInf
           }
       }
       tree[idx].m[0][0] = 0
       tree[idx].m[1][1] = float64(a[l]) / float64(x+y)
       tree[idx].m[2][2] = float64(a[l]) / float64(x)
       return
   }
   mid := (l + r) >> 1
   build(idx*2, l, mid)
   build(idx*2+1, mid, r)
   combine(&tree[idx], tree[idx*2], tree[idx*2+1])
}

func update(idx, l, r, pos, val int) {
   if pos < l || pos >= r {
       return
   }
   if r-l == 1 {
       a[pos] = val
       for i := 0; i < 3; i++ {
           for j := 0; j < 3; j++ {
               tree[idx].m[i][j] = negInf
           }
       }
       tree[idx].m[0][0] = 0
       tree[idx].m[1][1] = float64(val) / float64(x+y)
       tree[idx].m[2][2] = float64(val) / float64(x)
       return
   }
   mid := (l + r) >> 1
   update(idx*2, l, mid, pos, val)
   update(idx*2+1, mid, r, pos, val)
   combine(&tree[idx], tree[idx*2], tree[idx*2+1])
}

func query(idx, l, r, ql, qr int) Node {
   if r <= ql || l >= qr {
       var empty Node
       empty.m[0][0] = -1
       return empty
   }
   if l >= ql && r <= qr {
       return tree[idx]
   }
   mid := (l + r) >> 1
   left := query(idx*2, l, mid, ql, qr)
   right := query(idx*2+1, mid, r, ql, qr)
   if left.m[0][0] == -1 {
       return right
   }
   if right.m[0][0] == -1 {
       return left
   }
   var res Node
   combine(&res, left, right)
   return res
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, q int
   fmt.Fscan(reader, &n, &q)
   fmt.Fscan(reader, &x, &y)
   if x < y {
       x, y = y, x
   }
   a = make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   tree = make([]Node, 4*n)
   build(1, 0, n)
   for qi := 0; qi < q; qi++ {
       var typ int
       fmt.Fscan(reader, &typ)
       if typ == 1 {
           var pos, v int
           fmt.Fscan(reader, &pos, &v)
           pos--
           update(1, 0, n, pos, v)
       } else {
           var l, r int
           fmt.Fscan(reader, &l, &r)
           l--
           ans := query(1, 0, n, l, r)
           res := -1e18
           for i := 0; i < 3; i++ {
               for j := 0; j < 3; j++ {
                   if ans.m[i][j] > res {
                       res = ans.m[i][j]
                   }
               }
           }
           fmt.Fprintf(writer, "%.15f\n", res)
       }
   }
}

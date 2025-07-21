package main

import (
   "bufio"
   "fmt"
   "os"
)

const INF = 1 << 60

var (
   a    []int64
   b    []int64
   lazy []int64
   n    int
)

// build initializes the segment tree for size n0
func build(n0 int) {
   n = n0
   lazy = make([]int64, 4*n+10)
   for i := range lazy {
       lazy[i] = INF
   }
}

// update assigns offset to range [L,R] in b
func update(node, l, r, L, R int, offset int64) {
   if L <= l && r <= R {
       lazy[node] = offset
       return
   }
   if lazy[node] != INF {
       // push down existing tag
       lazy[node*2] = lazy[node]
       lazy[node*2+1] = lazy[node]
       lazy[node] = INF
   }
   mid := (l + r) >> 1
   if L <= mid {
       update(node*2, l, mid, L, R, offset)
   }
   if R > mid {
       update(node*2+1, mid+1, r, L, R, offset)
   }
}

// query returns the offset tag at position pos, or INF if none
func query(node, l, r, pos int) int64 {
   if lazy[node] != INF || l == r {
       return lazy[node]
   }
   mid := (l + r) >> 1
   if pos <= mid {
       return query(node*2, l, mid, pos)
   }
   return query(node*2+1, mid+1, r, pos)
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var m int
   fmt.Fscan(reader, &n, &m)
   a = make([]int64, n+1)
   b = make([]int64, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &b[i])
   }
   build(n)
   for i := 0; i < m; i++ {
       var t int
       fmt.Fscan(reader, &t)
       if t == 1 {
           var x, y, k int
           fmt.Fscan(reader, &x, &y, &k)
           offset := int64(x - y)
           update(1, 1, n, y, y+k-1, offset)
       } else {
           var x int
           fmt.Fscan(reader, &x)
           off := query(1, 1, n, x)
           if off == INF {
               fmt.Fprintln(writer, b[x])
           } else {
               idx := int64(x) + off
               fmt.Fprintln(writer, a[idx])
           }
       }
   }
}

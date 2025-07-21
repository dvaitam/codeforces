package main

import (
   "bufio"
   "fmt"
   "os"
)

func min(a, b int) int {
   if a < b {
       return a
   }
   return b
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   fmt.Fscan(reader, &n)
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   // build segment tree for range minimum query over a
   size := 1
   for size < n {
       size <<= 1
   }
   seg := make([]int, size*2)
   const INF = 1 << 60
   // init leaves
   for i := 0; i < size*2; i++ {
       seg[i] = INF
   }
   for i := 0; i < n; i++ {
       seg[size+i] = a[i]
   }
   // build internal
   for i := size - 1; i > 0; i-- {
       seg[i] = min(seg[i<<1], seg[i<<1|1])
   }

   // rmq on [l,r] inclusive, 0-based indices
   rmq := func(l, r int) int {
       l += size
       r += size
       res := INF
       for l <= r {
           if l&1 == 1 {
               res = min(res, seg[l])
               l++
           }
           if r&1 == 0 {
               res = min(res, seg[r])
               r--
           }
           l >>= 1
           r >>= 1
       }
       return res
   }

   var q int
   fmt.Fscan(reader, &q)
   for qi := 0; qi < q; qi++ {
       var r1, c1, r2, c2 int
       fmt.Fscan(reader, &r1, &c1, &r2, &c2)
       r1--
       r2--
       // same row
       if r1 == r2 {
           arow := a[r1]
           // horizontal distance on same row
           d := abs(c1 - c2)
           // via home then right
           d = min(d, 1 + c2)
           // via end then left
           d = min(d, 1 + (arow - c2))
           fmt.Fprintln(writer, d)
           continue
       }
       // vertical distance
       vd := abs(r1 - r2)
       // min a in rows between r1 and r2 (excluding r1, including r2 if moving down; excluding r1, including r2 if up)
       var mn int
       if r1 < r2 {
           mn = rmq(r1+1, r2)
       } else {
           mn = rmq(r2, r1-1)
       }
       a1 := a[r1]
       a2 := a[r2]
       best := INF
       // options for initial column x0: c1, 0, a1
       opts := []int{c1, 0, a1}
       for _, x0 := range opts {
           // cost to move from c1 to x0 on r1
           c1cost := abs(c1 - x0)
           c1cost = min(c1cost, 1 + x0)
           c1cost = min(c1cost, 1 + (a1 - x0))
           // after vertical moves, column becomes y
           y := x0
           if y > mn {
               y = mn
           }
           // cost to move from y to c2 on r2
           c2cost := abs(y - c2)
           c2cost = min(c2cost, 1 + c2)
           c2cost = min(c2cost, 1 + (a2 - c2))
           total := vd + c1cost + c2cost
           if total < best {
               best = total
           }
       }
       fmt.Fprintln(writer, best)
   }
}

func abs(x int) int {
   if x < 0 {
       return -x
   }
   return x
}

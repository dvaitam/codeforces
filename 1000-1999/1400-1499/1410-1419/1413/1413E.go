package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var t int
   if _, err := fmt.Fscan(in, &t); err != nil {
       return
   }
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   for ; t > 0; t-- {
       var a, b, c, d int64
       fmt.Fscan(in, &a, &b, &c, &d)
       // infinite if per-cast net damage positive
       if a > b*c {
           fmt.Fprintln(out, -1)
           continue
       }
       // precompute
       // w = floor(c/d), u = c mod d
       w := c / d
       u := c - w*d
       // function to compute S at given n, r
       compute := func(n, r int64) int64 {
           // t = n*d + r
           // number of casts = n+1
           // m = floor((t - c)/d)
           tval := n*d + r
           m := (tval - c) / d
           if tval < c {
               m = -1
           }
           if m > n {
               m = n
           }
           // full heals count
           full := m + 1
           if full < 0 {
               full = 0
           }
           // partial count K = n - m_eff
           K := n - m
           // sum of heal: full*c + sum_{j=0..K-1}(j*d + r)
           // sum j from 0..K-1 = (K-1)*K/2
           sumJR := (K - 1) * K / 2 * d
           sumR := r * K
           heal := full*c + sumJR + sumR
           // total heal multiplied by b
           heal *= b
           // damage
           dmg := (n + 1) * a
           return dmg - heal
       }
       // candidates
       best := a // at t=0
       // candidates for r = 0, n in {n1-1, n1, w}
       // n1 ~ floor(a/(b*d))
       n1 := a / (b * d)
       // try n1 - 1, n1, w
       cand := []int64{0, n1 - 1, n1, w}
       for _, n := range cand {
           if n < 0 || n > w {
               continue
           }
           s := compute(n, 0)
           if s > best {
               best = s
           }
       }
       // also candidate at t = c: n=w, r=u
       if u > 0 {
           s := compute(w, u)
           if s > best {
               best = s
           }
       }
       fmt.Fprintln(out, best)
   }
}

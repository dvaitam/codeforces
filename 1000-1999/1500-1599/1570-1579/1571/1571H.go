package main

import (
   "bufio"
   "fmt"
   "os"
)

const MOD = 998244353

func modAdd(a, b int) int {
   res := a + b
   if res >= MOD {
       res -= MOD
   }
   return res
}

func modSub(a, b int) int {
   res := a - b
   if res < 0 {
       res += MOD
   }
   return res
}

func modMul(a, b int) int {
   return int((int64(a) * int64(b)) % MOD)
}

func modPow(a, e int) int {
   res := 1
   base := a
   for e > 0 {
       if e&1 == 1 {
           res = modMul(res, base)
       }
       base = modMul(base, base)
       e >>= 1
   }
   return res
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, a, b int
   fmt.Fscan(reader, &n, &a, &b)
   xs := make([]int, n)
   ys := make([]int, n)
   ps := make([]int, n)
   inv1e6 := modPow(1000000, MOD-2)
   for i := 0; i < n; i++ {
       var x, y, pi int
       fmt.Fscan(reader, &x, &y, &pi)
       xs[i] = x
       ys[i] = y
       ps[i] = modMul(pi%MOD, inv1e6)
   }
   // distinct horizontal and vertical lines
   hMap := make(map[int]bool)
   vMap := make(map[int]bool)
   // diagonal keys
   d1Map := make(map[int]bool)
   d2Map := make(map[int]bool)
   for i := 0; i < n; i++ {
       vMap[xs[i]] = true
       hMap[ys[i]] = true
       d1Map[xs[i]-ys[i]] = true
       d2Map[xs[i]+ys[i]] = true
   }
   // slices of distinct values
   hs := make([]int, 0, len(hMap))
   vs := make([]int, 0, len(vMap))
   c1s := make([]int, 0, len(d1Map))
   c2s := make([]int, 0, len(d2Map))
   for y := range hMap {
       hs = append(hs, y)
   }
   for x := range vMap {
       vs = append(vs, x)
   }
   for c1 := range d1Map {
       c1s = append(c1s, c1)
   }
   for c2 := range d2Map {
       c2s = append(c2s, c2)
   }
   k1 := len(c1s)
   k2 := len(c2s)
   // map diagonal value to its index
   c1Index := make(map[int]int, k1)
   c2Index := make(map[int]int, k2)
   for i, c1 := range c1s {
       c1Index[c1] = i
   }
   for i, c2 := range c2s {
       c2Index[c2] = i
   }
   // point indices for each diagonal key
   s1 := make([][]int, k1)
   s2 := make([][]int, k2)
   for i := 0; i < n; i++ {
       i1 := c1Index[xs[i]-ys[i]]
       s1[i1] = append(s1[i1], i)
       i2 := c2Index[xs[i]+ys[i]]
       s2[i2] = append(s2[i2], i)
   }
   // compute F1, F2: prob no line
   F1 := make([]int, k1)
   for i := 0; i < k1; i++ {
       prod := 1
       for _, j := range s1[i] {
           prod = modMul(prod, modSub(1, ps[j]))
       }
       F1[i] = prod
   }
   F2 := make([]int, k2)
   for i := 0; i < k2; i++ {
       prod := 1
       for _, j := range s2[i] {
           prod = modMul(prod, modSub(1, ps[j]))
       }
       F2[i] = prod
   }
   // presence probabilities
   P1 := make([]int, k1)
   for i := 0; i < k1; i++ {
       P1[i] = modSub(1, F1[i])
   }
   P2 := make([]int, k2)
   for i := 0; i < k2; i++ {
       P2[i] = modSub(1, F2[i])
   }
   // Pboth for D1-D2: probability both diagonal lines present
   P12 := make([][]int, k1)
   for i := 0; i < k1; i++ {
       P12[i] = make([]int, k2)
       for j := 0; j < k2; j++ {
           // compute F12 = prod of (1-p) over union of s1[i] and s2[j]
           include := make([]bool, n)
           for _, idx := range s1[i] {
               include[idx] = true
           }
           for _, idx := range s2[j] {
               include[idx] = true
           }
           prod := 1
           for t := 0; t < n; t++ {
               if include[t] {
                   prod = modMul(prod, modSub(1, ps[t]))
               }
           }
           // P12 = 1 - F1 - F2 + F12
           tmp := modAdd(F1[i], F2[j])
           tmp = modSub(1, modSub(tmp, prod))
           P12[i][j] = tmp
       }
   }
   // expectation of number of lines
   expLines := modAdd(len(hs), len(vs))
   for i := 0; i < k1; i++ {
       expLines = modAdd(expLines, P1[i])
   }
   for i := 0; i < k2; i++ {
       expLines = modAdd(expLines, P2[i])
   }
   // compute expectation of intersections via Euler: sum_p (E[r_p] - P(r_p>=1))
   const factor = 100000
   pts := make(map[int]struct{})
   // collect all potential intersection points
   // H-V
   for _, x := range vs {
       for _, y := range hs {
           key := (2*x)*factor + 2*y
           pts[key] = struct{}{}
       }
   }
   // H-D1, H-D2
   for _, y := range hs {
       for _, c1 := range c1s {
           x := c1 + y
           if 0 <= x && x <= a {
               key := (2*x)*factor + 2*y
               pts[key] = struct{}{}
           }
       }
       for _, c2 := range c2s {
           x := c2 - y
           if 0 <= x && x <= a {
               key := (2*x)*factor + 2*y
               pts[key] = struct{}{}
           }
       }
   }
   // V-D1, V-D2
   for _, x := range vs {
       for _, c1 := range c1s {
           y0 := x - c1
           if 0 <= y0 && y0 <= b {
               key := (2*x)*factor + 2*y0
               pts[key] = struct{}{}
           }
       }
       for _, c2 := range c2s {
           y0 := c2 - x
           if 0 <= y0 && y0 <= b {
               key := (2*x)*factor + 2*y0
               pts[key] = struct{}{}
           }
       }
   }
   // D1-D2
   for _, c1 := range c1s {
       for _, c2 := range c2s {
           x2 := c1 + c2
           y2 := c2 - c1
           if 0 <= x2 && x2 <= 2*a && 0 <= y2 && y2 <= 2*b {
               key := x2*factor + y2
               pts[key] = struct{}{}
           }
       }
   }
   sumCore := 0
   // iterate points
   for key := range pts {
       x2 := key / factor
       y2 := key % factor
       // determine line candidates at this point
       sumPr := 0
       prodNone := 1
       // horizontal
       if y2%2 == 0 {
           y := y2 / 2
           if hMap[y] {
               sumPr = modAdd(sumPr, 1)
               prodNone = 0
           }
       }
       // vertical
       if x2%2 == 0 {
           x := x2 / 2
           if vMap[x] {
               sumPr = modAdd(sumPr, 1)
               prodNone = 0
           }
       }
       // D1
       if (x2-y2)%2 == 0 {
           c1 := (x2 - y2) / 2
           if idx, ok := c1Index[c1]; ok {
               sumPr = modAdd(sumPr, P1[idx])
               prodNone = modMul(prodNone, F1[idx])
           }
       }
       // D2
       if (x2+y2)%2 == 0 {
           c2 := (x2 + y2) / 2
           if idx, ok := c2Index[c2]; ok {
               sumPr = modAdd(sumPr, P2[idx])
               prodNone = modMul(prodNone, F2[idx])
           }
       }
       // probability at least one line
       pAny := modSub(1, prodNone)
       // add E[r_p] - P(r_p>=1)
       sumCore = modAdd(sumCore, modSub(sumPr, pAny))
   }
   res := modAdd(1, modAdd(expLines, sumCore))
   fmt.Fprintln(writer, res)
}

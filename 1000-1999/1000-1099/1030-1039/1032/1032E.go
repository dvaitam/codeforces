package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   fmt.Fscan(in, &n)
   a := make([]int, n)
   totalSum := 0
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &a[i])
       totalSum += a[i]
   }
   cntMap := make(map[int]int)
   for _, v := range a {
       cntMap[v]++
   }
   wVals := make([]int, 0, len(cntMap))
   wCnts := make([]int, 0, len(cntMap))
   for v, c := range cntMap {
       wVals = append(wVals, v)
       wCnts = append(wCnts, c)
   }
   D := len(wVals)
   B := (totalSum + 64) / 64
   dpFull := make([][]uint64, n+1)
   for t := 0; t <= n; t++ {
       dpFull[t] = make([]uint64, B)
   }
   dpFull[0][0] = 1
   for idx := 0; idx < D; idx++ {
       w := wVals[idx]
       c := wCnts[idx]
       for rep := 0; rep < c; rep++ {
           for t := n - 1; t >= 0; t-- {
               src := dpFull[t]
               dst := dpFull[t+1]
               shift := w
               ws := shift / 64
               bs := uint(shift % 64)
               if bs == 0 {
                   for j := B - 1; j >= ws; j-- {
                       dst[j] |= src[j-ws]
                   }
               } else {
                   for j := B - 1; j >= ws; j-- {
                       v := src[j-ws] << bs
                       if j-ws-1 >= 0 {
                           v |= src[j-ws-1] >> (64 - bs)
                       }
                       dst[j] |= v
                   }
               }
           }
       }
   }
   reachable := make([][]bool, n+1)
   for t := 0; t <= n; t++ {
       reachable[t] = make([]bool, totalSum+1)
       for b := 0; b < B; b++ {
           block := dpFull[t][b]
           if block == 0 {
               continue
           }
           for bit := 0; bit < 64; bit++ {
               if block&(1<<uint(bit)) != 0 {
                   s := b*64 + bit
                   if s <= totalSum {
                       reachable[t][s] = true
                   }
               }
           }
       }
   }
   varArr := make([][]int, n+1)
   for t := 0; t <= n; t++ {
       varArr[t] = make([]int, totalSum+1)
   }
   for idx := 0; idx < D; idx++ {
       dpExcl := make([][]uint64, n+1)
       for t := 0; t <= n; t++ {
           dpExcl[t] = make([]uint64, B)
       }
       dpExcl[0][0] = 1
       for j := 0; j < D; j++ {
           if j == idx {
               continue
           }
           w := wVals[j]
           c := wCnts[j]
           for rep := 0; rep < c; rep++ {
               for t := n - 1; t >= 0; t-- {
                   src := dpExcl[t]
                   dst := dpExcl[t+1]
                   shift := w
                   ws := shift / 64
                   bs := uint(shift % 64)
                   if bs == 0 {
                       for u := B - 1; u >= ws; u-- {
                           dst[u] |= src[u-ws]
                       }
                   } else {
                       for u := B - 1; u >= ws; u-- {
                           v := src[u-ws] << bs
                           if u-ws-1 >= 0 {
                               v |= src[u-ws-1] >> (64 - bs)
                           }
                           dst[u] |= v
                       }
                   }
               }
           }
       }
       w := wVals[idx]
       c := wCnts[idx]
       for t := 0; t <= n; t++ {
           for s := 0; s <= totalSum; s++ {
               if !reachable[t][s] {
                   continue
               }
               minx, maxx := c+1, -1
               for x := 0; x <= c; x++ {
                   tt := t - x
                   ss := s - w*x
                   if tt < 0 || ss < 0 {
                       break
                   }
                   b := ss / 64
                   bit := uint(ss % 64)
                   if dpExcl[tt][b]&(1<<bit) != 0 {
                       if x < minx {
                           minx = x
                       }
                       if x > maxx {
                           maxx = x
                       }
                   }
               }
               varArr[t][s] += maxx - minx
           }
       }
   }
   bestVar := n
   for t := 0; t <= n; t++ {
       for s := 0; s <= totalSum; s++ {
           if !reachable[t][s] {
               continue
           }
           if varArr[t][s] < bestVar {
               bestVar = varArr[t][s]
           }
       }
   }
   fmt.Println(n - bestVar)
}

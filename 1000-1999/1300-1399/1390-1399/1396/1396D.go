package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

const MOD = 1000000007

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, k int
   var L int64
   fmt.Fscan(in, &n, &k, &L)
   xs := make([]int64, n)
   ys := make([]int64, n)
   cs := make([]int, n)
   for i := 0; i < n; i++ {
       var x, y int64
       var c int
       fmt.Fscan(in, &x, &y, &c)
       xs[i] = x
       ys[i] = y
       cs[i] = c - 1
   }
   // compress x, y
   ux := uniqueSortInt64(xs)
   uy := uniqueSortInt64(ys)
   mx := len(ux)
   my := len(uy)
   xIdx := make(map[int64]int)
   for i, v := range ux {
       xIdx[v] = i
   }
   yIdx := make(map[int64]int)
   for i, v := range uy {
       yIdx[v] = i
   }
   // boundary choices
   leftX := make([]int64, mx)
   rightX := make([]int64, mx)
   for i := 0; i < mx; i++ {
       prev := int64(-1)
       if i > 0 {
           prev = ux[i-1]
       }
       leftX[i] = (ux[i] - prev) % MOD
       next := L
       if i+1 < mx {
           next = ux[i+1]
       }
       rightX[i] = (next - ux[i]) % MOD
   }
   leftY := make([]int64, my)
   rightY := make([]int64, my)
   for i := 0; i < my; i++ {
       prev := int64(-1)
       if i > 0 {
           prev = uy[i-1]
       }
       leftY[i] = (uy[i] - prev) % MOD
       next := L
       if i+1 < my {
           next = uy[i+1]
       }
       rightY[i] = (next - uy[i]) % MOD
   }
   // points per x idx
   vx := make([][]struct{ y, c int }, mx)
   for i := 0; i < n; i++ {
       xi := xIdx[xs[i]]
       yi := yIdx[ys[i]]
       vx[xi] = append(vx[xi], struct{ y, c int }{yi, cs[i]})
   }
   var ans int64
   // freq per color
   freq := make([]int, k)
   // rowColors per y idx
   rowColors := make([][]int, my)
   // used colors for reset
   used := make([]int, 0, n)
   for l := 0; l < mx; l++ {
       // reset rowColors
       for i := 0; i < my; i++ {
           rowColors[i] = rowColors[i][:0]
       }
       for r := l; r < mx; r++ {
           // add points at x=r
           for _, p := range vx[r] {
               rowColors[p.y] = append(rowColors[p.y], p.c)
           }
           // prefix sum of rightY
           sRight := make([]int64, my+1)
           for i := my - 1; i >= 0; i-- {
               sRight[i] = (rightY[i] + sRight[i+1]) % MOD
           }
           // sliding window on y
           used = used[:0]
           have := 0
           for i := range freq {
               freq[i] = 0
           }
           b, t := 0, 0
           lx := leftX[l]
           rx := rightX[r]
           for b < my {
               for t < my && have < k {
                   for _, c := range rowColors[t] {
                       if freq[c] == 0 {
                           have++
                           used = append(used, c)
                       }
                       freq[c]++
                   }
                   t++
               }
               if have == k {
                   yCnt := leftY[b] * sRight[t-1] % MOD
                   ans = (ans + lx*rx%MOD*yCnt%MOD) % MOD
               }
               // remove b
               for _, c := range rowColors[b] {
                   freq[c]--
                   if freq[c] == 0 {
                       have--
                   }
               }
               b++
           }
       }
   }
   fmt.Println(ans)
}

func uniqueSortInt64(a []int64) []int64 {
   m := make(map[int64]struct{}, len(a))
   for _, v := range a {
       m[v] = struct{}{}
   }
   res := make([]int64, 0, len(m))
   for v := range m {
       res = append(res, v)
   }
   sort.Slice(res, func(i, j int) bool { return res[i] < res[j] })
   return res
}

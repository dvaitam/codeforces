package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   fmt.Fscan(in, &n)
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &a[i])
   }
   // collect occurrences of each lucky value
   first := make(map[int]int)
   last := make(map[int]int)
   for i, x := range a {
       if isLucky(x) {
           pos := i + 1
           if _, ok := first[x]; !ok {
               first[x] = pos
           }
           last[x] = pos
       }
   }
   // no lucky values: count all non-intersecting pairs
   if len(first) == 0 {
       var total int64
       for k := 1; k < n; k++ {
           total += int64(k) * int64(n-k)
       }
       fmt.Println(total)
       return
   }
   // collect critical points
   pts := make([]int, 0, 2*len(first)+2)
   pts = append(pts, 1, n)
   for x, f := range first {
       l := last[x]
       pts = append(pts, f, l+1)
   }
   // unique and sort
   sort.Ints(pts)
   uniq := pts[:0]
   for _, v := range pts {
       if len(uniq) == 0 || uniq[len(uniq)-1] != v {
           uniq = append(uniq, v)
       }
   }
   // iterate intervals
   var ans int64
   vals := make([][2]int, 0, len(first))
   for x, f := range first {
       vals = append(vals, [2]int{f, last[x]})
   }
   for i := 0; i+1 < len(uniq); i++ {
       L := uniq[i]
       R := uniq[i+1] - 1
       if L > R || L >= n {
           continue
       }
       if R > n-1 {
           R = n - 1
       }
       lenK := R - L + 1
       // sum k and sum k^2 over [L..R]
       sumK := int64(L+R) * int64(lenK) / 2
       sumK2 := sumSquares(R) - sumSquares(L-1)
       total := int64(n)*sumK - sumK2
       // build S for k=L
       rects := make([][2]int, 0, len(vals))
       for _, p := range vals {
           if p[0] <= L && p[1] > L {
               rects = append(rects, p)
           }
       }
       bad := unionArea(rects, n)
       ans += total - bad*int64(lenK)
   }
   fmt.Println(ans)
}

func sumSquares(x int) int64 {
   if x <= 0 {
       return 0
   }
   xx := int64(x)
   return xx * (xx + 1) * (2*xx + 1) / 6
}

// unionArea computes union area of rectangles of form [1..u] x [v..n]
func unionArea(rects [][2]int, n int) int64 {
   if len(rects) == 0 {
       return 0
   }
   // sort by u increasing
   sort.Slice(rects, func(i, j int) bool {
       return rects[i][0] < rects[j][0]
   })
   var area int64
   maxV := n + 1
   // iterate from largest u
   for i := len(rects) - 1; i >= 0; i-- {
       u, v := rects[i][0], rects[i][1]
       if v < maxV {
           // new area (u - prev_u) * (n - v + 1)
           area += int64(u) * int64(n-v+1)
           if i < len(rects)-1 {
               // subtract overlap: prev_u * (n - maxV +1)
               prevU := rects[i+1][0]
               area -= int64(prevU) * int64(n-maxV+1)
           }
           if v < maxV {
               maxV = v
           }
       }
   }
   return area
}

func isLucky(x int) bool {
   if x <= 0 {
       return false
   }
   for x > 0 {
       d := x % 10
       if d != 4 && d != 7 {
           return false
       }
       x /= 10
   }
   return true
}

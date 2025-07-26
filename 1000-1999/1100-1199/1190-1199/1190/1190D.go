package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   xs := make([]int64, n)
   ys := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &xs[i], &ys[i])
   }
   // compress x coordinates
   ux := make([]int64, n)
   copy(ux, xs)
   sort.Slice(ux, func(i, j int) bool { return ux[i] < ux[j] })
   ux = unique(ux)
   m := len(ux)
   // points with compressed x index
   pts := make([]struct{ x int; y int64 }, n)
   for i := 0; i < n; i++ {
       xi := sort.Search(len(ux), func(j int) bool { return ux[j] >= xs[i] })
       pts[i] = struct{ x int; y int64 }{x: xi, y: ys[i]}
   }
   // sort by y descending
   sort.Slice(pts, func(i, j int) bool { return pts[i].y > pts[j].y })
   countX := make([]int, m)
   activeCols := 0
   var ans int64
   // process points by groups of equal y
   for i := 0; i < n; {
       yv := pts[i].y
       j := i
       for j < n && pts[j].y == yv {
           xi := pts[j].x
           if countX[xi] == 0 {
               activeCols++
           }
           countX[xi]++
           j++
       }
       // at threshold below this y, activeCols columns are non-empty
       ac := int64(activeCols)
       ans += ac * (ac + 1) / 2
       i = j
   }
   // output
   fmt.Fprintln(os.Stdout, ans)
}

// unique returns sorted unique slice of int64
func unique(a []int64) []int64 {
   if len(a) == 0 {
       return a
   }
   j := 1
   for i := 1; i < len(a); i++ {
       if a[i] != a[i-1] {
           a[j] = a[i]
           j++
       }
   }
   return a[:j]
}

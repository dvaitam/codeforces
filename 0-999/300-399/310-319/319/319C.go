package main

import (
   "bufio"
   "fmt"
   "os"
)

// Line represents y = m*x + c
type Line struct {
   m, c int64
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   a := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   b := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &b[i])
   }
   // dp[i] = min cost to cut tree i (0-based)
   dp := make([]int64, n)
   // Convex hull for lines with decreasing slopes, queries with increasing x
   hull := make([]Line, 0, n)
   pos := 0
   // initial line from tree 1 (index 0)
   hull = append(hull, Line{m: b[0], c: dp[0]})
   for i := 1; i < n; i++ {
       x := a[i]
       // query
       for pos+1 < len(hull) && eval(hull[pos], x) >= eval(hull[pos+1], x) {
           pos++
       }
       dp[i] = eval(hull[pos], x)
       // add line for this i
       newLine := Line{m: b[i], c: dp[i]}
       // maintain hull: remove last if unnecessary
       for len(hull) >= 2 && isBad(hull[len(hull)-2], hull[len(hull)-1], newLine) {
           hull = hull[:len(hull)-1]
           if pos >= len(hull) {
               pos = len(hull) - 1
           }
       }
       hull = append(hull, newLine)
   }
   // answer is dp[n-1]
   fmt.Fprintln(writer, dp[n-1])
}

// eval returns value of line at x
func eval(l Line, x int64) int64 {
   return l.m*x + l.c
}

// isBad checks if line l2 is unnecessary between l1 and l3
func isBad(l1, l2, l3 Line) bool {
   // intersection(l1,l2) >= intersection(l2,l3)
   return (l2.c - l1.c) * (l2.m - l3.m) >= (l3.c - l2.c) * (l1.m - l2.m)
}

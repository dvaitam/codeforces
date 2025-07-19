package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   fmt.Fscan(reader, &n)
   x1s := make([]int, n)
   y1s := make([]int, n)
   x2s := make([]int, n)
   y2s := make([]int, n)

   const INF = 2000000000
   var m [5]int
   var s [5]int
   // m[1]=min x2, m[2]=min y2, m[3]=max x1, m[4]=max y1
   m[1], m[2], s[1], s[2] = INF, INF, INF, INF
   m[3], m[4], s[3], s[4] = -INF, -INF, -INF, -INF

   for i := 0; i < n; i++ {
       var x1, y1, x2, y2 int
       fmt.Fscan(reader, &x1, &y1, &x2, &y2)
       x1s[i], y1s[i], x2s[i], y2s[i] = x1, y1, x2, y2
       // update max of x1 (left boundary)
       if m[3] <= x1 {
           s[3] = m[3]
           m[3] = x1
       } else if x1 >= s[3] {
           s[3] = x1
       }
       // update max of y1 (bottom boundary)
       if m[4] <= y1 {
           s[4] = m[4]
           m[4] = y1
       } else if y1 >= s[4] {
           s[4] = y1
       }
       // update min of x2 (right boundary)
       if m[1] >= x2 {
           s[1] = m[1]
           m[1] = x2
       } else if x2 <= s[1] {
           s[1] = x2
       }
       // update min of y2 (top boundary)
       if m[2] >= y2 {
           s[2] = m[2]
           m[2] = y2
       } else if y2 <= s[2] {
           s[2] = y2
       }
   }

   for i := 0; i < n; i++ {
       // exclude i-th rectangle
       U, R, D, L := m[1], m[2], m[3], m[4]
       if x2s[i] == m[1] {
           U = s[1]
       }
       if y2s[i] == m[2] {
           R = s[2]
       }
       if x1s[i] == m[3] {
           D = s[3]
       }
       if y1s[i] == m[4] {
           L = s[4]
       }
       if D <= U && L <= R {
           fmt.Fprintln(writer, D, L)
           return
       }
   }
}

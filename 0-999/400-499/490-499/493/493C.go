package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m int
   fmt.Fscan(reader, &n)
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   fmt.Fscan(reader, &m)
   b := make([]int, m)
   for i := 0; i < m; i++ {
       fmt.Fscan(reader, &b[i])
   }
   sort.Ints(a)
   sort.Ints(b)

   // initial scores: all throws > 0 (d=0) are worth 3
   score1 := 3 * n
   score2 := 3 * m
   best1, best2 := score1, score2
   bestDiff := score1 - score2

   i, j := 0, 0
   // process thresholds at each distinct distance
   for i < n || j < m {
       // next threshold
       var d int
       if i < n && (j >= m || a[i] < b[j]) {
           d = a[i]
       } else if j < m && (i >= n || b[j] < a[i]) {
           d = b[j]
       } else {
           // a[i] == b[j]
           d = a[i]
       }
       // count conversions for team1
       k1 := 0
       for i < n && a[i] == d {
           k1++
           i++
       }
       // count conversions for team2
       k2 := 0
       for j < m && b[j] == d {
           k2++
           j++
       }
       // update scores: each converted throw changes from 3 to 2 => -1 per throw
       if k1 > 0 {
           score1 -= k1
       }
       if k2 > 0 {
           score2 -= k2
       }
       diff := score1 - score2
       // update best
       if diff > bestDiff || (diff == bestDiff && score1 > best1) {
           bestDiff = diff
           best1 = score1
           best2 = score2
       }
   }

   // output result
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fprintf(writer, "%d:%d", best1, best2)
}

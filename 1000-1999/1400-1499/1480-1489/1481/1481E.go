package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   // Track first and last occurrence and count for each color
   first := make([]int, n+1)
   last := make([]int, n+1)
   cnt := make([]int, n+1)
   const inf = 1<<60
   for i := 1; i <= n; i++ {
       first[i] = n + 1
   }
   for i, v := range a {
       cnt[v]++
       if first[v] == n+1 {
           first[v] = i
       }
       last[v] = i
   }
   // Collect colors present
   type colInfo struct{ first, last, cnt int }
   cols := make([]colInfo, 0, n)
   for v := 1; v <= n; v++ {
       if cnt[v] > 0 {
           cols = append(cols, colInfo{first[v], last[v], cnt[v]})
       }
   }
   sort.Slice(cols, func(i, j int) bool {
       return cols[i].first < cols[j].first
   })
   ans := 0
   segL, segR, segMax := -1, -1, 0
   for _, c := range cols {
       if segL == -1 {
           // start first segment
           segL, segR, segMax = c.first, c.last, c.cnt
       } else if c.first > segR {
           // finish previous segment
           segLen := segR - segL + 1
           ans += segLen - segMax
           // start new
           segL, segR, segMax = c.first, c.last, c.cnt
       } else {
           // extend current segment
           if c.last > segR {
               segR = c.last
           }
           if c.cnt > segMax {
               segMax = c.cnt
           }
       }
   }
   if segL != -1 {
       segLen := segR - segL + 1
       ans += segLen - segMax
   }
   fmt.Fprintln(writer, ans)
}

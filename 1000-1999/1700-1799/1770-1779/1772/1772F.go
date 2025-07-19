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

   var n, m, k int
   fmt.Fscan(reader, &n, &m, &k)

   // Read k+1 pictures of size n x m
   a := make([][][]byte, k+2)
   for i := range a {
       a[i] = make([][]byte, n)
       for j := 0; j < n; j++ {
           a[i][j] = make([]byte, m)
       }
   }
   for i := 1; i <= k+1; i++ {
       for j := 0; j < n; j++ {
           var s string
           fmt.Fscan(reader, &s)
           for l := 0; l < m; l++ {
               a[i][j][l] = s[l]
           }
       }
   }

   type pair struct{ x, idx int }
   v := make([]pair, 0, k+1)
   // Count valid recolor operations in each picture
   for i := 1; i <= k+1; i++ {
       cnt := 0
       for j := 1; j < n-1; j++ {
           for l := 1; l < m-1; l++ {
               if a[i][j-1][l] == a[i][j+1][l] &&
                  a[i][j][l+1] == a[i][j][l-1] &&
                  a[i][j-1][l] == a[i][j][l+1] &&
                  a[i][j][l] != a[i][j+1][l] {
                   cnt++
               }
           }
       }
       v = append(v, pair{cnt, i})
   }
   sort.Slice(v, func(i, j int) bool { return v[i].x < v[j].x })

   // Initial picture is the one with maximum possible moves
   best := v[len(v)-1].idx
   fmt.Fprintln(writer, best)

   moves := k
   current := best
   var ans []string
   // Reconstruct operations backwards
   for i := len(v) - 2; i >= 0; i-- {
       r := v[i].idx
       for j := 1; j < n-1; j++ {
           for l := 1; l < m-1; l++ {
               if a[r][j][l] != a[current][j][l] {
                   ans = append(ans, fmt.Sprintf("1 %d %d", j+1, l+1))
                   moves++
               }
           }
       }
       ans = append(ans, fmt.Sprintf("2 %d", r))
       current = r
   }

   fmt.Fprintln(writer, moves)
   for _, op := range ans {
       fmt.Fprintln(writer, op)
   }
}

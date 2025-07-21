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

   var n, m, k int
   fmt.Fscan(reader, &n, &m, &k)
   grid := make([]string, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &grid[i])
   }

   // Prefix sum of 'a' letters
   psa := make([][]int, n+1)
   for i := 0; i <= n; i++ {
       psa[i] = make([]int, m+1)
   }
   for i := 1; i <= n; i++ {
       for j := 1; j <= m; j++ {
           psa[i][j] = psa[i-1][j] + psa[i][j-1] - psa[i-1][j-1]
           if grid[i][j-1] == 'a' {
               psa[i][j]++
           }
       }
   }

   var ans int64
   // iterate over pairs of rows
   for x1 := 1; x1 < n; x1++ {
       for x2 := x1 + 1; x2 <= n; x2++ {
           // positions for each letter where corners match
           var pos [26][]int
           for y := 1; y <= m; y++ {
               c1 := grid[x1][y-1]
               c2 := grid[x2][y-1]
               if c1 == c2 {
                   pos[c1-'a'] = append(pos[c1-'a'], y)
               }
           }
           // for each letter, count valid column pairs
           for c := 0; c < 26; c++ {
               P := pos[c]
               t := len(P)
               if t < 2 {
                   continue
               }
               r := 1
               for l := 0; l < t-1; l++ {
                   if r < l+1 {
                       r = l + 1
                   }
                   for r < t {
                       y1 := P[l]
                       y2 := P[r]
                       // count 'a's in rectangle [x1..x2] x [y1..y2]
                       cnt := psa[x2][y2] - psa[x1-1][y2] - psa[x2][y1-1] + psa[x1-1][y1-1]
                       if cnt <= k {
                           r++
                       } else {
                           break
                       }
                   }
                   // valid pairs starting at l: indices (l, l+1)...(l, r-1)
                   valid := r - l - 1
                   if valid > 0 {
                       ans += int64(valid)
                   }
               }
           }
       }
   }
   fmt.Fprintln(writer, ans)
}

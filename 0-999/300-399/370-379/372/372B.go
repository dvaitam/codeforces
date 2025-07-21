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

   var n, m, q int
   if _, err := fmt.Fscan(reader, &n, &m, &q); err != nil {
       return
   }
   grid := make([][]byte, n+1)
   for i := 1; i <= n; i++ {
       var line string
       fmt.Fscan(reader, &line)
       grid[i] = []byte(" " + line) // 1-indexed
   }
   // prefix sum of ones
   var onesPS [41][41]int
   for i := 1; i <= n; i++ {
       for j := 1; j <= m; j++ {
           onesPS[i][j] = int(grid[i][j]-'0') + onesPS[i-1][j] + onesPS[i][j-1] - onesPS[i-1][j-1]
       }
   }
   // dp f[a][b][c][d]: number of all-zero subrectangles within [a..c]x[b..d]
   var f [41][41][41][41]uint32
   // helper p for current bottom-right
   var p [42][42]uint32
   // build dp
   for c := 1; c <= n; c++ {
       for d := 1; d <= m; d++ {
           // reset boundaries for p
           for i := 1; i <= c+1; i++ {
               p[i][d+1] = 0
           }
           for j := 1; j <= d+1; j++ {
               p[c+1][j] = 0
           }
           // compute p and f for this (c,d)
           for a := c; a >= 1; a-- {
               for b := d; b >= 1; b-- {
                   // is rectangle [a..c][b..d] all zeros?
                   ok := 0
                   if onesPS[c][d]-onesPS[a-1][d]-onesPS[c][b-1]+onesPS[a-1][b-1] == 0 {
                       ok = 1
                   }
                   p[a][b] = p[a+1][b] + p[a][b+1] - p[a+1][b+1] + uint32(ok)
                   // inclusion-exclusion on f
                   var sum uint32 = p[a][b]
                   if c > a {
                       sum += f[a][b][c-1][d]
                   }
                   if d > b {
                       sum += f[a][b][c][d-1]
                   }
                   if c > a && d > b {
                       sum -= f[a][b][c-1][d-1]
                   }
                   f[a][b][c][d] = sum
               }
           }
       }
   }
   // answer queries
   for i := 0; i < q; i++ {
       var a, b, c, d int
       fmt.Fscan(reader, &a, &b, &c, &d)
       fmt.Fprintln(writer, f[a][b][c][d])
   }
}

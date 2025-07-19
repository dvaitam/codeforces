package main

import (
   "bufio"
   "fmt"
   "os"
)

var reader = bufio.NewReader(os.Stdin)
var writer = bufio.NewWriter(os.Stdout)

// readInt reads next integer from stdin.
func readInt() int {
   sign := 1
   b, err := reader.ReadByte()
   for err == nil && (b < '0' || b > '9') && b != '-' {
       b, err = reader.ReadByte()
   }
   if err != nil {
       return 0
   }
   if b == '-' {
       sign = -1
       b, _ = reader.ReadByte()
   }
   n := 0
   for err == nil && b >= '0' && b <= '9' {
       n = n*10 + int(b-'0')
       b, err = reader.ReadByte()
   }
   return n * sign
}

func main() {
   defer writer.Flush()
   n := readInt()
   m := readInt()
   // a[i][j]: probability for row i, column j
   a := make([][]float64, n+1)
   for i := 0; i <= n; i++ {
       a[i] = make([]float64, m+1)
   }
   for i := 1; i <= n; i++ {
       for j := 1; j <= m; j++ {
           a[i][j] = float64(readInt()) / 1000.0
       }
   }
   // g[j][i][k]: dynamic programming for column j, row i, state k
   g := make([][][2]float64, m+1)
   for j := 0; j <= m; j++ {
       g[j] = make([][2]float64, n+1)
   }
   G := make([]bool, m+1)
   var Ans float64
   // initial DP build
   for j := 1; j <= m; j++ {
       idx := 0
       if G[j] {
           idx = 1
       }
       for i := 1; i <= n; i++ {
           other := idx ^ 1
           g[j][i][idx] = (g[j][i-1][other] + 1) * a[i][j] + g[j][i-1][idx] * (1 - a[i][j])
       }
   }
   // iterate n times to choose best column flips
   for t := 1; t <= n; t++ {
       Mx := 0.0
       p := 1
       for j := 1; j <= m; j++ {
           idx := 0
           if G[j] {
               idx = 1
           }
           diff := g[j][n][idx] - g[j][n][idx^1]
           if diff > Mx {
               Mx = diff
               p = j
           }
       }
       Ans += Mx
       // flip state for column p and rebuild its DP
       G[p] = !G[p]
       idx := 0
       if G[p] {
           idx = 1
       }
       other := idx ^ 1
       for i := 1; i <= n; i++ {
           g[p][i][idx] = (g[p][i-1][other] + 1) * a[i][p] + g[p][i-1][idx] * (1 - a[i][p])
       }
   }
   fmt.Fprintf(writer, "%.15f\n", Ans)
}

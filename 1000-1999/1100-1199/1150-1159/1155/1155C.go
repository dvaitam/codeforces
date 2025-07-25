package main

import (
   "bufio"
   "fmt"
   "os"
)

func gcd(a, b int64) int64 {
   for b != 0 {
       a, b = b, a%b
   }
   if a < 0 {
       return -a
   }
   return a
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   xs := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &xs[i])
   }
   ps := make([]int64, m)
   for j := 0; j < m; j++ {
       fmt.Fscan(reader, &ps[j])
   }
   // Compute gcd of differences x[i] - x[0]
   var d int64 = xs[1] - xs[0]
   for i := 2; i < n; i++ {
       d = gcd(d, xs[i]-xs[0])
   }
   // Find a p_j dividing d
   idx := -1
   var pj int64
   for j, p := range ps {
       if d%p == 0 {
           idx = j
           pj = p
           break
       }
   }
   if idx == -1 {
       fmt.Fprintln(writer, "NO")
       return
   }
   // Compute first ring minute y such that y ≡ xs[0] (mod pj), y >= 1
   r := xs[0] % pj
   var y int64
   if r == 0 {
       y = pj
   } else {
       y = r
   }
   fmt.Fprintln(writer, "YES")
   // j is 1-based index
   fmt.Fprintf(writer, "%d %d", y, idx+1)
}

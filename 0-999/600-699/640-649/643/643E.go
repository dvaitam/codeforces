package main

import (
   "bufio"
   "fmt"
   "os"
)

const H = 40

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var q int
   fmt.Fscan(reader, &q)
   n := q + 2
   dp := make([]float64, n*H)
   mul := make([]float64, n*H)
   p := make([]int, n)
   var path [H]int

   // initial state
   dp[0*H+0] = 1.0
   p[0] = -1
   for i := 0; i < n; i++ {
       for j := 0; j < H; j++ {
           mul[i*H+j] = 1.0
       }
   }

   cur := 1
   for i := 0; i < q; i++ {
       var typ, v int
       fmt.Fscan(reader, &typ, &v)
       v-- // zero-based index
       if typ == 1 {
           p[cur] = v
           dp[cur*H+0] = 1.0
           u := cur
           pcur := 0
           // build path up to H ancestors
           for h := 0; h < H && u != -1; h++ {
               path[pcur] = u
               pcur++
               u = p[u]
           }
           // adjust multipliers
           for k := 2; k < pcur; k++ {
               node := path[k]
               prev := path[k-1]
               mul[node*H+k] /= (1 - dp[prev*H+(k-1)]/2)
           }
           // update dp values
           for k := 1; k < pcur; k++ {
               node := path[k]
               prev := path[k-1]
               mul[node*H+k] *= (1 - dp[prev*H+(k-1)]/2)
               dp[node*H+k] = 1 - mul[node*H+k]
           }
           cur++
       } else if typ == 2 {
           sum := 0.0
           base := v * H
           for j := 1; j < H; j++ {
               sum += dp[base+j]
           }
           fmt.Fprintf(writer, "%.17f\n", sum)
       }
   }
}

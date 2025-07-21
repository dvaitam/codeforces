package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
)

func main() {
   var a, n int
   fmt.Scan(&a, &n)
   N := a + n - 1
   rem := make([]int32, N+1)
   for i := 1; i <= N; i++ {
       rem[i] = int32(i)
   }
   lim := int(math.Sqrt(float64(N)))
   for p := 2; p <= lim; p++ {
       sq := p * p
       for j := sq; j <= N; j += sq {
           for rem[j] % int32(sq) == 0 {
               rem[j] /= int32(sq)
           }
       }
   }
   var total uint64
   for x := a; x <= N; x++ {
       total += uint64(rem[x])
   }
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fprintln(writer, total)
}

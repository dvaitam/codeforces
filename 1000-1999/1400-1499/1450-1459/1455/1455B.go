package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var t int
   if _, err := fmt.Fscan(reader, &t); err != nil {
       return
   }
   for i := 0; i < t; i++ {
       var x int64
       fmt.Fscan(reader, &x)
       // find minimal N such that S = N*(N+1)/2 >= x and S-x != 1
       // start from approximate N
       // use float to estimate
       N := int64(math.Ceil((math.Sqrt(float64(8*x+1)) - 1) / 2))
       for {
           S := N * (N + 1) / 2
           if S >= x && S - x != 1 {
               fmt.Fprintln(writer, N)
               break
           }
           N++
       }
   }
}

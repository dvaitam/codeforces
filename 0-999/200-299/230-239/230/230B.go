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
   var n int
   fmt.Fscan(reader, &n)
   xs := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &xs[i])
   }
   const max = 1000000
   sieve := make([]bool, max+1)
   for i := 2; i <= max; i++ {
       sieve[i] = true
   }
   for i := 2; i*i <= max; i++ {
       if sieve[i] {
           for j := i * i; j <= max; j += i {
               sieve[j] = false
           }
       }
   }
   for _, x := range xs {
       if x < 4 {
           fmt.Fprintln(writer, "NO")
           continue
       }
       rtf := math.Sqrt(float64(x))
       rt := int64(rtf + 0.5)
       if rt*rt == x && rt <= max && sieve[int(rt)] {
           fmt.Fprintln(writer, "YES")
       } else {
           fmt.Fprintln(writer, "NO")
       }
   }
}

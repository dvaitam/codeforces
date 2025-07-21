package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var a, b, c int
   if _, err := fmt.Fscan(reader, &a, &b, &c); err != nil {
       return
   }
   N := a * b * c
   // Precompute number of divisors for 1..N
   d := make([]int, N+1)
   for i := 1; i <= N; i++ {
       for j := i; j <= N; j += i {
           d[j]++
       }
   }
   const mod = 1 << 30
   var sum int64
   // Sum d[i*j*k] for i=1..a, j=1..b, k=1..c
   for i := 1; i <= a; i++ {
       for j := 1; j <= b; j++ {
           ij := i * j
           for k := 1; k <= c; k++ {
               sum += int64(d[ij*k])
               if sum >= mod {
                   sum %= mod
               }
           }
       }
   }
   sum %= mod
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fprint(writer, sum)
}

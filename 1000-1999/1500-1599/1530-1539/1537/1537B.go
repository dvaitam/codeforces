package main

import (
   "bufio"
   "fmt"
   "os"
)

func abs64(x int64) int64 {
   if x < 0 {
       return -x
   }
   return x
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var t int
   fmt.Fscan(reader, &t)
   for ; t > 0; t-- {
       var n, m, i, j int64
       fmt.Fscan(reader, &n, &m, &i, &j)
       // distances to corners
       d11 := abs64(i-1) + abs64(j-1)
       dnn := abs64(i-n) + abs64(j-m)
       d1m := abs64(i-1) + abs64(j-m)
       dn1 := abs64(i-n) + abs64(j-1)
       // compare sums without common term (n-1 + m-1)
       if d11 + dnn >= d1m + dn1 {
           // use corners (1,1) and (n,m)
           fmt.Fprintf(writer, "1 1 %d %d\n", n, m)
       } else {
           // use corners (1,m) and (n,1)
           fmt.Fprintf(writer, "1 %d %d 1\n", m, n)
       }
   }
}

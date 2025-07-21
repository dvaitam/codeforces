package main

import (
   "bufio"
   "fmt"
   "math/bits"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   p := make([]int, n+1)
   for i := range p {
       p[i] = -1
   }
   // Build permutation maximizing sum of i xor p[i]
   for i := n; i >= 0; i-- {
       if p[i] != -1 {
           continue
       }
       if i == 0 {
           p[0] = 0
           break
       }
       k := bits.Len(uint(i)) - 1
       b := (1 << (k + 1)) - 1
       j := b - i
       if j >= 0 && j <= n && p[j] == -1 {
           p[i] = j
           p[j] = i
       } else {
           p[i] = i
       }
   }
   // Compute beauty
   var sum int64
   for i := 0; i <= n; i++ {
       sum += int64(i ^ p[i])
   }
   fmt.Fprintln(writer, sum)
   for i, v := range p {
       if i > 0 {
           writer.WriteByte(' ')
       }
       fmt.Fprint(writer, v)
   }
   fmt.Fprintln(writer)
}

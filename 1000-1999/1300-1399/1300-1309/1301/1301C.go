package main

import (
   "bufio"
   "fmt"
   "os"
)

// nextInt reads next integer from bufio.Reader
func nextInt(r *bufio.Reader) int {
   var c byte
   var err error
   for {
       c, err = r.ReadByte()
       if err != nil {
           return 0
       }
       if (c >= '0' && c <= '9') || c == '-' {
           break
       }
   }
   sign := 1
   if c == '-' {
       sign = -1
       c, _ = r.ReadByte()
   }
   x := 0
   for ; c >= '0' && c <= '9'; c, _ = r.ReadByte() {
       x = x*10 + int(c-'0')
   }
   return x * sign
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   t := nextInt(reader)
   for i := 0; i < t; i++ {
       n := int64(nextInt(reader))
       m := int64(nextInt(reader))
       // total substrings
       total := n * (n + 1) / 2
       // zeros count
       z := n - m
       // number of gaps where zeros can be: m+1
       k := m + 1
       // distribute zeros as evenly as possible
       a := z / k
       r := z % k
       // compute substrings consisting only of zeros
       // for blocks of size L: L*(L+1)/2
       // r blocks of size a+1, k-r blocks of size a
       zeroSub := (k-r)*(a*(a+1)/2) + r*((a+1)*(a+2)/2)
       // answer is total minus zero-only substrings
       ans := total - zeroSub
       fmt.Fprintln(writer, ans)
   }
}

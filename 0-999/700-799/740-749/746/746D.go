package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var N, K, A, B int
   if _, err := fmt.Fscan(reader, &N, &K, &A, &B); err != nil {
       return
   }
   // small holds char for smaller count, big for larger
   small, big := 'G', 'B'
   if A > B {
       A, B = B, A
       small, big = 'B', 'G'
   }
   // impossible if even distributing cannot avoid > K
   if (B-1)/K > A {
       fmt.Fprintln(writer, "NO")
       return
   }
   cnt := K
   // build result
   for B > A {
       if cnt == 0 {
           // place one of smaller
           writer.WriteByte(byte(small))
           cnt = K
           A--
       } else {
           writer.WriteByte(byte(big))
           cnt--
           B--
       }
   }
   // now B == A, alternate small+big
   for A > 0 {
       writer.WriteByte(byte(small))
       writer.WriteByte(byte(big))
       A--
   }
   writer.WriteByte('\n')
}

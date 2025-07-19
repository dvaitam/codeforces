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

   var t int
   if _, err := fmt.Fscan(reader, &t); err != nil {
       return
   }
   for ; t > 0; t-- {
       var n int
       var k int64
       fmt.Fscan(reader, &n, &k)
       if k%2 != 0 {
           fmt.Fprintln(writer, "No")
           continue
       }
       // initialize permutation 1..n
       p := make([]int, n)
       for i := 0; i < n; i++ {
           p[i] = i + 1
       }
       rem := k
       for i := 0; i < n-i-1; i++ {
           j := n - i - 1
           d := int64(j - i)
           if d*2 < rem {
               // full swap to ends
               p[i], p[j] = p[j], p[i]
               rem -= d * 2
           } else {
               // partial swap
               shift := int(rem / 2)
               p[i], p[i+shift] = p[i+shift], p[i]
               rem = 0
               break
           }
       }
       if rem == 0 {
           fmt.Fprintln(writer, "Yes")
           for i, v := range p {
               if i > 0 {
                   writer.WriteByte(' ')
               }
               fmt.Fprint(writer, v)
           }
           fmt.Fprintln(writer)
       } else {
           fmt.Fprintln(writer, "No")
       }
   }
}

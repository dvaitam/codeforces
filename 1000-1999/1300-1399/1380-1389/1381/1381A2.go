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
   fmt.Fscan(reader, &t)
   for t > 0 {
       t--
       var n int
       fmt.Fscan(reader, &n)
       var a, b string
       fmt.Fscan(reader, &a, &b)

       sa := []byte(a)
       sb := []byte(b)
       sa = append(sa, '0')
       sb = append(sb, '0')

       ops1 := make([]int, 0, 2*n)
       ops2 := make([]int, 0, 2*n)
       for i := 1; i <= n; i++ {
           if sa[i] != sa[i-1] {
               ops1 = append(ops1, i)
           }
           if sb[i] != sb[i-1] {
               ops2 = append(ops2, i)
           }
       }
       for i := len(ops2) - 1; i >= 0; i-- {
           ops1 = append(ops1, ops2[i])
       }

       k := len(ops1)
       fmt.Fprint(writer, k)
       for _, v := range ops1 {
           fmt.Fprint(writer, " ", v)
       }
       fmt.Fprintln(writer)
   }
}

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

   var n, m int
   fmt.Fscan(reader, &n, &m)
   var a, b []byte
   if n < m {
       for i := 0; i < n; i++ {
           a = append(a, '1')
           b = append(b, '1')
       }
       diff := 2*n - m
       if diff >= 0 {
           for i := 0; i < diff; i++ {
               b = append(b, '0')
           }
       } else {
           ones := m - 2*n
           pre := make([]byte, ones)
           for i := 0; i < ones; i++ {
               pre[i] = '1'
           }
           b = append(pre, b...)
       }
   } else {
       for i := 0; i < n; i++ {
           a = append(a, '1')
       }
       for i := 0; i < m-1; i++ {
           b = append(b, '9')
       }
       for i := 0; i < n-m; i++ {
           b = append(b, '8')
       }
       b = append(b, '9')
   }
   fmt.Fprintln(writer, string(b))
   fmt.Fprint(writer, string(a))
}

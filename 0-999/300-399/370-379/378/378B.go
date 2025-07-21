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

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   a := make([]int, n)
   b := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i], &b[i])
   }
   // auto qualifiers: top m in each semifinal
   m := n / 2
   resA := make([]byte, n)
   resB := make([]byte, n)
   for i := 0; i < n; i++ {
       resA[i] = '0'
       resB[i] = '0'
   }
   for i := 0; i < m; i++ {
       resA[i] = '1'
       resB[i] = '1'
   }
   // time qualifiers for k=0: best n times among all participants
   i, j := 0, 0
   cnt := 0
   for cnt < n && (i < n || j < n) {
       // pick smaller time
       if j == n || (i < n && a[i] < b[j]) {
           resA[i] = '1'
           i++
       } else {
           resB[j] = '1'
           j++
       }
       cnt++
   }
   // output
   writer.Write(resA)
   writer.WriteByte('\n')
   writer.Write(resB)
   writer.WriteByte('\n')
}

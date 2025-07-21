package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m, c int
   if _, err := fmt.Fscan(reader, &n, &m, &c); err != nil {
       return
   }
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   b := make([]int, m)
   for i := 0; i < m; i++ {
       fmt.Fscan(reader, &b[i])
   }
   prefixB := make([]int, m+1)
   for i := 0; i < m; i++ {
       prefixB[i+1] = prefixB[i] + b[i]
   }
   diff := n - m
   for k := 0; k < n; k++ {
       L := k - diff
       if L < 0 {
           L = 0
       }
       R := k
       if R >= m {
           R = m - 1
       }
       sum := prefixB[R+1] - prefixB[L]
       sum %= c
       if sum < 0 {
           sum += c
       }
       a[k] = (a[k] + sum) % c
   }
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   for i := 0; i < n; i++ {
       if i > 0 {
           writer.WriteByte(' ')
       }
       fmt.Fprint(writer, a[i])
   }
   writer.WriteByte('\n')
}

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
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }

   // p[i]: furthest index j>=i such that a[i..j] is non-decreasing
   p := make([]int, n)
   p[n-1] = n - 1
   for i := n - 2; i >= 0; i-- {
       if a[i] <= a[i+1] {
           p[i] = p[i+1]
       } else {
           p[i] = i
       }
   }
   // q[i]: furthest index j>=i such that a[i..j] is non-increasing
   q := make([]int, n)
   q[n-1] = n - 1
   for i := n - 2; i >= 0; i-- {
       if a[i] >= a[i+1] {
           q[i] = q[i+1]
       } else {
           q[i] = i
       }
   }

   for k := 0; k < m; k++ {
       var l, r int
       fmt.Fscan(reader, &l, &r)
       l--
       r--
       x := p[l]
       if x >= r || q[x] >= r {
           fmt.Fprintln(writer, "Yes")
       } else {
           fmt.Fprintln(writer, "No")
       }
   }
}

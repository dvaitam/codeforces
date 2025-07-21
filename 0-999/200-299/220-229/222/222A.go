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

   var n, k int
   if _, err := fmt.Fscan(reader, &n, &k); err != nil {
       return
   }
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   // target value x is initial k-th element
   x := a[k-1]
   // length of first cycle in c
   L := n - k + 1
   // find u: first position in c[1..L] where c[u] != x; c[i]=a[k+i-2]
   u := L + 1
   for i := 1; i <= L; i++ {
       if a[k+i-2] != x {
           u = i
           break
       }
   }
   // find j: last index in a where a[j] != x
   j := 0
   for i := n; i >= 1; i-- {
       if a[i-1] != x {
           j = i
           break
       }
   }
   t := j
   // try t in [0..n-1]
   if t == 0 {
       fmt.Fprintln(writer, 0)
       return
   }
   if t <= n-1 {
       if t < L {
           if t <= u-1 {
               fmt.Fprintln(writer, t)
               return
           }
       } else {
           if u > L {
               fmt.Fprintln(writer, t)
               return
           }
       }
   }
   // if unable and infinite c-prefix constant, t>=n works
   if u > L {
       fmt.Fprintln(writer, n)
   } else {
       fmt.Fprintln(writer, -1)
   }
}

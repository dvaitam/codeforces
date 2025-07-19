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
   fmt.Fscan(reader, &n, &k)
   a := make([]int, n+1)
   for i := 1; i <= n; i++ {
       a[i] = 1
   }
   a[n] = 0
   mx := 1
   for step := 0; step < k; step++ {
       for i := 1; i <= n; i++ {
           if a[i] == n-i {
               fmt.Fprint(writer, n, " ")
           } else if a[i]+mx >= n-i {
               fmt.Fprint(writer, i+a[i], " ")
               a[i] = n - i
           } else {
               fmt.Fprint(writer, 1, " ")
               a[i] += mx
           }
       }
       fmt.Fprintln(writer)
       mx *= 2
   }
}

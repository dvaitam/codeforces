package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
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
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   sort.Ints(a)
   // Check triangle condition for largest three
   if n < 3 || a[n-3]+a[n-2] <= a[n-1] {
       fmt.Fprintln(writer, "NO")
       return
   }
   b := make([]int, n)
   half := n / 2
   for k := 0; k < half; k++ {
       b[k] = a[2*k]
       b[n-1-k] = a[2*k+1]
   }
   if n%2 == 1 {
       b[half] = a[n-1]
   }
   fmt.Fprintln(writer, "YES")
   for i, v := range b {
       if i > 0 {
           writer.WriteByte(' ')
       }
       fmt.Fprint(writer, v)
   }
   writer.WriteByte('\n')
}

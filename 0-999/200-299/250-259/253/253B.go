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
   maxKeep := 0
   j := 0
   for i := 0; i < n; i++ {
       // advance j while within 2 * a[i]
       for j < n && a[j] <= 2*a[i] {
           j++
       }
       // window [i, j-1]
       if j-i > maxKeep {
           maxKeep = j - i
       }
   }
   // minimum to remove = total - maximum kept
   res := n - maxKeep
   fmt.Fprintln(writer, res)
}

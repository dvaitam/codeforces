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

   var t int
   if _, err := fmt.Fscan(reader, &t); err != nil {
       return
   }
   for ti := 0; ti < t; ti++ {
       var n int
       fmt.Fscan(reader, &n)
       total := 2 * n
       a := make([]int, total)
       for i := 0; i < total; i++ {
           fmt.Fscan(reader, &a[i])
       }
       sort.Ints(a)
       // minimum difference between median of two odd-sized classes
       // is difference between middle two elements
       diff := a[n] - a[n-1]
       fmt.Fprintln(writer, diff)
   }
}

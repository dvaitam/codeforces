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
   for ; t > 0; t-- {
       var n int
       fmt.Fscan(reader, &n)
       a := make([]int, n)
       for i := 0; i < n; i++ {
           fmt.Fscan(reader, &a[i])
       }
       sort.Ints(a)
       // pair smallest element with next n/2 elements
       smallest := a[0]
       half := n / 2
       for i := 1; i <= half; i++ {
           fmt.Fprintf(writer, "%d %d\n", a[i], smallest)
       }
   }
}

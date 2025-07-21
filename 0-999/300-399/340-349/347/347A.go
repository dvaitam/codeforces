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
   fmt.Fscan(reader, &n)
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   // Sort to easily pick min, max, and lexicographically smallest middle
   sort.Ints(a)
   if n == 0 {
       return
   }
   // Output max at front
   fmt.Fprint(writer, a[n-1])
   // Output sorted middle elements
   for i := 1; i < n-1; i++ {
       fmt.Fprint(writer, " ", a[i])
   }
   // Output min at end if exists
   if n > 1 {
       fmt.Fprint(writer, " ", a[0])
   }
   fmt.Fprint(writer, "\n")
}

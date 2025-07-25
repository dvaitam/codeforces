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
   oddCount := 0
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
       if a[i]&1 != 0 {
           oddCount++
       }
   }
   // If both even and odd present, we can sort arbitrarily
   if oddCount > 0 && oddCount < n {
       sort.Ints(a)
   }
   // Output result
   for i, v := range a {
       if i > 0 {
           writer.WriteByte(' ')
       }
       fmt.Fprint(writer, v)
   }
   writer.WriteByte('\n')
}

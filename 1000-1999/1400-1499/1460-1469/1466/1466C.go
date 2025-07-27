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

   var t int
   if _, err := fmt.Fscan(reader, &t); err != nil {
       return
   }
   for i := 0; i < t; i++ {
       var s string
       fmt.Fscan(reader, &s)
       n := len(s)
       // convert to byte slice for modifications
       a := []byte(s)
       var changes int
       for j := 0; j < n; j++ {
           if j > 0 && a[j] == a[j-1] || j > 1 && a[j] == a[j-2] {
               changes++
               // mark as changed to a sentinel not equal to any lowercase letter
               a[j] = '?' 
           }
       }
       fmt.Fprintln(writer, changes)
   }
}

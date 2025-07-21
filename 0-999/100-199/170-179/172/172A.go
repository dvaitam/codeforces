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

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   var first string
   fmt.Fscan(reader, &first)
   commonLen := len(first)

   for i := 1; i < n; i++ {
       var s string
       fmt.Fscan(reader, &s)
       // compute common prefix length with first
       j := 0
       // s and first have same length
       for j < commonLen && j < len(s) && s[j] == first[j] {
           j++
       }
       commonLen = j
       if commonLen == 0 {
           break
       }
   }
   fmt.Fprintln(writer, commonLen)
}

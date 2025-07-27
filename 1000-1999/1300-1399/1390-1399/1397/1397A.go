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
   for ; t > 0; t-- {
       var n int
       fmt.Fscan(reader, &n)
       freq := make([]int, 26)
       for i := 0; i < n; i++ {
           var s string
           fmt.Fscan(reader, &s)
           for _, ch := range s {
               freq[ch-'a']++
           }
       }
       ok := true
       for _, cnt := range freq {
           if cnt % n != 0 {
               ok = false
               break
           }
       }
       if ok {
           fmt.Fprintln(writer, "YES")
       } else {
           fmt.Fprintln(writer, "NO")
       }
   }
}

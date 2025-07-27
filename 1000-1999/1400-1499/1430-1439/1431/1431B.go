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
       var s string
       fmt.Fscan(reader, &s)
       var ans int
       // Count each 'w' which must be underlined
       for i := 0; i < len(s); i++ {
           if s[i] == 'w' {
               ans++
           }
       }
       // Count runs of consecutive 'v'
       for i := 0; i < len(s); {
           if s[i] != 'v' {
               i++
               continue
           }
           // start of run
           j := i
           for j < len(s) && s[j] == 'v' {
               j++
           }
           // run length = j - i
           ans += (j - i) / 2
           i = j
       }
       fmt.Fprintln(writer, ans)
   }
}

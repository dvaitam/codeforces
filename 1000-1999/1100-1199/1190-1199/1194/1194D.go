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

   var T int
   if _, err := fmt.Fscan(reader, &T); err != nil {
       return
   }
   for i := 0; i < T; i++ {
       var n, k int64
       fmt.Fscan(reader, &n, &k)
       if k%3 != 0 {
           if n%3 == 0 {
               writer.WriteString("Bob\n")
           } else {
               writer.WriteString("Alice\n")
           }
       } else {
           x := n % (k + 1)
           if x%3 == 0 && x != k {
               writer.WriteString("Bob\n")
           } else {
               writer.WriteString("Alice\n")
           }
       }
   }
}

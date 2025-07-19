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

   for {
       var n int
       if _, err := fmt.Fscan(reader, &n); err != nil {
           return
       }
       // Use stack to pair adjacent equal parities
       stack := make([]int, 0, n)
       for i := 0; i < n; i++ {
           var m int
           fmt.Fscan(reader, &m)
           m %= 2
           l := len(stack)
           if l > 0 && stack[l-1] == m {
               stack = stack[:l-1]
           } else {
               stack = append(stack, m)
           }
       }
       if len(stack) > 1 {
           writer.WriteString("NO\n")
       } else {
           writer.WriteString("YES\n")
       }
   }
}

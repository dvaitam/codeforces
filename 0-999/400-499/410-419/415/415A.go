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

   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   buttons := make([]int, m)
   for i := 0; i < m; i++ {
       fmt.Fscan(reader, &buttons[i])
   }
   result := make([]int, n+1)
   for _, btn := range buttons {
       for i := btn; i <= n; i++ {
           if result[i] == 0 {
               result[i] = btn
           }
       }
   }
   for i := 1; i <= n; i++ {
       fmt.Fprint(writer, result[i])
       if i < n {
           fmt.Fprint(writer, " ")
       }
   }
   fmt.Fprintln(writer)
}

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
   // find maximum k such that k*(k-1)/2 <= n
   k := 1
   for (k*(k-1))/2 <= n {
       k++
   }
   k--
   // prepare guest lists for days 1..k
   guests := make([][]int, k)
   id := 1
   for i := 0; i < k; i++ {
       for j := i + 1; j < k; j++ {
           guests[i] = append(guests[i], id)
           guests[j] = append(guests[j], id)
           id++
       }
   }
   // output
   fmt.Fprintln(writer, k)
   for i := 0; i < k; i++ {
       // each guest list for day i+1
       for j, x := range guests[i] {
           if j > 0 {
               writer.WriteByte(' ')
           }
           fmt.Fprint(writer, x)
       }
       fmt.Fprintln(writer)
   }
}

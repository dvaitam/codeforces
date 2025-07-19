package main

import (
   "bufio"
   "fmt"
   "os"
)

func abs(x int) int {
   if x < 0 {
       return -x
   }
   return x
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   if n == 1 || n == 2 {
       fmt.Fprintln(writer, 1)
       fmt.Fprintln(writer, "1 1")
       return
   }
   // adjust n to make n%3 == 2
   if n%3 == 1 {
       n++
   }
   orig := n
   if n%3 == 0 {
       n--
   }
   // now n%3 == 2
   m := n
   // collect positions
   type pair struct{ a, b int }
   var res []pair
   third := m / 3
   for i := 0; i < m; i++ {
       j1 := third - i
       if j1 >= 0 && j1 < m {
           res = append(res, pair{i + 1, j1 + 1})
       }
       j2 := m - 1 - i
       if abs(i-j2) <= third {
           res = append(res, pair{i + 1, j2 + 1})
       }
   }
   if orig%3 == 0 {
       res = append(res, pair{orig, orig})
   }
   // output
   fmt.Fprintln(writer, len(res))
   for _, p := range res {
       fmt.Fprintf(writer, "%d %d\n", p.a, p.b)
   }
}

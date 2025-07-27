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
   fmt.Fscan(reader, &n)
   a := make([]int, n+1)
   // helper to ask a query
   query := func(op string, i, j int) int {
       fmt.Fprintf(writer, "%s %d %d\n", op, i, j)
       writer.Flush()
       var res int
       fmt.Fscan(reader, &res)
       return res
   }
   // get pairwise AND and OR for first three elements
   and12 := query("AND", 1, 2)
   or12 := query("OR", 1, 2)
   and13 := query("AND", 1, 3)
   or13 := query("OR", 1, 3)
   and23 := query("AND", 2, 3)
   or23 := query("OR", 2, 3)
   // sums of pairs
   sum12 := and12 + or12
   sum13 := and13 + or13
   sum23 := and23 + or23
   // reconstruct first three values
   a1 := (sum12 + sum13 - sum23) / 2
   a[1] = a1
   a[2] = sum12 - a1
   a[3] = sum13 - a1
   // reconstruct the rest using XOR with first element
   for i := 4; i <= n; i++ {
       x := query("XOR", 1, i)
       a[i] = a1 ^ x
   }
   // output the result
   fmt.Fprintf(writer, "!")
   for i := 1; i <= n; i++ {
       fmt.Fprintf(writer, " %d", a[i])
   }
   fmt.Fprintln(writer)
}

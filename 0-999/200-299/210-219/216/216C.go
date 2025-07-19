package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m, k int
   if _, err := fmt.Fscan(reader, &n, &m, &k); err != nil {
       return
   }
   // Special case
   if n == 2 && m == 2 && k == 1 {
       fmt.Println(4)
       fmt.Println("1 2 3 4")
       return
   }
   // Determine if an extra segment can be added
   q := 0
   if n+n == n+m+1 && k == 1 {
       q = 1
   }
   if n+n < n+m+1 {
       q = 1
   }
   total := k*2 + q
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fprintln(writer, total)
   // Print k times 1
   for i := 0; i < k; i++ {
       fmt.Fprint(writer, "1 ")
   }
   // Print n
   fmt.Fprintf(writer, "%d ", n)
   // Print (k-1) times n+1
   for i := 0; i < k-1; i++ {
       fmt.Fprintf(writer, "%d ", n+1)
   }
   // Adjustment for equal n and m in single segment case
   a := 0
   if n == m && k == 1 {
       a = 1
   }
   // Print final segment if applicable
   if q == 1 {
       fmt.Fprintf(writer, "%d", n+m-a)
   }
}

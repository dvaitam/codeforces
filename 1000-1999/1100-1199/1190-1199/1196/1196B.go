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

   var q int
   if _, err := fmt.Fscan(reader, &q); err != nil {
       return
   }
   for qi := 0; qi < q; qi++ {
       var n, k int
       fmt.Fscan(reader, &n, &k)
       oddPos := make([]int, 0, n)
       for i := 1; i <= n; i++ {
           var x int
           fmt.Fscan(reader, &x)
           if x&1 == 1 {
               oddPos = append(oddPos, i)
           }
       }
       // need at least k odd and total odd parity == k parity
       if len(oddPos) < k || (len(oddPos)%2) != (k%2) {
           fmt.Fprintln(writer, "NO")
           continue
       }
       fmt.Fprintln(writer, "YES")
       // output first k-1 segment ends at odd positions
       for i := 0; i < k-1; i++ {
           fmt.Fprintf(writer, "%d ", oddPos[i])
       }
       // last segment ends at n
       fmt.Fprintf(writer, "%d\n", n)
   }
}

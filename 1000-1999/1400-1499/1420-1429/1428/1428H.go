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
   k := n * m
   // Perform a full sweep rotating ring 0 clockwise
   f := make([]int, k)
   for i := 0; i < k; i++ {
       // Query type 1: rotate ring 1 (0-indexed ring 0) clockwise by 1
       fmt.Fprintf(writer, "1 1 1\n")
       writer.Flush()
       // Read the number of unblocked lasers
       var x int
       if _, err := fmt.Fscan(reader, &x); err != nil {
           return
       }
       f[i] = x
   }
   // Find the maximum value (should correspond to alignments)
   maxVal := 0
   for _, v := range f {
       if v > maxVal {
           maxVal = v
       }
   }
   // Collect up to n-1 positions with f == maxVal
   res := make([]int, 0, n-1)
   for i, v := range f {
       if v == maxVal {
           // position is i+1 rotations
           res = append(res, i+1)
           if len(res) == n-1 {
               break
           }
       }
   }
   // Output the result
   fmt.Fprintf(writer, "2")
   for _, p := range res {
       fmt.Fprintf(writer, " %d", p)
   }
   fmt.Fprintf(writer, "\n")
}

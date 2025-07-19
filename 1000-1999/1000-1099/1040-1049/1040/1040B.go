package main

import (
   "bufio"
   "fmt"
   "os"
)

func min(a, b int) int {
   if a < b {
       return a
   }
   return b
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, k int
   if _, err := fmt.Fscan(reader, &n, &k); err != nil {
       return
   }
   // If small n, only one station
   if n < (k+1)*2 {
       pos := min(n, k+1)
       fmt.Fprintln(writer, 1)
       fmt.Fprintln(writer, pos)
       return
   }
   // Try to find optimal starting offset
   step := 2*k + 1
   baseMin := (k + 1) * 2
   found := false
   for i := baseMin; i <= baseMin+step; i++ {
       if (n - i) % step == 0 {
           printPositions(writer, n, k, i)
           found = true
           break
       }
   }
   if !found {
       // fallback: every position
       fmt.Fprintln(writer, n)
       for i := 1; i <= n; i++ {
           fmt.Fprint(writer, i)
           if i < n {
               fmt.Fprint(writer, " ")
           }
       }
       fmt.Fprintln(writer)
   }
}

// printPositions computes and prints positions given initial segment length i
func printPositions(writer *bufio.Writer, n, k, seg int) {
   var positions []int
   // compute first position
   // seg includes two edge ranges of size k+1
   rem := seg - (k+1)*2
   start := 1 + min(rem, k)
   positions = append(positions, start)
   // subsequent positions
   step := 2*k + 1
   for p := start + step; p <= n; p += step {
       positions = append(positions, p)
   }
   // output
   fmt.Fprintln(writer, len(positions))
   for i, v := range positions {
       if i > 0 {
           fmt.Fprint(writer, " ")
       }
       fmt.Fprint(writer, v)
   }
   fmt.Fprintln(writer)
}

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

   var t int
   if _, err := fmt.Fscan(reader, &t); err != nil {
       return
   }
   for ; t > 0; t-- {
       solve(reader, writer)
   }
}

func solve(reader *bufio.Reader, writer *bufio.Writer) {
   var n int
   fmt.Fscan(reader, &n)
   a := make([]int64, n)
   b := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   minDiff := int64(1e18)
   maxSum := int64(-1e18)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &b[i])
       diff := a[i] - b[i]
       sum := a[i] + b[i]
       if diff < minDiff {
           minDiff = diff
       }
       if sum > maxSum {
           maxSum = sum
       }
   }
   // result is (minDiff + maxSum) / 2.0
   res := float64(minDiff+maxSum) * 0.5
   fmt.Fprintf(writer, "%.6f\n", res)
}

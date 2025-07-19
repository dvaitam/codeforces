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
   half := n / 2
   b := make([]int64, half+2)
   for i := 1; i <= half; i++ {
       fmt.Fscan(reader, &b[i])
   }

   a := make([]int64, n+2)
   // Initialize first and last elements
   a[1] = 0
   a[n] = b[1]
   // Construct the array to minimize lexicographically
   for i := 2; i <= half; i++ {
       j := n - i + 1
       a[i] = a[i-1]
       a[j] = b[i] - a[i]
       if a[j] > a[j+1] {
           diff := a[j] - a[j+1]
           a[i] += diff
           a[j] = a[j+1]
       }
   }

   // Output result
   for i := 1; i <= n; i++ {
       if i > 1 {
           writer.WriteByte(' ')
       }
       fmt.Fprintf(writer, "%d", a[i])
   }
   writer.WriteByte('\n')
}

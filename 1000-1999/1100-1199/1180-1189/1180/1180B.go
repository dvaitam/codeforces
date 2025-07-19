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
   a := make([]int, n)
   // Read and initial transform
   allZero := true
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
       if a[i] >= 0 {
           a[i] = ^a[i]
       }
       // Check if ~(a[i]) == 0
       if ^a[i] != 0 {
           allZero = false
       }
   }
   if allZero {
       // All elements were zero or -1 originally
       if n%2 == 1 {
           for i := 0; i < n; i++ {
               a[i] = 0
           }
       }
   } else {
       if n%2 == 1 {
           // Find minimal a[i]
           minVal := a[0]
           minIdx := 0
           for i := 1; i < n; i++ {
               if a[i] < minVal {
                   minVal = a[i]
                   minIdx = i
               }
           }
           a[minIdx] = ^a[minIdx]
       }
   }
   // Output
   if n > 0 {
       fmt.Fprint(writer, a[0])
       for i := 1; i < n; i++ {
           // print space and value
           fmt.Fprint(writer, " ", a[i])
       }
   }
}

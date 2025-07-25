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

   var n, k int
   fmt.Fscan(reader, &n, &k)
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   total := 0
   for _, v := range a {
       total += v
   }
   if total%k != 0 {
       fmt.Fprintln(writer, "No")
       return
   }
   target := total / k
   parts := make([]int, 0, k)
   currSum := 0
   currCount := 0
   for _, v := range a {
       currSum += v
       currCount++
       if currSum == target {
           parts = append(parts, currCount)
           currSum = 0
           currCount = 0
       } else if currSum > target {
           fmt.Fprintln(writer, "No")
           return
       }
   }
   if currSum != 0 || len(parts) != k {
       fmt.Fprintln(writer, "No")
       return
   }
   fmt.Fprintln(writer, "Yes")
   for i, v := range parts {
       if i > 0 {
           fmt.Fprint(writer, " ")
       }
       fmt.Fprint(writer, v)
   }
   fmt.Fprintln(writer)
}

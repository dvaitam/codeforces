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

   var k int
   if _, err := fmt.Fscan(reader, &k); err != nil {
       return
   }
   n := k + 1
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   ambIdx := -1
   for i := 1; i < n; i++ {
       if a[i-1] > 1 && a[i] > 1 {
           ambIdx = i
           break
       }
   }
   if ambIdx == -1 {
       fmt.Fprintln(writer, "perfect")
       return
   }
   fmt.Fprintln(writer, "ambiguous")

   // first parent assignment (perfect tree)
   sum := 0
   for i := 0; i < n; i++ {
       for j := 0; j < a[i]; j++ {
           fmt.Fprint(writer, sum, " ")
       }
       sum += a[i]
   }
   fmt.Fprintln(writer)

   // second assignment (introduce ambiguity)
   sum = 0
   for i := 0; i < n; i++ {
       for j := 0; j < a[i]; j++ {
           if i == ambIdx && j == 0 {
               fmt.Fprint(writer, sum-1, " ")
           } else {
               fmt.Fprint(writer, sum, " ")
           }
       }
       sum += a[i]
   }
   fmt.Fprintln(writer)

package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   day := 0
   p := 0
   for p < n {
       day++
       goal := a[p]
       // expand segment until reaching goal
       for p < goal-1 {
           p++
           if a[p] > goal {
               goal = a[p]
           }
       }
       p++
   }
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fprintln(writer, day)
}

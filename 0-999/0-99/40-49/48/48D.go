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
   maxVal := 0
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
       if a[i] > maxVal {
           maxVal = a[i]
       }
   }
   // count frequencies
   b := make([]int, maxVal+1)
   for _, v := range a {
       if v >= 0 && v < len(b) {
           b[v]++
       }
   }
   // check non-increasing frequency
   ok := true
   for i := 2; i < len(b); i++ {
       if b[i] > b[i-1] {
           ok = false
           break
       }
   }
   if !ok {
       fmt.Fprintln(writer, -1)
       return
   }
   // print number of groups = count of 1s
   if len(b) > 1 {
       fmt.Fprintln(writer, b[1])
   } else {
       fmt.Fprintln(writer, 0)
   }
   // output assignment for each element in order
   for _, v := range a {
       // print current bucket number
       fmt.Fprint(writer, b[v], " ")
       // decrement
       b[v]--
   }
   fmt.Fprintln(writer)
}

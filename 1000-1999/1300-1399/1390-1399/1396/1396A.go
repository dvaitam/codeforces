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
   a := make([]int64, n)
   for i := range a {
       fmt.Fscan(reader, &a[i])
   }

   if n == 1 {
       // Three operations to zero the single element
       for z := 0; z < 3; z++ {
           fmt.Fprintln(writer, "1 1")
           fmt.Fprintln(writer, -a[0])
           a[0] = 0
       }
       return
   }

   // Operation 1: Zero the first element
   fmt.Fprintln(writer, "1 1")
   fmt.Fprintln(writer, -a[0])
   a[0] = 0

   // Operation 2: For all i in [1..n], subtract a[i]*n
   fmt.Fprintf(writer, "1 %d\n", n)
   for i := 0; i < n; i++ {
       fmt.Fprintf(writer, "%d ", -a[i]*int64(n))
   }
   fmt.Fprintln(writer)

   // Operation 3: For all i in [2..n], add a[i]*(n-1)
   fmt.Fprintf(writer, "2 %d\n", n)
   for i := 1; i < n; i++ {
       fmt.Fprintf(writer, "%d ", a[i]*int64(n-1))
   }
   fmt.Fprintln(writer)
}

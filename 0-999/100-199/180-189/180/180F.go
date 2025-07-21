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
   a := make([]int, n+1)
   b := make([]int, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &b[i])
   }
   invA := make([]int, n+1)
   for i := 1; i <= n; i++ {
       invA[a[i]] = i
   }
   p := make([]int, n+1)
   for i := 1; i <= n; i++ {
       // p[i] = b[ invA[i] ]
       p[i] = b[invA[i]]
   }
   for i := 1; i <= n; i++ {
       if i > 1 {
           writer.WriteByte(' ')
       }
       fmt.Fprintf(writer, "%d", p[i])
   }
   writer.WriteByte('\n')
}

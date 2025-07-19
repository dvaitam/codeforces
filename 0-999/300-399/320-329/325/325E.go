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
   // n must be even
   if n&1 == 1 {
       fmt.Fprintln(writer, -1)
       return
   }
   a := make([]int, n+1)
   vis := make([]bool, n)
   // initialize last two
   a[n] = 0
   a[n-1] = n / 2
   vis[n/2] = true
   // build sequence backwards
   for i := n - 1; i > 0; i-- {
       x1 := a[i] / 2
       x2 := (a[i] + n) / 2
       if vis[x1] {
           a[i-1] = x2
           vis[x2] = true
       } else if vis[x2] {
           a[i-1] = x1
           vis[x1] = true
       } else if x1 > x2 {
           a[i-1] = x1
           vis[x1] = true
       } else {
           a[i-1] = x2
           vis[x2] = true
       }
   }
   // output sequence
   for i := 0; i <= n; i++ {
       fmt.Fprintf(writer, "%d", a[i])
       if i < n {
           writer.WriteByte(' ')
       }
   }
   writer.WriteByte('\n')
}

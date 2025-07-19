package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, k int
   if _, err := fmt.Fscan(reader, &n, &k); err != nil {
       return
   }
   L := make([]int, n+1)
   R := make([]int, n+1)
   p := make([]int, n+1)
   vis := make([]bool, n+1)
   a0 := make([]int, n+1)
   a1 := make([]int, n+1)
   // prepare intervals
   for i := 1; i <= n; i += k - 1 {
       end := i + k - 2
       if end > n {
           end = n
       }
       for j := i; j <= end; j++ {
           L[j] = i
           R[j] = end
       }
   }
   // read values and assign
   total := n * k
   for i := 1; i <= total; i++ {
       var x int
       fmt.Fscan(reader, &x)
       if vis[x] {
           continue
       }
       if p[x] != 0 {
           a0[x] = p[x]
           a1[x] = i
           vis[x] = true
           for j := L[x]; j <= R[x]; j++ {
               p[j] = 0
           }
       } else {
           p[x] = i
       }
   }
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   for i := 1; i <= n; i++ {
       fmt.Fprintln(writer, a0[i], a1[i])
   }
}

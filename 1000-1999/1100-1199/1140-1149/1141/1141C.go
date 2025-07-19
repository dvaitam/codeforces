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
   v := make([]int, n-1)
   var sum int64
   var pos int
   for i := 0; i < n-1; i++ {
       fmt.Fscan(reader, &v[i])
       sum += int64(v[i])
       if sum > 0 {
           pos++
       }
   }
   // initial element
   first := n - pos
   ans := make([]int, 0, n)
   ans = append(ans, first)
   vis := make([]bool, n+1)
   if first < 1 || first > n {
       fmt.Fprint(writer, -1)
       return
   }
   vis[first] = true
   // build permutation
   for i := 1; i < n; i++ {
       next := ans[i-1] + v[i-1]
       if next < 1 || next > n || vis[next] {
           fmt.Fprint(writer, -1)
           return
       }
       ans = append(ans, next)
       vis[next] = true
   }
   // output
   for i, x := range ans {
       if i > 0 {
           writer.WriteByte(' ')
       }
       fmt.Fprint(writer, x)
   }
}

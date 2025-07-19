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
   f := make([]int, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &f[i])
   }
   for i := 1; i <= n; i++ {
       if f[f[i]] != f[i] {
           writer.WriteString("-1")
           return
       }
   }
   p := make(map[int]int)
   g := make([]int, n+1)
   h := make([]int, 0, n)
   s := 0
   for i := 1; i <= n; i++ {
       fi := f[i]
       idx, ok := p[fi]
       if !ok {
           s++
           idx = s
           p[fi] = s
           h = append(h, fi)
       }
       g[i] = idx
   }
   fmt.Fprintln(writer, s)
   for i := 1; i <= n; i++ {
       if i > 1 {
           writer.WriteByte(' ')
       }
       fmt.Fprint(writer, g[i])
   }
   writer.WriteByte('\n')
   for i := 0; i < s; i++ {
       if i > 0 {
           writer.WriteByte(' ')
       }
       fmt.Fprint(writer, h[i])
   }
   writer.WriteByte('\n')
}

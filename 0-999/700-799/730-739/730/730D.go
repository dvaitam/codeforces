package main

import (
   "bufio"
   "fmt"
   "io"
   "os"
)

func minInt64(a, b int64) int64 {
   if a < b {
       return a
   }
   return b
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   var r int64
   if _, err := fmt.Fscan(reader, &n, &r); err == io.EOF {
       return
   } else if err != nil {
       panic(err)
   }
   l := make([]int64, n+1)
   t := make([]int64, n+1)
   sul := make([]int64, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &l[i])
       sul[i] = sul[i-1] + l[i]
   }
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &t[i])
   }
   need := make([]int64, n+1)
   for i := 1; i <= n; i++ {
       need[i] = 2*l[i] - t[i]
       if need[i] < 0 {
           need[i] = 0
       }
       if need[i] > l[i] {
           fmt.Fprintln(writer, -1)
           return
       }
   }
   al := make([]int64, n+2)
   var ans int64
   pos := make([]int64, 0, 100000)
   for i := 1; i <= n; i++ {
       if al[i] >= need[i] {
           continue
       }
       still := need[i] - al[i]
       cnt := (still + r - 1) / r
       duo := cnt*r - still
       // generate positions
       start := sul[i] - still
       for j := start; j < sul[i]; j += r {
           if len(pos) >= 100000 {
               break
           }
           pos = append(pos, j)
       }
       ans += cnt
       // distribute leftover 'duo' to future al
       for j := i + 1; j <= n && duo > 0; j++ {
           del := minInt64(duo, l[j])
           duo -= del
           al[j] = del
       }
   }
   // output answer
   fmt.Fprintln(writer, ans)
   if ans <= 100000 {
       // print adjusted positions
       for i, v := range pos {
           // formula: 2*v - r*i
           x := 2*v - r*int64(i)
           if i > 0 {
               writer.WriteByte(' ')
           }
           fmt.Fprint(writer, x)
       }
       if len(pos) > 0 {
           writer.WriteByte('\n')
       }
   }
}

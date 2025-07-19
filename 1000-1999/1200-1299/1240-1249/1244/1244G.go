package main

import (
   "bufio"
   "fmt"
   "os"
   "strconv"
   "strings"
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
   var k int64
   if _, err := fmt.Fscan(reader, &n, &k); err != nil {
       return
   }
   sumMin := int64(n) * int64(n+1) / 2
   if k < sumMin {
       fmt.Fprintln(writer, -1)
       return
   }
   res := sumMin
   // p and q are 1-indexed
   p := make([]int, n+2)
   used := make([]bool, n+2)
   pos := 0
   for res < k && pos < n {
       pos++
       maxDet := int64(n) - 2*int64(pos) + 1
       det := minInt64(k-res, maxDet)
       if det <= 0 {
           pos--
           break
       }
       p[pos] = pos + int(det)
       res += det
   }
   // mark used for first track
   for i := 1; i <= pos; i++ {
       used[p[i]] = true
   }
   // fill remaining for p
   next := pos + 1
   for i := 1; i <= n; i++ {
       if !used[i] {
           p[next] = i
           next++
       }
   }
   // build q
   q := make([]int, n+2)
   cur := pos
   for i := 1; i <= pos; i++ {
       q[cur] = i
       cur--
   }
   for i := pos + 1; i <= n; i++ {
       q[i] = i
   }
   // output
   var sb strings.Builder
   sb.WriteString(strconv.FormatInt(res, 10))
   sb.WriteByte('\n')
   for i := 1; i <= n; i++ {
       sb.WriteString(strconv.Itoa(p[i]))
       if i < n {
           sb.WriteByte(' ')
       }
   }
   sb.WriteByte('\n')
   for i := 1; i <= n; i++ {
       sb.WriteString(strconv.Itoa(q[i]))
       if i < n {
           sb.WriteByte(' ')
       }
   }
   sb.WriteByte('\n')
   writer.WriteString(sb.String())
}

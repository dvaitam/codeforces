package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReaderSize(os.Stdin, 1<<20)
   writer := bufio.NewWriterSize(os.Stdout, 1<<20)
   defer writer.Flush()
   var t int
   fmt.Fscan(reader, &t)
   for tc := 0; tc < t; tc++ {
       var k int
       var s, a, b string
       fmt.Fscan(reader, &k, &s, &a, &b)
       n := len(s)
       p := make([]int, k)
       used := make([]bool, k)
       for i := 0; i < k; i++ {
           p[i] = -1
       }
       sArr := make([]int, n)
       aArr := make([]int, n)
       bArr := make([]int, n)
       for i := 0; i < n; i++ {
           sArr[i] = int(s[i] - 'a')
           aArr[i] = int(a[i] - 'a')
           bArr[i] = int(b[i] - 'a')
       }
       var dfs func(int, bool, bool) bool
       dfs = func(pos int, fa, fb bool) bool {
           for i := pos; i < n; i++ {
               ci, ai, bi := sArr[i], aArr[i], bArr[i]
               if p[ci] != -1 {
                   x := p[ci]
                   if (!fa && x < ai) || (!fb && x > bi) {
                       return false
                   }
                   if !fa && x > ai {
                       fa = true
                   }
                   if !fb && x < bi {
                       fb = true
                   }
                   continue
               }
               lo, hi := 0, k-1
               if !fa {
                   lo = ai
               }
               if !fb {
                   hi = bi
               }
               for x := lo; x <= hi; x++ {
                   if used[x] {
                       continue
                   }
                   newFa := fa || (x > ai)
                   newFb := fb || (x < bi)
                   p[ci], used[x] = x, true
                   if newFa && newFb {
                       // fill remaining arbitrarily
                       idx := 0
                       for j := 0; j < k; j++ {
                           if p[j] == -1 {
                               for used[idx] {
                                   idx++
                               }
                               p[j], used[idx] = idx, true
                           }
                       }
                       return true
                   }
                   if dfs(i+1, newFa, newFb) {
                       return true
                   }
                   // rollback
                   used[x] = false
                   p[ci] = -1
               }
               return false
           }
           return true
       }
       ok := dfs(0, false, false)
       if ok {
           // fill any remaining
           idx := 0
           for j := 0; j < k; j++ {
               if p[j] == -1 {
                   for used[idx] {
                       idx++
                   }
                   p[j], used[idx] = idx, true
               }
           }
           fmt.Fprintln(writer, "YES")
           out := make([]byte, k)
           for i := 0; i < k; i++ {
               out[i] = byte(p[i] + 'a')
           }
           writer.Write(out)
           writer.WriteByte('\n')
       } else {
           fmt.Fprintln(writer, "NO")
       }
   }
}

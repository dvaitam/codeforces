package main

import (
   "bufio"
   "fmt"
   "os"
)

const maxCode = 27 * 27 * 27 * 27

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   words := make([]string, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &words[i])
   }
   // build graph: each word to possible codes
   edges := make([][]int, n)
   // factors for base-27 positions
   f1, f2, f3, f4 := 1, 27, 27*27, 27*27*27
   for i, w := range words {
       l := len(w)
       for a := 0; a < l; a++ {
           for b := a + 1; b <= l; b++ {
               // determine c start similar to C++ min(b+1, l)
               cstart := b + 1
               if cstart > l {
                   cstart = l
               }
               for c := cstart; c <= l; c++ {
                   // determine d start similar to C++ min(c+1, l)
                   dstart := c + 1
                   if dstart > l {
                       dstart = l
                   }
                   for d := dstart; d <= l; d++ {
                       var h int
                       h += (int(w[a]-'a') + 1) * f1
                       if b < l {
                           h += (int(w[b]-'a') + 1) * f2
                       }
                       if c < l {
                           h += (int(w[c]-'a') + 1) * f3
                       }
                       if d < l {
                           h += (int(w[d]-'a') + 1) * f4
                       }
                       edges[i] = append(edges[i], h)
                   }
               }
           }
       }
   }
   // matching arrays
   match := make([]int, maxCode)
   for i := range match {
       match[i] = -1
   }
   matchI := make([]int, n)

   var dfs func(int, []bool) bool
   dfs = func(u int, vis []bool) bool {
       if vis[u] {
           return false
       }
       vis[u] = true
       for _, v := range edges[u] {
           if match[v] < 0 || dfs(match[v], vis) {
               match[v] = u
               matchI[u] = v
               return true
           }
       }
       return false
   }

   // perform matching
   cnt := 0
   for i := 0; i < n; i++ {
       vis := make([]bool, n)
       if dfs(i, vis) {
           cnt++
       }
   }
   if cnt < n {
       fmt.Fprintln(writer, -1)
       return
   }
   // output codes
   for i := 0; i < n; i++ {
       tmp := matchI[i]
       // reconstruct letters
       for tmp > 0 {
           c := tmp % 27
           writer.WriteByte(byte('a' + c - 1))
           tmp /= 27
       }
       writer.WriteByte('\n')
   }
}

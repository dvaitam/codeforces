package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

type pair struct {
   pos int
   cnt int
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   var ts int
   fmt.Fscan(reader, &ts)
   for ts > 0 {
       ts--
       var n, k int
       fmt.Fscan(reader, &n, &k)
       pos := make([]int, k)
       cnt := make([]int, k)
       for i := 0; i < k; i++ {
           fmt.Fscan(reader, &pos[i])
       }
       for i := 0; i < k; i++ {
           fmt.Fscan(reader, &cnt[i])
       }
       ok := true
       lp, lc := 0, 0
       for i := 0; i < k; i++ {
           if pos[i]-lp < cnt[i]-lc {
               ok = false
               break
           }
           if pos[i] <= 3 {
               if cnt[i] != pos[i] {
                   ok = false
                   break
               }
           } else {
               if cnt[i] < 3 {
                   ok = false
                   break
               }
           }
           lp = pos[i]
           lc = cnt[i]
       }
       if !ok {
           fmt.Fprintln(writer, "NO")
           continue
       }
       p := make([]pair, k)
       for i := 0; i < k; i++ {
           p[i] = pair{pos: pos[i] - 1, cnt: cnt[i]}
       }
       sort.Slice(p, func(i, j int) bool { return p[i].pos < p[j].pos })
       res := make([]byte, n)
       for i := range res {
           res[i] = '-'
       }
       if n > 0 {
           res[0] = 'a'
       }
       if n > 1 {
           res[1] = 'b'
       }
       if n > 2 {
           res[2] = 'c'
       }
       pal := 3
       cur := byte('d')
       for _, pr := range p {
           if pr.pos < 3 {
               continue
           }
           times := pr.cnt - pal
           for i := 0; i < times; i++ {
               idx := pr.pos - i
               if idx >= 0 && idx < n {
                   res[idx] = cur
               }
           }
           pal = pr.cnt
           cur++
       }
       cur = 'a'
       for i := 3; i < n; i++ {
           if res[i] == '-' {
               res[i] = cur
               cur++
               if cur > 'c' {
                   cur = 'a'
               }
           }
       }
       fmt.Fprintln(writer, "YES")
       writer.Write(res)
       writer.WriteByte('\n')
   }
}

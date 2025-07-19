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

   var t int
   if _, err := fmt.Fscan(reader, &t); err != nil {
       return
   }
   const INF = 1000000000
   for t > 0 {
       t--
       var s string
       fmt.Fscan(reader, &s)
       var n int
       fmt.Fscan(reader, &n)
       arr := make([]string, n)
       for i := 0; i < n; i++ {
           fmt.Fscan(reader, &arr[i])
       }
       L := len(s)
       ss := make([][2]int, L)
       for i := 0; i < L; i++ {
           ss[i][0], ss[i][1] = INF, INF
       }
       for i, pat := range arr {
           plen := len(pat)
           for j := 0; j+plen <= L; j++ {
               if s[j:j+plen] == pat {
                   end := j + plen - 1
                   if ss[end][1] > j {
                       ss[end][0] = i
                       ss[end][1] = j
                   }
               }
           }
       }
       type pair struct{ idx, pos int }
       var ans []pair
       rpos := L - 1
       for rpos >= 0 {
           mv := rpos
           cs := -1
           for i := rpos; i < L; i++ {
               if ss[i][1] <= rpos && ss[i][1]-1 < mv {
                   mv = ss[i][1] - 1
                   cs = ss[i][0]
               }
           }
           if mv >= rpos {
               break
           }
           ans = append(ans, pair{cs, mv + 1})
           rpos = mv
       }
       if rpos >= 0 {
           fmt.Fprintln(writer, -1)
       } else {
           fmt.Fprintln(writer, len(ans))
           for _, p := range ans {
               fmt.Fprintln(writer, p.idx+1, p.pos+1)
           }
       }
   }
}

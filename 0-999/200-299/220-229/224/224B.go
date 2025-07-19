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
   var N, K int
   if _, err := fmt.Fscan(reader, &N, &K); err != nil {
       return
   }
   data := make([]int, N)
   for i := 0; i < N; i++ {
       fmt.Fscan(reader, &data[i])
   }
   cnt := make(map[int]int)
   for _, v := range data {
       cnt[v]++
   }
   distinct := len(cnt)
   if distinct < K {
       fmt.Fprintln(writer, "-1 -1")
       return
   }
   // find left boundary
   l := 0
   for ; l < N; l++ {
       v := data[l]
       cnt[v]--
       if cnt[v] == 0 {
           distinct--
           if distinct < K {
               cnt[v]++
               distinct++
               break
           }
       }
   }
   // find right boundary
   r := N - 1
   for ; r >= 0; r-- {
       v := data[r]
       cnt[v]--
       if cnt[v] == 0 {
           distinct--
           if distinct < K {
               cnt[v]++
               distinct++
               break
           }
       }
   }
   fmt.Fprintf(writer, "%d %d\n", l+1, r+1)
}

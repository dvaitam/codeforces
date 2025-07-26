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

   var s string
   if _, err := fmt.Fscan(reader, &s); err != nil {
       return
   }
   n := len(s)
   // count single characters
   cnt := make([]int64, 26)
   // count ordered pairs
   pairs := make([][]int64, 26)
   for i := 0; i < 26; i++ {
       pairs[i] = make([]int64, 26)
   }
   for i := 0; i < n; i++ {
       c := s[i] - 'a'
       for j := 0; j < 26; j++ {
           pairs[j][c] += cnt[j]
       }
       cnt[c]++
   }
   // find maximum
   var ans int64
   for i := 0; i < 26; i++ {
       if cnt[i] > ans {
           ans = cnt[i]
       }
   }
   for i := 0; i < 26; i++ {
       for j := 0; j < 26; j++ {
           if pairs[i][j] > ans {
               ans = pairs[i][j]
           }
       }
   }
   fmt.Fprintln(writer, ans)

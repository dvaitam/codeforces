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
   for ti := 0; ti < t; ti++ {
       var n int
       fmt.Fscan(reader, &n)
       var s string
       fmt.Fscan(reader, &s)

       // map from prefix sum to count
       freq := make(map[int]int)
       sum := 0
       freq[0] = 1
       var ans int64
       for i := 0; i < n; i++ {
           sum += int(s[i]-'0') - 1
           if c, ok := freq[sum]; ok {
               ans += int64(c)
           }
           freq[sum]++
       }
       fmt.Fprintln(writer, ans)
   }
}

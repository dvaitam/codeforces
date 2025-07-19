package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   freq := make(map[int]int)
   var ans int64
   for i := 0; i < n; i++ {
       var s string
       fmt.Fscan(reader, &s)
       mask := 0
       for j := 0; j < len(s); j++ {
           mask ^= 1 << (s[j] - 'a')
       }
       // count exact matches
       if c, ok := freq[mask]; ok {
           ans += int64(c)
       }
       // count masks differing by one bit
       for b := 0; b < 26; b++ {
           ans += int64(freq[mask^(1<<b)])
       }
       freq[mask]++
   }
   fmt.Println(ans)
}

package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, x int
   if _, err := fmt.Fscan(reader, &n, &x); err != nil {
       return
   }
   freq := make(map[int]int)
   var ans int64
   for i := 0; i < n; i++ {
       var a int
       fmt.Fscan(reader, &a)
       need := a ^ x
       if cnt, ok := freq[need]; ok {
           ans += int64(cnt)
       }
       freq[a]++
   }
   fmt.Println(ans)
}

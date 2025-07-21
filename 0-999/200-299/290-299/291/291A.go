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
   var id int
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &id)
       if id != 0 {
           freq[id]++
       }
   }
   pairs := 0
   for _, cnt := range freq {
       if cnt > 2 {
           fmt.Println(-1)
           return
       }
       if cnt == 2 {
           pairs++
       }
   }
   fmt.Println(pairs)
}

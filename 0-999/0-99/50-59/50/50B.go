package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var s string
   if _, err := fmt.Fscan(reader, &s); err != nil {
       return
   }
   var freq [256]int64
   for i := 0; i < len(s); i++ {
       freq[s[i]]++
   }
   var ans int64
   for _, v := range freq {
       ans += v * v
   }
   fmt.Println(ans)
}

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
   n := len(s)
   if n < 2 {
       fmt.Println(-1)
       return
   }
   last := s[n-1]
   earliest := -1
   latestGreater := -1
   for i := 0; i < n-1; i++ {
       c := s[i]
       if (c-'0')%2 == 0 {
           if c < last && earliest < 0 {
               earliest = i
           }
           if c > last {
               latestGreater = i
           }
       }
   }
   idx := -1
   if earliest >= 0 {
       idx = earliest
   } else if latestGreater >= 0 {
       idx = latestGreater
   } else {
       fmt.Println(-1)
       return
   }
   bs := []byte(s)
   bs[idx], bs[n-1] = bs[n-1], bs[idx]
   if bs[0] == '0' {
       fmt.Println(-1)
       return
   }
   fmt.Println(string(bs))
}

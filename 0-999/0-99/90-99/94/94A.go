package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var enc string
   if _, err := fmt.Fscan(reader, &enc); err != nil {
       return
   }
   codes := make([]string, 10)
   for i := 0; i < 10; i++ {
       fmt.Fscan(reader, &codes[i])
   }
   var result [8]byte
   for j := 0; j < 8; j++ {
       segment := enc[j*10 : (j+1)*10]
       for d, code := range codes {
           if code == segment {
               result[j] = byte('0' + d)
               break
           }
       }
   }
   fmt.Println(string(result[:]))
}

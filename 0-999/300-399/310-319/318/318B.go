package main

import (
   "bufio"
   "fmt"
   "io"
   "os"
   "strings"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   s, err := reader.ReadString('\n')
   if err != nil && err != io.EOF {
       panic(err)
   }
   s = strings.TrimSpace(s)
   b := []byte(s)
   var heavyCount, result int64
   n := len(b)
   for i := 0; i+4 < n; i++ {
       // check for "heavy"
       if b[i] == 'h' && b[i+1] == 'e' && b[i+2] == 'a' && b[i+3] == 'v' && b[i+4] == 'y' {
           heavyCount++
       } else if b[i] == 'm' && b[i+1] == 'e' && b[i+2] == 't' && b[i+3] == 'a' && b[i+4] == 'l' {
           result += heavyCount
       }
   }
   fmt.Println(result)
}

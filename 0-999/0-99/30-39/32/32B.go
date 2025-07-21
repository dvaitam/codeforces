package main

import (
   "bufio"
   "fmt"
   "os"
   "strings"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   input, err := reader.ReadString('\n')
   if err != nil && err.Error() != "EOF" {
       fmt.Fprintln(os.Stderr, "error reading input:", err)
       return
   }
   s := strings.TrimSpace(input)
   var b strings.Builder
   for i := 0; i < len(s); i++ {
       switch s[i] {
       case '.':
           b.WriteByte('0')
       case '-':
           if i+1 < len(s) && s[i+1] == '.' {
               b.WriteByte('1')
           } else {
               b.WriteByte('2')
           }
           i++
       }
   }
   fmt.Println(b.String())
}

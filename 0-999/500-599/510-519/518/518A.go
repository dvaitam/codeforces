package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var s, t string
   if _, err := fmt.Fscan(reader, &s, &t); err != nil {
       return
   }
   n := len(s)
   // compute next string after s
   // if s is all 'z', no next
   next := []rune(s)
   i := n - 1
   for ; i >= 0; i-- {
       if next[i] != 'z' {
           next[i]++
           break
       }
       next[i] = 'a'
   }
   if i < 0 {
       fmt.Println("No such string")
       return
   }
   res := string(next)
   // check lex order
   if res < t {
       fmt.Println(res)
   } else {
       fmt.Println("No such string")
   }
}

package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   fmt.Fscan(reader, &n)
   var s string
   fmt.Fscan(reader, &s)
   // convert to slice of bytes for easy removal
   a := []byte(s)
   removed := 0
   // iterate from 'z' to 'b'
   for c := byte('z'); c > byte('a'); c-- {
       for {
           pos := -1
           // find rightmost occurrence of c with neighbor c-1
           for i := 0; i < len(a); i++ {
               if a[i] != c {
                   continue
               }
               if (i > 0 && a[i-1] == c-1) || (i+1 < len(a) && a[i+1] == c-1) {
                   pos = i
               }
           }
           if pos == -1 {
               break
           }
           // remove a[pos]
           a = append(a[:pos], a[pos+1:]...)
           removed++
       }
   }
   fmt.Println(removed)
}

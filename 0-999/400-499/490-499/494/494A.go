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
   // count opens, fixed closes, and positions of '#'
   open := 0
   close := 0
   hashCount := 0
   for i := 0; i < n; i++ {
       switch s[i] {
       case '(':
           open++
       case ')':
           close++
       case '#':
           hashCount++
       }
   }
   // initial assign 1 ')' to each '#'
   // total closes now close + hashCount
   totalClose := close + hashCount
   // compute extra to assign to last '#'
   extra := open - totalClose
   if extra < 0 {
       fmt.Println(-1)
       return
   }
   // prepare assignments
   assigns := make([]int, hashCount)
   for i := 0; i < hashCount; i++ {
       assigns[i] = 1
   }
   if hashCount > 0 {
       assigns[hashCount-1] += extra
   }
   // simulate to check validity
   bal := 0
   hi := 0 // index for assigns
   for i := 0; i < n; i++ {
       switch s[i] {
       case '(':
           bal++
       case ')':
           bal--
       case '#':
           // subtract assigned closes
           bal -= assigns[hi]
           hi++
       }
       if bal < 0 {
           fmt.Println(-1)
           return
       }
   }
   if bal != 0 {
       fmt.Println(-1)
       return
   }
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   for _, v := range assigns {
       fmt.Fprintln(out, v)
   }
}

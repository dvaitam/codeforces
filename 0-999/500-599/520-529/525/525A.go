package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   var s string
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   // read string of keys and doors
   fmt.Fscan(reader, &s)
   // counts of available keys a-z
   keys := make([]int, 26)
   bought := 0
   // traverse rooms 1 to n-1
   for i := 0; i < n-1; i++ {
       // at index 2*i: key in current room
       k := s[2*i]
       keys[k-'a']++
       // at index 2*i+1: door to next room
       d := s[2*i+1]
       idx := d - 'A'
       if keys[idx] > 0 {
           keys[idx]--
       } else {
           bought++
       }
   }
   fmt.Println(bought)
}

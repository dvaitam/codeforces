package main

import (
   "bufio"
   "fmt"
   "os"
   "strings"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   s, err := reader.ReadString('\n')
   if err != nil && len(s) == 0 {
       return
   }
   s = strings.TrimRight(s, "\n")
   pivot := strings.IndexByte(s, '^')
   var leftTorque, rightTorque int64
   // compute left torque
   for i := 0; i < pivot; i++ {
       c := s[i]
       if c >= '1' && c <= '9' {
           w := int64(c - '0')
           leftTorque += w * int64(pivot-i)
       }
   }
   // compute right torque
   n := len(s)
   for i := pivot + 1; i < n; i++ {
       c := s[i]
       if c >= '1' && c <= '9' {
           w := int64(c - '0')
           rightTorque += w * int64(i-pivot)
       }
   }
   switch {
   case leftTorque > rightTorque:
       fmt.Println("left")
   case leftTorque < rightTorque:
       fmt.Println("right")
   default:
       fmt.Println("balance")
   }
}

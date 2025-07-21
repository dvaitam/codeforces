package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, l int
   if _, err := fmt.Fscan(reader, &n, &l); err != nil {
       return
   }
   a := make([]int, n)
   maxA := 0
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
       if a[i] > maxA {
           maxA = a[i]
       }
   }
   maxArea := 0
   // try each possible piece length d >= l
   for d := l; d <= maxA; d++ {
       total := 0
       for _, ai := range a {
           total += ai / d
       }
       area := total * d
       if area > maxArea {
           maxArea = area
       }
   }
   fmt.Println(maxArea)
}

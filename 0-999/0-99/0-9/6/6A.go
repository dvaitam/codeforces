package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var sticks [4]int
   for i := 0; i < 4; i++ {
       if _, err := fmt.Fscan(reader, &sticks[i]); err != nil {
           return
       }
   }
   triangle := false
   segment := false
   // check all combinations of three sticks
   for i := 0; i < 4; i++ {
       for j := i + 1; j < 4; j++ {
           for k := j + 1; k < 4; k++ {
               // get three lengths
               a, b, c := sticks[i], sticks[j], sticks[k]
               // sort a, b, c so that a <= b <= c
               if a > b { a, b = b, a }
               if b > c { b, c = c, b }
               if a > b { a, b = b, a }
               if a + b > c {
                   triangle = true
               } else if a + b == c {
                   segment = true
               }
           }
       }
   }
   switch {
   case triangle:
       fmt.Println("TRIANGLE")
   case segment:
       fmt.Println("SEGMENT")
   default:
       fmt.Println("IMPOSSIBLE")
   }
}
 

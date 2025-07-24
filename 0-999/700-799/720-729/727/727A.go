package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var a, b int64
   if _, err := fmt.Fscan(reader, &a, &b); err != nil {
       return
   }
   // reverse operations from b to a
   seq := make([]int64, 0, 64)
   seq = append(seq, b)
   for b > a {
       if b%10 == 1 {
           b = (b - 1) / 10
           seq = append(seq, b)
       } else if b%2 == 0 {
           b = b / 2
           seq = append(seq, b)
       } else {
           break
       }
   }
   if b != a {
       fmt.Println("NO")
       return
   }
   // successful, print sequence from a to original b
   fmt.Println("YES")
   k := len(seq)
   fmt.Println(k)
   // reverse seq
   for i := k - 1; i >= 0; i-- {
       fmt.Print(seq[i])
       if i > 0 {
           fmt.Print(' ')
       }
   }
   fmt.Println()
}

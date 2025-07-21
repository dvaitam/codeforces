package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   for k := 1; ; k++ {
       t := k * (k + 1) / 2
       if t == n {
           fmt.Println("YES")
           return
       }
       if t > n {
           fmt.Println("NO")
           return
       }
   }
}

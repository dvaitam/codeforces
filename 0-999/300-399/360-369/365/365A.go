package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, k int
   if _, err := fmt.Fscan(in, &n, &k); err != nil {
       return
   }
   count := 0
   for i := 0; i < n; i++ {
       var s string
       fmt.Fscan(in, &s)
       found := make([]bool, k+1)
       needed := k + 1
       for _, c := range s {
           d := int(c - '0')
           if d <= k && !found[d] {
               found[d] = true
               needed--
           }
       }
       if needed == 0 {
           count++
       }
   }
   fmt.Println(count)
}

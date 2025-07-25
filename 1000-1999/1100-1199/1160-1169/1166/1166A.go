package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   counts := make([]int, 26)
   for i := 0; i < n; i++ {
       var name string
       fmt.Fscan(reader, &name)
       if len(name) > 0 {
           counts[name[0]-'a']++
       }
   }
   result := 0
   for _, c := range counts {
       a := c / 2
       b := c - a
       result += a*(a-1)/2 + b*(b-1)/2
   }
   fmt.Println(result)
}

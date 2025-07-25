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
   seen := make(map[int]struct{})
   count := 0
   for {
       if _, ok := seen[n]; ok {
           break
       }
       seen[n] = struct{}{}
       count++
       n = f(n)
   }
   fmt.Println(count)
}

// f increments x by 1, then removes all trailing zeros
func f(x int) int {
   x++
   for x%10 == 0 {
       x /= 10
   }
   return x
}

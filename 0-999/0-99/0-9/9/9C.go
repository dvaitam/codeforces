package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int64
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   var count int64
   queue := []int64{1}
   for i := 0; i < len(queue); i++ {
       x := queue[i]
       if x > n {
           continue
       }
       count++
       x10 := x * 10
       if x10 <= n {
           queue = append(queue, x10)
       }
       if x10+1 <= n {
           queue = append(queue, x10+1)
       }
   }
   fmt.Println(count)
}

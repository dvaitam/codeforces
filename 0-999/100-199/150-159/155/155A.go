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
   var score int
   // Read first score
   if n <= 0 {
       fmt.Println(0)
       return
   }
   fmt.Fscan(reader, &score)
   best, worst := score, score
   count := 0
   for i := 1; i < n; i++ {
       if _, err := fmt.Fscan(reader, &score); err != nil {
           break
       }
       if score > best {
           count++
           best = score
       } else if score < worst {
           count++
           worst = score
       }
   }
   fmt.Println(count)
}

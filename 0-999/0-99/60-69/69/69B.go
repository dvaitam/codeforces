package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   l := make([]int, m)
   r := make([]int, m)
   t := make([]int, m)
   c := make([]int, m)
   for i := 0; i < m; i++ {
       fmt.Fscan(reader, &l[i], &r[i], &t[i], &c[i])
   }
   total := 0
   for j := 1; j <= n; j++ {
       // initialize bestTime larger than any possible ti (<=1000)
       bestTime := 1000000000
       bestIdx := -1
       for i := 0; i < m; i++ {
           if l[i] <= j && j <= r[i] {
               if t[i] < bestTime || (t[i] == bestTime && i < bestIdx) {
                   bestTime = t[i]
                   bestIdx = i
               }
           }
       }
       if bestIdx != -1 {
           total += c[bestIdx]
       }
   }
   fmt.Println(total)
}

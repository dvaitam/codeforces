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
   names := make([]string, n)
   scores := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &names[i], &scores[i])
   }
   // compute final scores
   total := make(map[string]int)
   for i := 0; i < n; i++ {
       total[names[i]] += scores[i]
   }
   // find max score
   maxScore := -1000000000 // lower than any possible total (min total >= -1e6)
   for _, s := range total {
       if s > maxScore {
           maxScore = s
       }
   }
   // collect candidates with maxScore
   candidates := make(map[string]bool)
   for name, s := range total {
       if s == maxScore {
           candidates[name] = true
       }
   }
   // find winner by earliest reaching at least maxScore
   roundTotal := make(map[string]int)
   for i := 0; i < n; i++ {
       name := names[i]
       roundTotal[name] += scores[i]
       if candidates[name] && roundTotal[name] >= maxScore {
           fmt.Println(name)
           return
       }
   }
}

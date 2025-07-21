package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, k int
   if _, err := fmt.Fscan(reader, &n, &k); err != nil {
       return
   }
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   // Map value to positions
   pos := make(map[int][]int)
   curMax := 0
   var ans int64
   for i, v := range a {
       idx := i + 1
       lst := pos[v]
       lst = append(lst, idx)
       pos[v] = lst
       if len(lst) >= k {
           // position of k-th last occurrence for v
           t := lst[len(lst)-k]
           if t > curMax {
               curMax = t
           }
       }
       if curMax > 0 {
           ans += int64(curMax)
       }
   }
   fmt.Println(ans)
}

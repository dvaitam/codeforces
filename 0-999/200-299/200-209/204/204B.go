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
   fronts := make(map[int]int)
   backs := make(map[int]int)
   // Read cards
   for i := 0; i < n; i++ {
       var a, b int
       fmt.Fscan(reader, &a, &b)
       fronts[a]++
       if b != a {
           backs[b]++
       }
   }
   needed := (n + 1) / 2
   const INF = 1<<60
   ans := INF
   // Consider colors seen on front
   for color, fcnt := range fronts {
       if fcnt >= needed {
           ans = 0
           break
       }
       need := needed - fcnt
       if bcnt, ok := backs[color]; ok && bcnt >= need {
           if need < ans {
               ans = need
           }
       }
   }
   // Consider colors only on back
   for color, bcnt := range backs {
       if _, seen := fronts[color]; seen {
           continue
       }
       if bcnt >= needed && needed < ans {
           ans = needed
       }
   }
   if ans == INF {
       fmt.Println(-1)
   } else {
       fmt.Println(ans)
   }
}

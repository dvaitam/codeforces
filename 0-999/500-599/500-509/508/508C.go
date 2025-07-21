package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var m, t, r int
   if _, err := fmt.Fscan(reader, &m, &t, &r); err != nil {
       return
   }
   w := make([]int, m)
   for i := 0; i < m; i++ {
       fmt.Fscan(reader, &w[i])
   }
   const offset = 300
   const size = 1000
   used := make([]bool, size)
   total := 0
   for _, wi := range w {
       cnt := 0
       // count already burning candles at time wi
       for s := wi - t; s <= wi-1; s++ {
           idx := s + offset
           if idx >= 0 && idx < size && used[idx] {
               cnt++
           }
       }
       if cnt >= r {
           continue
       }
       need := r - cnt
       // light new candles as late as possible
       for s := wi - 1; s >= wi-t && need > 0; s-- {
           idx := s + offset
           if idx >= 0 && idx < size && !used[idx] {
               used[idx] = true
               total++
               need--
           }
       }
       if need > 0 {
           fmt.Println(-1)
           return
       }
   }
   fmt.Println(total)
}

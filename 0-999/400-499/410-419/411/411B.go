package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, m, k int
   if _, err := fmt.Fscan(reader, &n, &m, &k); err != nil {
       return
   }
   instr := make([][]int, n)
   for i := 0; i < n; i++ {
       instr[i] = make([]int, m)
       for j := 0; j < m; j++ {
           fmt.Fscan(reader, &instr[i][j])
       }
   }

   coreLocked := make([]bool, n)
   coreLockTime := make([]int, n)
   cellLocked := make([]bool, k+1)

   for cycle := 1; cycle <= m; cycle++ {
       // map cell -> list of cores writing this cycle
       writes := make(map[int][]int)
       for i := 0; i < n; i++ {
           if coreLocked[i] {
               continue
           }
           x := instr[i][cycle-1]
           if x == 0 {
               continue
           }
           if cellLocked[x] {
               coreLocked[i] = true
               coreLockTime[i] = cycle
           } else {
               writes[x] = append(writes[x], i)
           }
       }
       // process write conflicts
       for cell, cores := range writes {
           if len(cores) >= 2 {
               cellLocked[cell] = true
               for _, ci := range cores {
                   if !coreLocked[ci] {
                       coreLocked[ci] = true
                       coreLockTime[ci] = cycle
                   }
               }
           }
       }
   }

   for i := 0; i < n; i++ {
       fmt.Fprintln(writer, coreLockTime[i])
   }
}

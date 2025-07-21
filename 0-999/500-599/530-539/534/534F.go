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
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   b := make([]int, m)
   for j := 0; j < m; j++ {
       fmt.Fscan(reader, &b[j])
   }
   // generate all row patterns matching segment counts
   patterns := make([][]uint, n)
   for i := 0; i < n; i++ {
       patterns[i] = genPatterns(m, a[i])
   }
   // state for columns
   segCount := make([]int, m)
   prev := make([]int, m)
   sol := make([]uint, n)
   var solved bool
   var dfs func(int)
   dfs = func(row int) {
       if solved {
           return
       }
       if row == n {
           // check final column counts
           for j := 0; j < m; j++ {
               if segCount[j] != b[j] {
                   return
               }
           }
           // output solution
           w := bufio.NewWriter(os.Stdout)
           defer w.Flush()
           for i := 0; i < n; i++ {
               for j := 0; j < m; j++ {
                   if (sol[i]>>uint(j))&1 == 1 {
                       w.WriteByte('*')
                   } else {
                       w.WriteByte('.')
                   }
               }
               w.WriteByte('\n')
           }
           solved = true
           return
       }
       // try each pattern
       for _, pat := range patterns[row] {
           // backup column state
           backupSeg := make([]int, m)
           backupPrev := make([]int, m)
           copy(backupSeg, segCount)
           copy(backupPrev, prev)
           ok := true
           for j := 0; j < m; j++ {
               bit := int((pat >> uint(j)) & 1)
               if bit == 1 && prev[j] == 0 {
                   segCount[j]++
               }
               if segCount[j] > b[j] {
                   ok = false
                   break
               }
               prev[j] = bit
           }
           if ok {
               sol[row] = pat
               dfs(row + 1)
           }
           // restore state
           copy(segCount, backupSeg)
           copy(prev, backupPrev)
           if solved {
               return
           }
       }
   }
   dfs(0)
}

// genPatterns returns all bitmasks of length m with exactly want segments
func genPatterns(m, want int) []uint {
   var res []uint
   limit := 1 << m
   for mask := 0; mask < limit; mask++ {
       cnt := 0
       prev := 0
       for j := 0; j < m; j++ {
           bit := (mask >> j) & 1
           if bit == 1 && prev == 0 {
               cnt++
           }
           prev = bit
           if cnt > want {
               break
           }
       }
       if cnt == want {
           res = append(res, uint(mask))
       }
   }
   return res
}

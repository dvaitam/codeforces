package main

import (
   "bufio"
   "fmt"
   "os"
)

func minInt64(a, b int64) int64 {
   if a < b {
       return a
   }
   return b
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var m, n, k int64
   if _, err := fmt.Fscan(reader, &m, &n, &k); err != nil {
       return
   }
   // maximum possible moves
   maxMoves := 4*m*n - 2*(m + n)
   if k > maxMoves {
       fmt.Fprintln(writer, "NO")
       return
   }
   fmt.Fprintln(writer, "YES")
   type cmd struct { cnt int64; s string }
   var res []cmd
   var run int64
   // move right on first row
   if n > 1 {
       cnt := minInt64(k, n-1)
       if cnt > 0 {
           res = append(res, cmd{cnt, "R"})
           run += cnt
       }
   }
   // sweep columns
   for i := n; i > 1 && run < k; i-- {
       rem := k - run
       downMoves := (rem) / 3
       if downMoves == 0 {
           if m > 1 {
               if rem == 1 {
                   res = append(res, cmd{1, "D"})
               } else if rem == 2 {
                   res = append(res, cmd{1, "DL"})
               }
               run += rem
           }
       } else if downMoves < m-1 {
           // partial vertical sweeps
           if downMoves > 0 {
               res = append(res, cmd{downMoves, "DLR"})
               run += downMoves * 3
           }
           if run < k {
               // one more down
               res = append(res, cmd{1, "D"})
               run++
               if run < k {
                   res = append(res, cmd{1, "U"})
                   run++
               }
           }
       } else {
           // full vertical sweeps
           if m > 1 {
               cnt := m - 1
               res = append(res, cmd{cnt, "DLR"})
               run += cnt * 3
               if run < k {
                   upCnt := minInt64(k-run, m-1)
                   res = append(res, cmd{upCnt, "U"})
                   run += upCnt
               }
           }
       }
       if run < k {
           // move left to next column
           res = append(res, cmd{1, "L"})
           run++
       }
   }
   // final column adjustments
   if run < k && m > 1 {
       // go down
       rem := k - run
       cnt := minInt64(rem, m-1)
       if cnt > 0 {
           res = append(res, cmd{cnt, "D"})
           run += cnt
       }
       if run < k {
           rem2 := k - run
           cnt2 := minInt64(rem2, m-1)
           if cnt2 > 0 {
               res = append(res, cmd{cnt2, "U"})
               run += cnt2
           }
       }
   }
   // output
   fmt.Fprintln(writer, len(res))
   for _, c := range res {
       fmt.Fprintln(writer, c.cnt, c.s)
   }
}

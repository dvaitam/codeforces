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
   var t int
   fmt.Fscan(reader, &t)
   for ; t > 0; t-- {
       var s string
       fmt.Fscan(reader, &s)
       n := len(s)
       idx := make([][]int, 4)
       for i := range idx {
           idx[i] = make([]int, 0, n)
       }
       for i, ch := range s {
           var v int
           switch ch {
           case 'A': v = 0
           case 'N': v = 1
           case 'T': v = 2
           case 'O': v = 3
           }
           idx[v] = append(idx[v], i)
       }
       var swaps [4][4]int64
       for i := 0; i < 4; i++ {
           for j := 0; j < 4; j++ {
               if i == j {
                   continue
               }
               ii, jj := 0, 0
               ni, nj := len(idx[i]), len(idx[j])
               for ii < ni && jj < nj {
                   if idx[i][ii] < idx[j][jj] {
                       swaps[i][j] += int64(jj)
                       ii++
                   } else {
                       jj++
                   }
               }
               swaps[i][j] += int64(ni-ii) * int64(jj)
           }
       }
       bestVal := int64(-1)
       order := [4]int{0, 1, 2, 3}
       for i := 0; i < 4; i++ {
           for j := 0; j < 4; j++ {
               if j == i {
                   continue
               }
               for k := 0; k < 4; k++ {
                   if k == i || k == j {
                       continue
                   }
                   l := 6 - i - j - k
                   if l < 0 || l >= 4 || l == i || l == j || l == k {
                       continue
                   }
                   var tans int64
                   tans += swaps[i][j] + swaps[i][k] + swaps[i][l]
                   tans += swaps[j][k] + swaps[j][l]
                   tans += swaps[k][l]
                   if tans > bestVal {
                       bestVal = tans
                       order = [4]int{i, j, k, l}
                   }
               }
           }
       }
       mp := []byte{'A', 'N', 'T', 'O'}
       res := make([]byte, 0, n)
       for _, c := range order {
           for range idx[c] {
               res = append(res, mp[c])
           }
       }
       writer.Write(res)
       writer.WriteByte('\n')
   }
}

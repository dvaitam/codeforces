package main

import (
   "bufio"
   "fmt"
   "os"
)

var (
   n   int
   a   []int
   ans int
   pw  [20]int
   c   [20]int
   b   [20][]bool
   v   []int
)

func dfs(k, ds int) {
   if k == 17 {
       if b[k][0] {
           c[17] = -1
           ds++
       } else if b[k][2] {
           c[17] = 1
           ds++
       } else {
           c[17] = 0
       }
       if ds < ans {
           ans = ds
           v = v[:0]
           for i := 0; i <= 17; i++ {
               if c[i] == 1 {
                   v = append(v, 1<<i)
               } else if c[i] == -1 {
                   v = append(v, -(1<<i))
               }
           }
       }
       c[17] = 0
       return
   }
   // check if any odd index has true
   flag := false
   for i := 1; i <= pw[18-k]; i += 2 {
       if b[k][i] {
           flag = true
           break
       }
   }
   if !flag {
       // all even only
       for i := 0; i <= pw[17-k]; i++ {
           b[k+1][i] = b[k][i<<1]
       }
       c[k] = 0
       dfs(k+1, ds)
       return
   }
   // try +1 (shift right)
   for i, j := pw[17-k], 0; i <= pw[18-k]; i, j = i+2, j+1 {
       b[k+1][pw[16-k]+j] = b[k][i] || b[k][i+1]
   }
   for i, j := pw[17-k]-1, 1; i >= 0; i, j = i-2, j+1 {
       b[k+1][pw[16-k]-j] = b[k][i] || b[k][i-1]
   }
   c[k] = 1
   dfs(k+1, ds+1)
   // try -1 (shift left)
   for i, j := pw[17-k], 0; i >= 0; i, j = i-2, j+1 {
       b[k+1][pw[16-k]-j] = b[k][i] || b[k][i-1]
   }
   for i, j := pw[17-k]+1, 1; i <= pw[18-k]; i, j = i+2, j+1 {
       b[k+1][pw[16-k]+j] = b[k][i] || b[k][i+1]
   }
   c[k] = -1
   dfs(k+1, ds+1)
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   fmt.Fscan(reader, &n)
   a = make([]int, n)
   pw[0] = 1
   for i := 1; i < 20; i++ {
       pw[i] = pw[i-1] << 1
   }
   maxSz := pw[18] + 1
   for i := range b {
       b[i] = make([]bool, maxSz)
   }
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
       b[0][a[i]+pw[17]] = true
   }
   ans = 30
   dfs(0, 0)
   // output
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fprintln(writer, ans)
   for i, val := range v {
       if i > 0 {
           fmt.Fprint(writer, " ")
       }
       fmt.Fprint(writer, val)
   }
   fmt.Fprintln(writer)
}

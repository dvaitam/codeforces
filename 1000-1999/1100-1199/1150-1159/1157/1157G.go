package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, m int
   if _, err := fmt.Fscan(in, &n, &m); err != nil {
       return
   }
   a := make([][]int, n)
   for i := 0; i < n; i++ {
       a[i] = make([]int, m)
       for j := 0; j < m; j++ {
           fmt.Fscan(in, &a[i][j])
       }
   }
   // forbidden transitions for columns: F[j][u][v] = true if c[j]=u, c[j+1]=v is forbidden
   F := make([][2][2]bool, m-1)
   for i := 0; i < n; i++ {
       for j := 0; j+1 < m; j++ {
           x := a[i][j]
           y := a[i][j+1]
           u := 1 - x
           v := y
           F[j][u][v] = true
       }
   }
   // dp over columns
   dp := make([][2]bool, m)
   parent := make([][2]int, m)
   // initial: column 0 can be 0 or 1
   dp[0][0] = true
   dp[0][1] = true
   for j := 0; j+1 < m; j++ {
       for b := 0; b < 2; b++ {
           if !dp[j][b] {
               continue
           }
           for nb := 0; nb < 2; nb++ {
               if !F[j][b][nb] {
                   if !dp[j+1][nb] {
                       dp[j+1][nb] = true
                       parent[j+1][nb] = b
                   }
               }
           }
       }
   }
   // check end
   c := make([]int, m)
   if m > 0 {
       if dp[m-1][0] {
           c[m-1] = 0
       } else if dp[m-1][1] {
           c[m-1] = 1
       } else {
           fmt.Println("NO")
           return
       }
       // backtrack
       for j := m - 1; j > 0; j-- {
           c[j-1] = parent[j][c[j]]
       }
   }
   // compute r
   r := make([]int, n)
   if n > 0 {
       r[0] = 0
       for i := 0; i+1 < n; i++ {
           d := a[i][m-1] ^ c[m-1] ^ a[i+1][0] ^ c[0]
           if d == 0 {
               r[i+1] = r[i]
           } else {
               r[i+1] = 0
           }
       }
   }
   // output
   fmt.Println("YES")
   // rows
   for i := 0; i < n; i++ {
       fmt.Fprint(os.Stdout, byte('0'+r[i]))
   }
   fmt.Fprintln(os.Stdout)
   // columns
   for j := 0; j < m; j++ {
       fmt.Fprint(os.Stdout, byte('0'+c[j]))
   }
   fmt.Fprintln(os.Stdout)
}

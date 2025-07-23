package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   fmt.Fscan(reader, &n)
   board := make([][]byte, n)
   for i := 0; i < n; i++ {
       var s string
       fmt.Fscan(reader, &s)
       board[i] = []byte(s)
   }
   // possible moves
   size := 2*n - 1
   center := n - 1
   valid := make([][]bool, size)
   for i := range valid {
       valid[i] = make([]bool, size)
       for j := range valid[i] {
           valid[i][j] = true
       }
   }
   // disallow zero shift
   valid[center][center] = false
   // eliminate invalid shifts
   for dx := -center; dx <= center; dx++ {
       for dy := -center; dy <= center; dy++ {
           kx := dx + center
           ky := dy + center
           if kx == center && ky == center {
               continue
           }
           ok := true
           for i := 0; i < n && ok; i++ {
               for j := 0; j < n; j++ {
                   if board[i][j] != 'o' {
                       continue
                   }
                   ii := i + dx
                   jj := j + dy
                   if ii >= 0 && ii < n && jj >= 0 && jj < n {
                       if board[ii][jj] == '.' {
                           ok = false
                           break
                       }
                   }
               }
           }
           valid[kx][ky] = ok
       }
   }
   // simulate attacks
   got := make([][]byte, n)
   for i := range got {
       got[i] = make([]byte, n)
       for j := range got[i] {
           got[i][j] = '.'
       }
   }
   for i := 0; i < n; i++ {
       for j := 0; j < n; j++ {
           if board[i][j] == 'o' {
               got[i][j] = 'o'
               for dx := -center; dx <= center; dx++ {
                   for dy := -center; dy <= center; dy++ {
                       kx := dx + center
                       ky := dy + center
                       if !valid[kx][ky] {
                           continue
                       }
                       ii := i + dx
                       jj := j + dy
                       if ii >= 0 && ii < n && jj >= 0 && jj < n {
                           if got[ii][jj] == '.' {
                               got[ii][jj] = 'x'
                           }
                       }
                   }
               }
           }
       }
   }
   // compare
   for i := 0; i < n; i++ {
       for j := 0; j < n; j++ {
           if board[i][j] == 'x' && got[i][j] != 'x' {
               fmt.Println("NO")
               return
           }
           if board[i][j] == '.' && got[i][j] == 'x' {
               fmt.Println("NO")
               return
           }
       }
   }
   // output moves board
   fmt.Println("YES")
   // moves board size x size
   for dx := -center; dx <= center; dx++ {
       for dy := -center; dy <= center; dy++ {
           kx := dx + center
           ky := dy + center
           if kx == center && ky == center {
               fmt.Print('o')
           } else if valid[kx][ky] {
               fmt.Print('x')
           } else {
               fmt.Print('.')
           }
       }
       fmt.Println()
   }
}

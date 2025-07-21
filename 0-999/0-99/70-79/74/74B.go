package main

import (
   "bufio"
   "fmt"
   "os"
   "strings"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   // read first line: n, m, k
   line, err := reader.ReadString('\n')
   if err != nil {
       return
   }
   line = strings.TrimSpace(line)
   var n, m, k int
   if _, err := fmt.Sscan(line, &n, &m, &k); err != nil {
       return
   }
   // read direction line
   dirLine, err := reader.ReadString('\n')
   if err != nil {
       return
   }
   dirLine = strings.TrimSpace(dirLine)
   dir := 0
   if strings.Contains(dirLine, "tail") {
       dir = 1
   } else {
       dir = -1
   }
   // read train state string
   stLine, err := reader.ReadString('\n')
   if err != nil && err.Error() != "EOF" {
       return
   }
   st := strings.TrimSpace(stLine)
   T := len(st)

   // reachable positions of stowaway
   S := make([]bool, n+1)
   S[m] = true
   // controller position and direction
   cpos := k

   for i := 1; i <= T; i++ {
       ch := st[i-1]
       if ch == '0' {
           // moving: stowaway moves first
           U := make([]bool, n+1)
           for p := 1; p <= n; p++ {
               if !S[p] {
                   continue
               }
               for _, dp := range []int{-1, 0, 1} {
                   q := p + dp
                   if q < 1 || q > n {
                       continue
                   }
                   if q == cpos {
                       // would be caught before controller moves
                       continue
                   }
                   U[q] = true
               }
           }
           // controller moves
           newpos := cpos + dir
           newdir := dir
           if newpos == 1 || newpos == n {
               newdir = -dir
           }
           // remove positions where controller arrives
           if newpos >= 1 && newpos <= n {
               U[newpos] = false
           }
           // check if no positions left
           ok := false
           for p := 1; p <= n; p++ {
               if U[p] {
                   ok = true
                   break
               }
           }
           if !ok {
               fmt.Printf("Controller %d", i)
               return
           }
           S = U
           cpos = newpos
           dir = newdir
       } else {
           // idle minute: stowaway leaves, controller moves, then re-enters
           newpos := cpos + dir
           newdir := dir
           if newpos == 1 || newpos == n {
               newdir = -dir
           }
           // terminal station at last minute: stowaway wins
           if i == T {
               fmt.Print("Stowaway")
               return
           }
           // re-enter into any wagon except controller's
           U := make([]bool, n+1)
           for p := 1; p <= n; p++ {
               if p != newpos {
                   U[p] = true
               }
           }
           S = U
           cpos = newpos
           dir = newdir
       }
   }
   // if loop ends, stowaway wins
   fmt.Print("Stowaway")
}

package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var a, b int64
   fmt.Fscan(reader, &a, &b)
   var s string
   fmt.Fscan(reader, &s)
   var dx, dy int64
   for _, ch := range s {
       switch ch {
       case 'U': dy++
       case 'D': dy--
       case 'L': dx--
       case 'R': dx++
       }
   }
   var px, py int64
   // check all prefix positions
   for i := 0; i <= len(s); i++ {
       rx := a - px
       ry := b - py
       if reachable(rx, ry, dx, dy) {
           fmt.Println("Yes")
           return
       }
       if i < len(s) {
           switch s[i] {
           case 'U': py++
           case 'D': py--
           case 'L': px--
           case 'R': px++
           }
       }
   }
   fmt.Println("No")
}

// reachable checks if there exists k >= 0 such that k*dx == rx and k*dy == ry
func reachable(rx, ry, dx, dy int64) bool {
   if dx == 0 && dy == 0 {
       return rx == 0 && ry == 0
   }
   if dx == 0 {
       if rx != 0 {
           return false
       }
       if ry%dy != 0 {
           return false
       }
       k := ry / dy
       return k >= 0
   }
   if dy == 0 {
       if ry != 0 {
           return false
       }
       if rx%dx != 0 {
           return false
       }
       k := rx / dx
       return k >= 0
   }
   if rx%dx != 0 || ry%dy != 0 {
       return false
   }
   kx := rx / dx
   ky := ry / dy
   return kx == ky && kx >= 0
}

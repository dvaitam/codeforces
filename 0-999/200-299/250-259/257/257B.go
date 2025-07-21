package main

import "fmt"

// simulate returns the number of blocks (color segments) when starting with start color
func simulate(r, b int, start byte) int {
   rRem, bRem := r, b
   last := start
   if start == 'R' {
       rRem--
   } else {
       bRem--
   }
   blocks := 1
   total := r + b
   for move := 2; move <= total; move++ {
       if move%2 == 1 {
           // Petya's turn: try to keep same color
           if last == 'R' {
               if rRem > 0 {
                   rRem--
               } else {
                   bRem--
                   last = 'B'
                   blocks++
               }
           } else {
               if bRem > 0 {
                   bRem--
               } else {
                   rRem--
                   last = 'R'
                   blocks++
               }
           }
       } else {
           // Vasya's turn: try to switch color
           if last == 'R' {
               if bRem > 0 {
                   bRem--
                   last = 'B'
                   blocks++
               } else {
                   rRem--
               }
           } else {
               if rRem > 0 {
                   rRem--
                   last = 'R'
                   blocks++
               } else {
                   bRem--
               }
           }
       }
   }
   return blocks
}

func main() {
   var n, m int
   if _, err := fmt.Scan(&n, &m); err != nil {
       return
   }
   // simulate starting with Red or Blue and choose minimal blocks
   blocksR := simulate(n, m, 'R')
   blocksB := simulate(n, m, 'B')
   blocks := blocksR
   if blocksB < blocksR {
       blocks = blocksB
   }
   // petya's points = total - blocks, vasya's = blocks - 1
   total := n + m
   petya := total - blocks
   vasya := blocks - 1
   fmt.Println(petya, vasya)
}

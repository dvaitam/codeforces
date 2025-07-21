package main

import (
   "bufio"
   "fmt"
   "os"
)

const MOD = 1000003

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   grid := make([]string, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &grid[i])
   }
   total := 1
   // Horizontal: rows
   for i := 0; i < n; i++ {
       row := grid[i]
       ways := 0
       // two patterns: p=0: L,R,L,R,... ; p=1: R,L,R,L,...
       for p := 0; p < 2; p++ {
           ok := true
           for j := 0; j < m; j++ {
               c := row[j]
               // expected horizontal
               var exp byte
               if (j%2 == 0) == (p == 0) {
                   exp = 'L'
               } else {
                   exp = 'R'
               }
               // actual if fixed
               if c == '1' || c == '3' {
                   if exp != 'L' {
                       ok = false
                       break
                   }
               } else if c == '2' || c == '4' {
                   if exp != 'R' {
                       ok = false
                       break
                   }
               }
           }
           if ok {
               ways++
           }
       }
       if ways == 0 {
           fmt.Println(0)
           return
       }
       total = (total * ways) % MOD
   }
   // Vertical: columns
   for j := 0; j < m; j++ {
       ways := 0
       for p := 0; p < 2; p++ {
           ok := true
           for i := 0; i < n; i++ {
               c := grid[i][j]
               var exp byte
               if (i%2 == 0) == (p == 0) {
                   exp = 'T'
               } else {
                   exp = 'B'
               }
               if c == '1' || c == '2' {
                   if exp != 'T' {
                       ok = false
                       break
                   }
               } else if c == '3' || c == '4' {
                   if exp != 'B' {
                       ok = false
                       break
                   }
               }
           }
           if ok {
               ways++
           }
       }
       if ways == 0 {
           fmt.Println(0)
           return
       }
       total = (total * ways) % MOD
   }
   fmt.Println(total)

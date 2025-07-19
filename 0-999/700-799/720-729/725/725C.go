package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var s string
   if _, err := fmt.Fscan(reader, &s); err != nil {
       return
   }
   // Check consecutive duplicates
   for i := 0; i < len(s)-1; i++ {
       if s[i] == s[i+1] {
           fmt.Println("Impossible")
           return
       }
   }
   // Count and find duplicate letter
   cnt := make([]int, 26)
   var dup byte
   for i := 0; i < len(s); i++ {
       idx := s[i] - 'A'
       cnt[idx]++
       if cnt[idx] == 2 {
           dup = s[i]
       }
   }
   // Positions of duplicate
   fst, snd := -1, -1
   for i := 0; i < len(s); i++ {
       if s[i] == dup {
           if fst < 0 {
               fst = i
           } else {
               snd = i
           }
       }
   }
   // Prepare board
   const cols = 13
   board := [2][cols]rune{}
   // Distance between duplicates
   dist := snd - fst - 1
   // Starting column for duplicate in row 0
   start := 12 - (dist / 2)
   board[0][start] = rune(dup)
   // Fill between duplicates
   flag := false
   ptr := start + 1
   for i := fst + 1; i < snd; i++ {
       if ptr == cols {
           ptr--
           flag = true
       }
       if !flag {
           board[0][ptr] = rune(s[i]); ptr++
       } else {
           board[1][ptr] = rune(s[i]); ptr--
       }
   }
   // Fill before first duplicate
   flag = false
   ptr = start - 1
   for i := fst - 1; i >= 0; i-- {
       if ptr < 0 {
           ptr++
           flag = true
       }
       if !flag {
           board[0][ptr] = rune(s[i]); ptr--
       } else {
           board[1][ptr] = rune(s[i]); ptr++
       }
   }
   // Fill after second duplicate
   ptr = start
   if dist%2 == 1 {
       ptr--
   }
   flag = true
   for i := snd + 1; i < len(s); i++ {
       if ptr < 0 {
           ptr++
           flag = false
       }
       if flag {
           board[1][ptr] = rune(s[i]); ptr--
       } else {
           board[0][ptr] = rune(s[i]); ptr++
       }
   }
   // Output
   w := bufio.NewWriter(os.Stdout)
   defer w.Flush()
   for r := 0; r < 2; r++ {
       for c := 0; c < cols; c++ {
           if board[r][c] == 0 {
               w.WriteRune(' ')
           } else {
               w.WriteRune(board[r][c])
           }
       }
       w.WriteByte('\n')
   }
}

package main

import (
   "bufio"
   "fmt"
   "os"
)

func abs(x int) int {
   if x < 0 {
       return -x
   }
   return x
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m, k int
   if _, err := fmt.Fscan(reader, &n, &m, &k); err != nil {
       return
   }
   draw := false
   rows := make([][]byte, n)
   for i := 0; i < n; i++ {
       line := make([]byte, 0, m)
       for len(line) < m {
           buf := make([]byte, m)
           cnt, err := reader.Read(buf)
           if err != nil {
               break
           }
           for j := 0; j < cnt && len(line) < m; j++ {
               if buf[j] != '\n' && buf[j] != '\r' {
                   line = append(line, buf[j])
               }
           }
       }
       rows[i] = line
       cnt := 0
       for _, c := range line {
           if c == 'G' || c == 'R' {
               cnt++
           }
       }
       if cnt <= 1 {
           draw = true
       }
   }
   if draw || k >= 2 {
       fmt.Println("Draw")
       return
   }
   // k == 1 and all rows have exactly 2 soldiers
   xor := 0
   for i := 0; i < n; i++ {
       var gpos, rpos int = -1, -1
       for j, c := range rows[i] {
           if c == 'G' {
               gpos = j
           } else if c == 'R' {
               rpos = j
           }
       }
       if gpos >= 0 && rpos >= 0 {
           d := abs(gpos - rpos) - 1
           xor ^= d
       }
   }
   if xor != 0 {
       fmt.Println("First")
   } else {
       fmt.Println("Second")
   }

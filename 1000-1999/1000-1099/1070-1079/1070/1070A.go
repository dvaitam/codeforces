package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   var d, s int
   if _, err := fmt.Fscan(in, &d, &s); err != nil {
       return
   }
   // visited[rem][sum]
   visited := make([][]bool, d)
   for i := 0; i < d; i++ {
       visited[i] = make([]bool, s+1)
   }
   // parent pointers
   parentRem := make([][]int16, d)
   parentSum := make([][]int16, d)
   parentDig := make([][]byte, d)
   for i := 0; i < d; i++ {
       parentRem[i] = make([]int16, s+1)
       parentSum[i] = make([]int16, s+1)
       parentDig[i] = make([]byte, s+1)
   }
   type state struct{ r, sum int16 }
   // BFS queue
   q := make([]state, 0, d*min(s+1, 10))
   // initialize with first digit 1..9
   for dig := 1; dig <= 9 && dig <= s; dig++ {
       r := dig % d
       sm := dig
       if !visited[r][sm] {
           visited[r][sm] = true
           parentRem[r][sm] = -1
           parentSum[r][sm] = -1
           parentDig[r][sm] = byte(dig)
           q = append(q, state{int16(r), int16(sm)})
       }
   }
   foundR, foundSum := int16(-1), int16(-1)
   head := 0
   for head < len(q) {
       cur := q[head]
       head++
       if cur.r == 0 && int(cur.sum) == s {
           foundR, foundSum = cur.r, cur.sum
           break
       }
       for dig := 0; dig <= 9; dig++ {
           newSum := int(cur.sum) + dig
           if newSum > s {
               break
           }
           newR := (int(cur.r)*10 + dig) % d
           if !visited[newR][newSum] {
               visited[newR][newSum] = true
               parentRem[newR][newSum] = cur.r
               parentSum[newR][newSum] = cur.sum
               parentDig[newR][newSum] = byte(dig)
               q = append(q, state{int16(newR), int16(newSum)})
           }
       }
   }
   if foundR < 0 {
       fmt.Fprintln(out, -1)
       return
   }
   // reconstruct number
   var digs []byte
   r := foundR
   sm := foundSum
   for {
       dgt := parentDig[r][sm]
       digs = append(digs, '0'+dgt)
       pr := parentRem[r][sm]
       ps := parentSum[r][sm]
       if ps < 0 {
           break
       }
       r, sm = pr, ps
   }
   // reverse digits
   for i, j := 0, len(digs)-1; i < j; i, j = i+1, j-1 {
       digs[i], digs[j] = digs[j], digs[i]
   }
   out.Write(digs)
}

func min(a, b int) int {
   if a < b {
       return a
   }
   return b
}

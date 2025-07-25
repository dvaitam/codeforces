package main

import (
   "bufio"
   "fmt"
   "os"
)

type state struct {
   carry int
   diff  int
   prev  int
   digit int
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var a int
   if _, err := fmt.Fscan(reader, &a); err != nil {
       return
   }
   // BFS on state (carry, diff)
   // diff range limited to [-maxDiff, maxDiff]
   maxCarry := a
   maxDiff := 9 * a
   diffRange := 2*maxDiff + 1
   diffOffset := maxDiff
   // visited[carry][diff+offset]
   visited := make([][]bool, maxCarry+1)
   for i := range visited {
       visited[i] = make([]bool, diffRange)
   }
   // queue of states
   var q []state
   // initial state
   q = append(q, state{carry: 0, diff: 0, prev: -1, digit: -1})
   visited[0][diffOffset] = true
   foundIdx := -1
   for idx := 0; idx < len(q); idx++ {
       st := q[idx]
       // skip initial idx=0 for found
       if idx > 0 && st.carry == 0 && st.diff == 0 {
           foundIdx = idx
           break
       }
       // try digits 0..9
       for d := 0; d <= 9; d++ {
           // skip leading zero block (d=0 at first step)
           if idx == 0 && d == 0 {
               continue
           }
           tmp := a*d + st.carry
           r := tmp % 10
           nc := tmp / 10
           nd := st.diff + d - a*r
           if nd < -maxDiff || nd > maxDiff || nc < 0 || nc > maxCarry {
               continue
           }
           off := nd + diffOffset
           if visited[nc][off] {
               continue
           }
           visited[nc][off] = true
           q = append(q, state{carry: nc, diff: nd, prev: idx, digit: d})
           // next iteration will catch carry=0,diff=0
       }
   }
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   if foundIdx < 0 {
       fmt.Fprintln(writer, -1)
       return
   }
   // reconstruct digits
   var digs []int
   for i := foundIdx; i > 0; i = q[i].prev {
       digs = append(digs, q[i].digit)
   }
   // reverse digs
   for i, j := 0, len(digs)-1; i < j; i, j = i+1, j-1 {
       digs[i], digs[j] = digs[j], digs[i]
   }
   if len(digs) == 0 {
       fmt.Fprintln(writer, -1)
       return
   }
   // print number
   for _, d := range digs {
       writer.WriteByte(byte('0' + d))
   }
   writer.WriteByte('\n')
}

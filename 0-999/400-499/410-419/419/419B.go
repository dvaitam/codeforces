package main

import (
   "bufio"
   "fmt"
   "os"
)

type evt struct {
   t  int
   op byte
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   // store events by user and global ops
   events := make([][]evt, n+1)
   firstOp := make([]byte, n+1)
   ops := make([]byte, m+1)

   for i := 1; i <= m; i++ {
       var op string
       var id int
       fmt.Fscan(reader, &op, &id)
       b := op[0]
       ops[i] = b
       if firstOp[id] == 0 {
           firstOp[id] = b
       }
       events[id] = append(events[id], evt{t: i, op: b})
   }
   // initial online status
   initialOnline := make([]bool, n+1)
   initCount := 0
   for id := 1; id <= n; id++ {
       if firstOp[id] == '-' {
           initialOnline[id] = true
           initCount++
       }
   }
   // cur and bad arrays
   cur := make([]int, m+1)
   bad := make([]int, m+1)
   prefBad := make([]int, m+1)
   cur[0] = initCount
   bad[0] = b2i(cur[0] > 0)
   prefBad[0] = bad[0]
   for i := 1; i <= m; i++ {
       if ops[i] == '+' {
           cur[i] = cur[i-1] + 1
       } else {
           cur[i] = cur[i-1] - 1
       }
       if cur[i] > 0 {
           bad[i] = 1
       }
       prefBad[i] = prefBad[i-1] + bad[i]
   }
   // check candidates
   res := make([]int, 0, n)
   for id := 1; id <= n; id++ {
       userEvents := events[id]
       lastT := 0
       online := initialOnline[id]
       possible := true
       for _, e := range userEvents {
           if !online {
               l, r := lastT, e.t-1
               if l <= r && hasBad(prefBad, l, r) {
                   possible = false
                   break
               }
           }
           // update status at event
           if e.op == '+' {
               online = true
           } else {
               online = false
           }
           lastT = e.t
       }
       if possible && !online {
           if lastT <= m && hasBad(prefBad, lastT, m) {
               possible = false
           }
       }
       if possible {
           res = append(res, id)
       }
   }
   // output
   fmt.Fprintln(writer, len(res))
   if len(res) > 0 {
       for i, id := range res {
           if i > 0 {
               writer.WriteByte(' ')
           }
           fmt.Fprint(writer, id)
       }
       writer.WriteByte('\n')
   }
}

func b2i(b bool) int {
   if b {
       return 1
   }
   return 0
}

// hasBad reports if prefBad[l..r] has any positive entries
func hasBad(prefBad []int, l, r int) bool {
   if l == 0 {
       return prefBad[r] > 0
   }
   return prefBad[r]-prefBad[l-1] > 0
}

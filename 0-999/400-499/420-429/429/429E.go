package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
   "strconv"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   // events: position and event id
   ev := make([]struct{ pos, id int }, 2*n)
   for i := 0; i < n; i++ {
       var a, b int
       fmt.Fscan(in, &a, &b)
       ev[2*i] = struct{ pos, id int }{a, 2 * i}
       ev[2*i+1] = struct{ pos, id int }{b + 1, 2*i + 1}
   }
   // sort events by position, then id
   sort.Slice(ev, func(i, j int) bool {
       if ev[i].pos != ev[j].pos {
           return ev[i].pos < ev[j].pos
       }
       return ev[i].id < ev[j].id
   })
   // build matching p between events
   p := make([]int, 2*n)
   for i := 0; i < n; i++ {
       a := ev[2*i].id
       b := ev[2*i+1].id
       p[a] = b
       p[b] = a
   }
   // c: 0=unvisited, 2 or 3 assigned
   c := make([]int, n)
   // iterative DFS stack of (segment, t)
   stackX := make([]int, 0, n)
   stackT := make([]int, 0, n)
   for i := 0; i < n; i++ {
       if c[i] != 0 {
           continue
       }
       // start new component
       stackX = append(stackX, i)
       stackT = append(stackT, 2)
       for len(stackX) > 0 {
           // pop
           idx := len(stackX) - 1
           x := stackX[idx]
           t := stackT[idx]
           stackX = stackX[:idx]
           stackT = stackT[:idx]
           if c[x] != 0 {
               continue
           }
           c[x] = t
           // neighbor from event 2*x
           u0 := p[2*x]
           y0 := u0 >> 1
           t0 := (u0 & 1) ^ t ^ 1
           if c[y0] == 0 {
               stackX = append(stackX, y0)
               stackT = append(stackT, t0)
           }
           // neighbor from event 2*x+1
           u1 := p[2*x+1]
           y1 := u1 >> 1
           t1 := (u1 & 1) ^ t
           if c[y1] == 0 {
               stackX = append(stackX, y1)
               stackT = append(stackT, t1)
           }
       }
   }
   // output: c[i] ^ 2 gives 0 or 1
   for i := 0; i < n; i++ {
       if i > 0 {
           out.WriteByte(' ')
       }
       // 2^2=0, 3^2=1
       res := c[i] ^ 2
       out.WriteString(strconv.Itoa(res))
   }
}

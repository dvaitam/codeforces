package main

import (
   "bufio"
   "fmt"
   "os"
)

type ctx struct {
   cond       int
   entryCost  int64
   innerCost  int64
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, s int
   if _, err := fmt.Fscan(reader, &n, &s); err != nil {
       return
   }
   const maxY = 200000
   costToRemove := make([]int64, maxY+1)
   // stack of contexts
   stack := make([]ctx, 0, 100)
   // root context
   stack = append(stack, ctx{cond: -1, entryCost: 0, innerCost: 0})
   for i := 0; i < n; i++ {
       var tok string
       fmt.Fscan(reader, &tok)
       if tok == "set" {
           var y int
           var v int64
           fmt.Fscan(reader, &y, &v)
           // record cost to remove this set
           costToRemove[y] += v
           if y == s {
               // must pay v if this set potentially runs
               idx := len(stack) - 1
               stack[idx].innerCost += v
           }
       } else if tok == "if" {
           var y int
           fmt.Fscan(reader, &y)
           // push new context for this block
           c := ctx{cond: y, entryCost: costToRemove[y], innerCost: 0}
           stack = append(stack, c)
       } else if tok == "end" {
           // pop context
           last := stack[len(stack)-1]
           stack = stack[:len(stack)-1]
           // process last into parent
           if last.cond == s {
               // block never executes since x!=s always, ignore innerCost
               continue
           }
           // cost to avoid entering block: remove all sets assigning x=cond before block
           avoidCost := last.entryCost
           // cost to remove s-sets inside if
           payCost := last.innerCost
           // choose minimal
           var add int64
           if avoidCost < payCost {
               add = avoidCost
           } else {
               add = payCost
           }
           // add to parent innerCost
           idx := len(stack) - 1
           stack[idx].innerCost += add
       }
   }
   // result in root context
   result := stack[0].innerCost
   fmt.Println(result)
}

package main

import (
   "bufio"
   "fmt"
   "os"
)

type node struct {
   left, right, sum int
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var h, q int
   fmt.Fscan(reader, &h)
   fmt.Fscan(reader, &q)
   // Preallocate nodes; use 1-based indexing
   maxNodes := q * (h + 2) + 10
   nodes := make([]node, maxNodes)
   lastnode := 1
   // path buffer
   path := make([]int, h+2)
   for i := 0; i < q; i++ {
       var op string
       fmt.Fscan(reader, &op)
       if op[0] == 'a' {
           var val, extra int
           fmt.Fscan(reader, &val, &extra)
           // build path to root
           plen := 0
           for {
               path[plen] = val
               plen++
               if val == 1 {
                   break
               }
               val /= 2
           }
           at := 1
           // traverse down path, skipping root duplicate at end
           for j := plen - 2; j >= 0; j-- {
               nodes[at].sum += extra
               v := path[j]
               if v%2 == 0 {
                   if nodes[at].left == 0 {
                       lastnode++
                       nodes[at].left = lastnode
                   }
                   at = nodes[at].left
               } else {
                   if nodes[at].right == 0 {
                       lastnode++
                       nodes[at].right = lastnode
                   }
                   at = nodes[at].right
               }
           }
           nodes[at].sum += extra
       } else {
           // query
           var res, pr float64
           pr = 1.0
           at := 1
           maxOther := 0
           for at != 0 {
               pr *= 0.5
               leftSum := 0
               if nodes[at].left != 0 {
                   leftSum = nodes[nodes[at].left].sum
               }
               rightSum := 0
               if nodes[at].right != 0 {
                   rightSum = nodes[nodes[at].right].sum
               }
               cur := nodes[at].sum - leftSum - rightSum
               // update result with weighted max
               stop := max(leftSum, rightSum) + cur
               res += pr * float64(max(maxOther, stop))
               // update best alternative sum
               maxOther = max(maxOther, min(leftSum, rightSum)+cur)
               // move to the heavier child
               if leftSum > rightSum {
                   at = nodes[at].left
               } else {
                   at = nodes[at].right
               }
           }
           res += pr * float64(maxOther)
           fmt.Printf("%.10f\n", res)
       }
   }
}

func max(a, b int) int {
   if a > b {
       return a
   }
   return b
}

func min(a, b int) int {
   if a < b {
       return a
   }
   return b
}

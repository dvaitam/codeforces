package main

import (
   "bufio"
   "fmt"
   "os"
)

type Prev struct {
   prev int
   ch   byte
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var a, b string
   if _, err := fmt.Fscan(reader, &a, &b); err != nil {
       return
   }
   la, lb := len(a), len(b)
   maxN := la + lb + 2
   strideB := maxN + 1
   strideJ := (lb + 1) * strideB
   // encode state id = i*strideJ + j*strideB + b
   startID := 0
   targetID := la*strideJ + lb*strideB + 0
   // BFS
   queue := make([]int, 0, 1024)
   queue = append(queue, startID)
   visited := make(map[int]bool)
   visited[startID] = true
   prevMap := make(map[int]Prev)
   prevMap[startID] = Prev{-1, 0}
   var head int
   found := false
   for head < len(queue) && !found {
       id := queue[head]
       head++
       i := id / strideJ
       rem := id - i*strideJ
       j := rem / strideB
       bbal := rem - j*strideB
       for _, ch := range []byte{'(', ')'} {
           nb := bbal
           if ch == '(' {
               nb = bbal + 1
           } else {
               if bbal == 0 {
                   continue
               }
               nb = bbal - 1
           }
           // advance i and j for subsequence
           ni := i
           if i < la && a[i] == ch {
               ni = i + 1
           }
           nj := j
           if j < lb && b[j] == ch {
               nj = j + 1
           }
           nid := ni*strideJ + nj*strideB + nb
           if !visited[nid] {
               visited[nid] = true
               prevMap[nid] = Prev{id, ch}
               queue = append(queue, nid)
               if nid == targetID {
                   found = true
                   break
               }
           }
       }
   }
   // reconstruct
   res := make([]byte, 0, 256)
   cur := targetID
   for cur != startID {
       p := prevMap[cur]
       res = append(res, p.ch)
       cur = p.prev
   }
   // reverse
   for l, r := 0, len(res)-1; l < r; l, r = l+1, r-1 {
       res[l], res[r] = res[r], res[l]
   }
   fmt.Fprint(os.Stdout, string(res))
}

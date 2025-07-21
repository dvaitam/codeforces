package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var s string
   fmt.Fscan(reader, &s)
   // color indices: R,O,Y,G,B,V -> 0..5
   colorMap := map[rune]int{
       'R': 0, 'O': 1, 'Y': 2,
       'G': 3, 'B': 4, 'V': 5,
   }
   var avail [6]int
   for _, ch := range s {
       if idx, ok := colorMap[ch]; ok {
           avail[idx]++
       }
   }
   // define basic rotations Rx, Ry, Rz
   // face order: 0=U,1=F,2=R,3=B,4=L,5=D
   Rx := []int{1, 5, 2, 0, 4, 3}
   Ry := []int{2, 1, 5, 3, 0, 4}
   Rz := []int{0, 2, 3, 4, 1, 5}
   // generate group
   perms := make([][]int, 0, 24)
   seen := make(map[string]bool)
   // identity
   id := []int{0, 1, 2, 3, 4, 5}
   queue := [][]int{id}
   seen[key(id)] = true
   perms = append(perms, id)
   gens := [][]int{Rx, Ry, Rz}
   for len(queue) > 0 {
       cur := queue[0]
       queue = queue[1:]
       for _, g := range gens {
           nxt := compose(cur, g)
           k := key(nxt)
           if !seen[k] {
               seen[k] = true
               perms = append(perms, nxt)
               queue = append(queue, nxt)
           }
       }
   }
   // should have 24 perms
   total := 0
   for _, perm := range perms {
       // get cycles
       vis := [6]bool{}
       var cycles []int
       for i := 0; i < 6; i++ {
           if !vis[i] {
               // trace cycle
               j := i
               cnt := 0
               for !vis[j] {
                   vis[j] = true
                   cnt++
                   j = perm[j]
               }
               if cnt > 0 {
                   cycles = append(cycles, cnt)
               }
           }
       }
       // count assignments
       used := [6]int{}
       total += dfsAssign(cycles, avail, used, 0)
   }
   // divide by group size
   fmt.Println(total / len(perms))
}

// compose p following by g: apply p then g
func compose(p, g []int) []int {
   r := make([]int, len(p))
   for i := range p {
       r[i] = g[p[i]]
   }
   return r
}

func key(p []int) string {
   b := make([]byte, len(p))
   for i, v := range p {
       b[i] = byte(v)
   }
   return string(b)
}

// dfsAssign assigns colors to cycles recursively
func dfsAssign(cycles []int, avail, used [6]int, idx int) int {
   if idx == len(cycles) {
       // check if used == avail
       for i := 0; i < 6; i++ {
           if used[i] != avail[i] {
               return 0
           }
       }
       return 1
   }
   cnt := 0
   l := cycles[idx]
   for c := 0; c < 6; c++ {
       if used[c]+l <= avail[c] {
           used[c] += l
           cnt += dfsAssign(cycles, avail, used, idx+1)
           used[c] -= l
       }
   }
   return cnt
}

package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var s, t string
   if _, err := fmt.Fscanln(reader, &s); err != nil {
       return
   }
   if _, err := fmt.Fscanln(reader, &t); err != nil {
       return
   }
   n, m := len(s), len(t)
   // BFS over reachable states
   type pair struct{ i, j int }
   // visited keyed by int64 (i<<32 | j)
   visited := make(map[int64]bool)
   queue := make([]pair, 0, 1024)
   enqueue := func(i, j int) {
       key := int64(i)<<32 | int64(j)
       if !visited[key] {
           visited[key] = true
           queue = append(queue, pair{i, j})
       }
   }
   enqueue(0, 0)
   colors := []byte{'R', 'G', 'B'}
   for head := 0; head < len(queue); head++ {
       p := queue[head]
       i, j := p.i, p.j
       for _, c := range colors {
           // cannot perform if animal at last with matching color
           if i < n && s[i] == c {
               if i == n-1 {
                   continue
               }
           }
           if j < m && t[j] == c {
               if j == m-1 {
                   continue
               }
           }
           ni, nj := i, j
           if i < n && s[i] == c {
               ni++
           }
           if j < m && t[j] == c {
               nj++
           }
           enqueue(ni, nj)
       }
   }
   fmt.Println(len(visited))
}

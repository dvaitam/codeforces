package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var m, k int
   if _, err := fmt.Fscan(reader, &m, &k); err != nil {
       return
   }
   // Read friend pairs and collect unique user IDs
   pairs := make([][2]int64, m)
   idSet := make(map[int64]struct{}, 2*m)
   ids := make([]int64, 0, 2*m)
   for i := 0; i < m; i++ {
       var a, b int64
       fmt.Fscan(reader, &a, &b)
       pairs[i] = [2]int64{a, b}
       if _, ok := idSet[a]; !ok {
           idSet[a] = struct{}{}
           ids = append(ids, a)
       }
       if _, ok := idSet[b]; !ok {
           idSet[b] = struct{}{}
           ids = append(ids, b)
       }
   }
   // Sort user IDs
   sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })
   n := len(ids)
   // Map IDs to indices
   id2idx := make(map[int64]int, n)
   for i, id := range ids {
       id2idx[id] = i
   }
   // Build adjacency sets
   friends := make([]map[int]struct{}, n)
   for i := 0; i < n; i++ {
       friends[i] = make(map[int]struct{})
   }
   for _, p := range pairs {
       u, v := id2idx[p[0]], id2idx[p[1]]
       friends[u][v] = struct{}{}
       friends[v][u] = struct{}{}
   }
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   // For each user, find suggested friends
   for xi, xID := range ids {
       d := len(friends[xi])
       // threshold = ceil(k * d / 100)
       threshold := (k*d + 99) / 100
       // Collect suggestions in sorted order
       suggestions := make([]int64, 0)
       for yi, yID := range ids {
           if yi == xi {
               continue
           }
           if _, ok := friends[xi][yi]; ok {
               continue
           }
           // Count common friends
           cnt := 0
           for f := range friends[xi] {
               if _, ok := friends[yi][f]; ok {
                   cnt++
               }
           }
           if cnt >= threshold {
               suggestions = append(suggestions, yID)
           }
       }
       // Output
       fmt.Fprintf(writer, "%d: %d", xID, len(suggestions))
       for _, yID := range suggestions {
           fmt.Fprintf(writer, " %d", yID)
       }
   fmt.Fprintln(writer)
   }
}

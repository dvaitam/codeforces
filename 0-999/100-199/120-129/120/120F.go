package main

import (
   "bufio"
   "fmt"
   "os"
   "strconv"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   // read integer n
   line, _, err := reader.ReadLine()
   if err != nil {
       panic(err)
   }
   n, err := strconv.Atoi(string(line))
   if err != nil {
       panic(err)
   }
   totalDiameter := 0
   for i := 0; i < n; i++ {
       // read next line for ni and edges
       line, _, err := reader.ReadLine()
       if err != nil {
           panic(err)
       }
       parts := splitFields(string(line))
       if len(parts) < 1 {
           panic("invalid input")
       }
       ni, err := strconv.Atoi(parts[0])
       if err != nil {
           panic(err)
       }
       // adjacency list
       adj := make([][]int, ni)
       // parts length should be 1 + 2*(ni-1)
       expected := 1 + 2*(ni-1)
       if len(parts) != expected {
           // maybe edges wrap to next lines, read rest
           // read ints until have expected
           vals := make([]string, 0, expected)
           vals = append(vals, parts...)
           for len(vals) < expected {
               line, _, err = reader.ReadLine()
               if err != nil {
                   panic(err)
               }
               more := splitFields(string(line))
               vals = append(vals, more...)
           }
           parts = vals
       }
       // parse edges
       for j := 0; j < ni-1; j++ {
           u, err1 := strconv.Atoi(parts[1+2*j])
           v, err2 := strconv.Atoi(parts[1+2*j+1])
           if err1 != nil || err2 != nil {
               panic("invalid edge")
           }
           u--
           v--
           adj[u] = append(adj[u], v)
           adj[v] = append(adj[v], u)
       }
       // compute diameter by two BFS
       // first BFS from node 0
       far, _ := bfsFarthest(0, adj)
       // second BFS from farthest
       _, dist := bfsFarthest(far, adj)
       totalDiameter += dist
   }
   fmt.Println(totalDiameter)
}

// splitFields splits a string by spaces
func splitFields(s string) []string {
   var res []string
   field := ""
   for i := 0; i < len(s); i++ {
       c := s[i]
       if c == ' ' || c == '\t' {
           if field != "" {
               res = append(res, field)
               field = ""
           }
       } else {
           field += string(c)
       }
   }
   if field != "" {
       res = append(res, field)
   }
   return res
}

// bfsFarthest returns the farthest node and its distance from start
func bfsFarthest(start int, adj [][]int) (int, int) {
   n := len(adj)
   dist := make([]int, n)
   for i := range dist {
       dist[i] = -1
   }
   queue := make([]int, 0, n)
   dist[start] = 0
   queue = append(queue, start)
   qhead := 0
   var farNode, farDist int
   for qhead < len(queue) {
       u := queue[qhead]
       qhead++
       for _, v := range adj[u] {
           if dist[v] == -1 {
               dist[v] = dist[u] + 1
               queue = append(queue, v)
           }
       }
       if dist[u] > farDist {
           farDist = dist[u]
           farNode = u
       }
   }
   return farNode, farDist
}

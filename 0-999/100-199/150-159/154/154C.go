package main

import (
   "bufio"
   "fmt"
   "io"
   "math/rand"
   "os"
   "sort"
   "time"
)

func main() {
   br := bufio.NewReader(os.Stdin)
   // fast integer reader
   readInt := func() (int, error) {
       var c byte
       var err error
       // skip spaces
       for {
           c, err = br.ReadByte()
           if err != nil {
               return 0, err
           }
           if c > ' ' {
               break
           }
       }
       sign := 1
       if c == '-' {
           sign = -1
           c, err = br.ReadByte()
           if err != nil {
               return 0, err
           }
       }
       x := 0
       for c >= '0' && c <= '9' {
           x = x*10 + int(c-'0')
           c, err = br.ReadByte()
           if err != nil && err != io.EOF {
               return 0, err
           }
           if err == io.EOF {
               break
           }
       }
       return sign * x, nil
   }
   n, err := readInt()
   if err != nil {
       return
   }
   m, err := readInt()
   if err != nil {
       return
   }
   // random values
   z := make([]uint64, n+1)
   rng := rand.New(rand.NewSource(time.Now().UnixNano()))
   for i := 1; i <= n; i++ {
       z[i] = rng.Uint64()
   }
   h := make([]uint64, n+1)
   edges := make([][2]int, m)
   for i := 0; i < m; i++ {
       u, err := readInt()
       if err != nil {
           return
       }
       v, err := readInt()
       if err != nil {
           return
       }
       edges[i][0], edges[i][1] = u, v
       h[u] += z[v]
       h[v] += z[u]
   }
   // count non-friend doubles (equal hashes)
   hs := make([]uint64, n)
   for i := 1; i <= n; i++ {
       hs[i-1] = h[i]
   }
   sort.Slice(hs, func(i, j int) bool { return hs[i] < hs[j] })
   var ans int64
   for i := 0; i < n; {
       j := i + 1
       for j < n && hs[j] == hs[i] {
           j++
       }
       c := int64(j - i)
       ans += c * (c - 1) / 2
       i = j
   }
   // count friend doubles
   for _, e := range edges {
       u, v := e[0], e[1]
       if h[u]-z[v] == h[v]-z[u] {
           ans++
       }
   }
   // output
   fmt.Println(ans)
}

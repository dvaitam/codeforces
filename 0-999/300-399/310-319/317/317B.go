package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   var n, t int
   if _, err := fmt.Fscan(in, &n, &t); err != nil {
       return
   }
   // sandpile on infinite grid
   m := make(map[int64]int)
   // pack coordinates into key
   pack := func(x, y int) int64 {
       return (int64(x) << 32) | int64(uint32(y))
   }
   unpack := func(key int64) (int, int) {
       x := int(key >> 32)
       y := int(int32(key))
       return x, y
   }
   origin := pack(0, 0)
   m[origin] = n
   queue := make([]int64, 0, n)
   if n >= 4 {
       queue = append(queue, origin)
   }
   for i := 0; i < len(queue); i++ {
       key := queue[i]
       cnt := m[key]
       k := cnt / 4
       if k == 0 {
           continue
       }
       m[key] = cnt - k*4
       x, y := unpack(key)
       for _, d := range [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}} {
           nx, ny := x+d[0], y+d[1]
           nkey := pack(nx, ny)
           prev := m[nkey]
           m[nkey] = prev + k
           if prev < 4 && m[nkey] >= 4 {
               queue = append(queue, nkey)
           }
       }
   }
   // answer queries
   for i := 0; i < t; i++ {
       var xi, yi int
       fmt.Fscan(in, &xi, &yi)
       key := pack(xi, yi)
       if v, ok := m[key]; ok {
           fmt.Fprintln(out, v)
       } else {
           fmt.Fprintln(out, 0)
       }
   }
}

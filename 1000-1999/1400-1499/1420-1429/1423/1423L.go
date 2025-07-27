package main

import (
   "bufio"
   "fmt"
   "os"
   "strconv"
   "strings"
   "math/bits"
)
const ML = 16

func main() {
   reader := bufio.NewReader(os.Stdin)
   line, _ := reader.ReadString('\n')
   parts := strings.Fields(line)
   N, _ := strconv.Atoi(parts[0])
   S, _ := strconv.Atoi(parts[1])
   D, _ := strconv.Atoi(parts[2])
   // number of uint64 needed
   L := (N + 63) >> 6
   // read switches
   // read switches
   switchMasks := make([][ML]uint64, S)
   for i := 0; i < S; i++ {
       line, _ = reader.ReadString('\n')
       parts = strings.Fields(line)
       ci, _ := strconv.Atoi(parts[0])
       var mask [ML]uint64
       for j := 0; j < ci; j++ {
           idx, _ := strconv.Atoi(parts[1+j])
           idx--
           mask[idx>>6] ^= 1 << (idx & 63)
       }
       switchMasks[i] = mask
   }
   // split switches
   s1 := S / 2
   s2 := S - s1
   len1 := 1 << s1
   len2 := 1 << s2
   // precompute half1
   bm1 := make([][ML]uint64, len1)
   cost1 := make([]int, len1)
   // bm1[0] is zero
   for m := 1; m < len1; m++ {
       lb := m & -m
       p := bits.TrailingZeros(uint(lb))
       prev := m ^ lb
       prevbm := bm1[prev]
       var cur [ML]uint64
       // copy and apply switch p
       for j := 0; j < L; j++ {
           cur[j] = prevbm[j] ^ switchMasks[p][j]
       }
       // rest of cur[j] are zero by default
       bm1[m] = cur
       cost1[m] = cost1[prev] + 1
   }
   // precompute half2 slice and cost
   bm2 := make([][ML]uint64, len2)
   cost2slice := make([]int, len2)
   // bm2[0] is zero
   for m := 1; m < len2; m++ {
       lb := m & -m
       p := bits.TrailingZeros(uint(lb))
       prev := m ^ lb
       prevbm := bm2[prev]
       var cur [ML]uint64
       for j := 0; j < L; j++ {
           cur[j] = prevbm[j] ^ switchMasks[s1+p][j]
       }
       bm2[m] = cur
       cost2slice[m] = cost2slice[prev] + 1
   }
   // build map2 for half2
   // build map2 for half2
   map2 := make(map[[ML]uint64]int, len2)
   for m := 0; m < len2; m++ {
       key := bm2[m]
       c := cost2slice[m]
       if prev, ok := map2[key]; !ok || c < prev {
           map2[key] = c
       }
   }
   // process days
   out := bufio.NewWriter(os.Stdout)
   for d := 0; d < D; d++ {
       line, _ = reader.ReadString('\n')
       parts = strings.Fields(line)
       ti, _ := strconv.Atoi(parts[0])
       var bmask [ML]uint64
       for j := 0; j < ti; j++ {
           idx, _ := strconv.Atoi(parts[1+j])
           idx--
           bmask[idx>>6] ^= 1 << (idx & 63)
       }
       // search
       best := -1
       for i := 0; i < len1; i++ {
           // need mask2 = bm1[i] xor bmask
           var key [ML]uint64
           for j := 0; j < L; j++ {
               key[j] = bm1[i][j] ^ bmask[j]
           }
           if c2, ok := map2[key]; ok {
               c := cost1[i] + c2
               if best < 0 || c < best {
                   best = c
               }
           }
       }
       fmt.Fprintln(out, best)
   }
   out.Flush()
}

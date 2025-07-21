package main

import (
   "bufio"
   "fmt"
   "math/bits"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   var sa, sb string
   fmt.Fscan(reader, &sa)
   fmt.Fscan(reader, &sb)
   na, nb := len(sa), len(sb)
   // prefix sums of ones
   preA := make([]int, na+1)
   preB := make([]int, nb+1)
   for i := 0; i < na; i++ {
       preA[i+1] = preA[i]
       if sa[i] == '1' {
           preA[i+1]++
       }
   }
   for i := 0; i < nb; i++ {
       preB[i+1] = preB[i]
       if sb[i] == '1' {
           preB[i+1]++
       }
   }
   // build blocks
   const W = 64
   naBlk := (na + W - 1) / W
   nbBlk := (nb + W - 1) / W
   aBlk := make([]uint64, naBlk+1)
   bBlk := make([]uint64, nbBlk+1)
   for i := 0; i < na; i++ {
       if sa[i] == '1' {
           aBlk[i/W] |= 1 << (uint(i) & (W - 1))
       }
   }
   for i := 0; i < nb; i++ {
       if sb[i] == '1' {
           bBlk[i/W] |= 1 << (uint(i) & (W - 1))
       }
   }
   var q int
   fmt.Fscan(reader, &q)
   for qi := 0; qi < q; qi++ {
       var p1, p2, ln int
       fmt.Fscan(reader, &p1, &p2, &ln)
       // count ones in a and b
       sa1 := preA[p1+ln] - preA[p1]
       sb1 := preB[p2+ln] - preB[p2]
       // count intersection
       inter := countAnd(aBlk, bBlk, p1, p2, ln)
       // hamming distance = onesA + onesB - 2*inter
       fmt.Fprint(writer, sa1+sb1-2*inter)
       if qi+1 < q {
           writer.WriteByte(' ')
       }
   }
   writer.WriteByte('\n')
}

func countAnd(aBlk, bBlk []uint64, p1, p2, length int) int {
   const W = 64
   wa := p1 / W
   ba := p1 & (W - 1)
   wb := p2 / W
   bb := p2 & (W - 1)
   // number of 64-bit words needed
   words := (ba + length + W - 1) / W
   var cnt int
   for i := 0; i < words; i++ {
       // extract word from a
       aw := aBlk[wa+i] >> uint(ba)
       if ba != 0 {
           aw |= aBlk[wa+i+1] << uint(W-ba)
       }
       // extract word from b
       bw := bBlk[wb+i] >> uint(bb)
       if bb != 0 {
           bw |= bBlk[wb+i+1] << uint(W-bb)
       }
       // mask last word
       if i == words-1 {
           r := (ba + length) & (W - 1)
           if r != 0 {
               mask := uint64((1<<uint(r)) - 1)
               aw &= mask
               bw &= mask
           }
       }
       cnt += bits.OnesCount64(aw & bw)
   }
   return cnt
}

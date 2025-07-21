package main

import (
   "bufio"
   "fmt"
   "math/big"
   "os"
)

const MAXBITS = 2005

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   var m int
   if _, err := fmt.Fscan(reader, &m); err != nil {
       return
   }
   // basisVec[pos] holds basis vector with highest bit pos
   basisVec := make([]*big.Int, MAXBITS)
   basisMask := make([]*big.Int, MAXBITS)
   for i := 0; i < m; i++ {
       var s string
       fmt.Fscan(reader, &s)
       X := new(big.Int)
       X.SetString(s, 10)
       // Query: check representability
       Xq := new(big.Int).Set(X)
       maskQ := new(big.Int)
       for pos := Xq.BitLen() - 1; pos >= 0; pos-- {
           if pos >= MAXBITS {
               continue
           }
           if Xq.Bit(pos) == 1 {
               if basisVec[pos] == nil {
                   break
               }
               Xq.Xor(Xq, basisVec[pos])
               maskQ.Xor(maskQ, basisMask[pos])
           }
       }
       if Xq.BitLen() == 0 {
           // representable
           // collect indices
           var idxs []int
           for j := 0; j < i; j++ {
               if maskQ.Bit(j) == 1 {
                   idxs = append(idxs, j)
               }
           }
           fmt.Fprint(writer, len(idxs))
           for _, j := range idxs {
               fmt.Fprint(writer, " ", j)
           }
           fmt.Fprintln(writer)
       } else {
           // not representable
           fmt.Fprintln(writer, 0)
       }
       // Insert into basis
       Xi := new(big.Int).Set(X)
       maskI := new(big.Int)
       maskI.SetBit(maskI, i, 1)
       for pos := Xi.BitLen() - 1; pos >= 0; pos-- {
           if pos >= MAXBITS {
               continue
           }
           if Xi.Bit(pos) == 1 {
               if basisVec[pos] == nil {
                   // new independent vector
                   basisVec[pos] = new(big.Int).Set(Xi)
                   basisMask[pos] = new(big.Int).Set(maskI)
                   break
               }
               Xi.Xor(Xi, basisVec[pos])
               maskI.Xor(maskI, basisMask[pos])
           }
       }
   }
}

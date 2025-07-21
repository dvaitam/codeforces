package main
import (
	"fmt"
	"time"
	"math/bits"
)
var n int = 13
var count uint64
func dfs(i int, colmask uint32, diagmask uint32) {
	if i==n {
		count++;
		return
	}
	avail := ^colmask & ((1<<n)-1)
	for avail!=0 {
		j := bits.TrailingZeros32(avail)
		avail &= avail-1
		d := (i+j)%n
		if diagmask&(1<<d)==0 {
			dfs(i+1, colmask|(1<<j), diagmask|(1<<d))
		}
	}
}
func main(){ t0:=time.Now(); dfs(0,0,0); fmt.Println("g=",count, "time",time.Since(t0)); }

package main

import (
	"bytes"
	"fmt"
	"math/big"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded source for the reference solution (was 126D.go).
const solutionSource = `package main

import (
   "bufio"
   "fmt"
   "math/big"
   "os"
)

// Count decompositions of n into distinct Fibonacci numbers F1=1, F2=2, Fi=Fi-1+Fi-2
var fib []int64
var sumFib []int64
var memo map[key]*big.Int

type key struct{ n int64; k int }

// f(n,k) returns number of ways to write n using fib[0..k]
func f(n int64, k int) *big.Int {
   if n == 0 {
       return big.NewInt(1)
   }
   if k < 0 {
       return big.NewInt(0)
   }
   // prune: if sum of all fib[0..k] < n, impossible
   if sumFib[k] < n {
       return big.NewInt(0)
   }
   key := key{n, k}
   if v, ok := memo[key]; ok {
       return v
   }
   var res *big.Int
   if fib[k] > n {
       res = f(n, k-1)
   } else {
       a := f(n, k-1)
       b := f(n-fib[k], k-1)
       res = new(big.Int).Add(a, b)
   }
   memo[key] = res
   return res
}

func main() {
   // precompute fibs up to >1e18
   fib = make([]int64, 0, 100)
   fib = append(fib, 1, 2)
   for {
       i := len(fib)
       next := fib[i-1] + fib[i-2]
       if next < 0 || next > 1e18 {
           break
       }
       fib = append(fib, next)
   }
   // prefix sums
   sumFib = make([]int64, len(fib))
   for i, v := range fib {
       if i == 0 {
           sumFib[i] = v
       } else {
           sumFib[i] = sumFib[i-1] + v
       }
   }
   memo = make(map[key]*big.Int)

   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   var t int
   if _, err := fmt.Fscan(in, &t); err != nil {
       return
   }
   for i := 0; i < t; i++ {
       var n int64
       fmt.Fscan(in, &n)
       // find largest k with fib[k] <= n
       k := len(fib) - 1
       for k >= 0 && fib[k] > n {
           k--
       }
       ans := f(n, k)
       out.WriteString(ans.String())
       out.WriteByte('\n')
   }
}
`

const testcasesRaw = `1 800330560049
4 530889784860 351305571888 738262019349 882990998296
2 464368953796 93061711152
5 962013344778 954556899283 600300332388 976434342708 887874639766
2 8763160604 421448183536
3 33546762542 890658434453 568374026914
5 926997791565 652009656528 594419254959 123322596851 679055963705
5 232757391105 705760880888 444765071106 43569492135 991610559956
1 38308927472
3 967260183494 624538550165 711653100292
1 229400993115
3 112057942787 798580481004 862235809330
2 258871966324 725523749136
2 424549975507 216617710106
2 307185923566 783482175090
2 562473902132 889495211280
5 933566795201 699570982920 368782656162 279161690854 255590652537
4 234415998701 58266151250 994569128436 118803213808
3 706230936473 859429189673 107535001997
3 125974096068 311907051062 262413215623
3 630453663177 103357963160 525667764155
2 360375409749 142700224193
2 43542021631 625838576101
1 525994647839
5 260067874115 955305161197 237286226990 896029756682 365490134338
2 325198088172 490064413151
4 522479900882 864413100207 208215687002 497378243393
2 234516825050 626345068134
4 473952935815 426327683712 142451838969 655948891997
5 122364409019 261125423376 994751223307 748807255917 622165869830
3 335039237599 955398571660 235740227234
4 318792841744 346085528544 395241854239 37860374193
4 807142242917 28998949571 874335935323 603979974054
1 54757190821
1 329094194678
4 705798903309 987955270164 418998985913 460017470673
3 810868959124 876360853427 48700739786
1 718758670117
1 758579249065
1 695256046358
5 110081709569 630648394443 595945493956 232076058535 638198062444
4 308498745368 116396686344 133919987149 171209683083
5 254430025594 439202997687 56711200570 582062783050 934103982960
1 252852151567
2 757605676916 725566931026
4 270867947389 958758188550 168012411435 959641735538
4 280978324066 891671980257 602119880577 605749019159
4 769465602701 813867188408 891518703797 52135922420
3 641609357224 732205217641 654794263318
4 847534869138 956231412558 255692134443 72337543401
3 587812469697 259852308196 333817815839
2 562051008660 298218230791
3 419996293134 817427910961 914144443342
2 106700590381 327720102544
2 266724231937 178386859988
3 460475606631 510214885528 935731813023
2 572409944926 403461773102
4 72426745522 483082845134 105998271136 184618730816
1 578996590180
3 9727885774 431172839096 26623972158
5 863987801624 455526076081 366900004135 471995827546 793728513898
4 292232270361 523661559065 650999204416 996493326890
1 163360940868
5 618416080641 181404766673 948756580737 368186292443 152150166248
5 272672271713 84343964315 160070857661 213117355129 147087303272
1 789822994692
1 302545206540
3 578653711498 416699151727 886767986170
5 850276306968 183189585253 180553891644 947701060487 571270271535
4 100305990738 953390753696 27115867131 747119051462
4 561439947614 284464970153 416557140681 472700267270
4 540988962554 915539757731 692418157317 945365244310
5 105294210930 238471980629 848809020936 935065901928 935163170841
5 442052328231 302029469010 808445717937 66990113506 984745984946
1 776904706399
5 917342277741 205309279543 502972624480 24810937858 388600195227
3 994005589537 39399181654 578588276473
2 179751746401 822186596155
2 662716541138 588197036847
5 56054482378 26757671745 844271713400 738128115809 406969657017
3 891011629253 798664907821 311863495293
3 363030371371 412291861660 145941465729
1 483798911742
5 902431204998 106227096086 171743885677 370073478886 238004951405
5 795647552903 688706740471 280693023897 959958005808 11766880737
5 982706969412 72720694387 725598428959 727655584986 631837405360
4 869556722879 820656975060 790702103465 921079967216
2 809094012613 447469505369
5 38474629420 465334829305 310433448530 90800385063 699967262939
5 598905260363 347545092230 880098534316 208715865815 732550816643
1 244097311614
4 797737577552 556753332497 147782769794 595084151464
1 217838046897
3 225725536263 61268537162 401637239624
3 143383527587 210611405799 329100176951
2 324047457533 42963470092
4 603169479419 760253614205 743744508591 399875420729
4 958767787518 295345719141 441390207964 388696293368
5 795235464016 356874252309 731067790577 959467178392 415579963007
3 554889497668 832000709913 623785971386`

var _ = solutionSource

type testCase []int64

var fib []int64
var sumFib []int64

// f counts representations of n using fib[0..k] (distinct), matching solution logic.
func f(n int64, k int, memo map[[2]int64]*big.Int) *big.Int {
	if n == 0 {
		return big.NewInt(1)
	}
	if k < 0 {
		return big.NewInt(0)
	}
	if sumFib[k] < n {
		return big.NewInt(0)
	}
	key := [2]int64{n, int64(k)}
	if v, ok := memo[key]; ok {
		return v
	}
	var res *big.Int
	if fib[k] > n {
		res = f(n, k-1, memo)
	} else {
		a := f(n, k-1, memo)
		b := f(n-fib[k], k-1, memo)
		res = new(big.Int).Add(a, b)
	}
	memo[key] = res
	return res
}

func parseTests() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcasesRaw), "\n")
	tests := make([]testCase, 0, len(lines))
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 1 {
			return nil, fmt.Errorf("invalid line %d", idx+1)
		}
		t, _ := strconv.Atoi(fields[0])
		if len(fields) != t+1 {
			return nil, fmt.Errorf("invalid test count on line %d", idx+1)
		}
		arr := make([]int64, t)
		for i := 0; i < t; i++ {
			val, _ := strconv.ParseInt(fields[i+1], 10, 64)
			arr[i] = val
		}
		tests = append(tests, arr)
	}
	return tests, nil
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func prepareFib() {
	fib = []int64{1, 2}
	for {
		i := len(fib)
		next := fib[i-1] + fib[i-2]
		if next < 0 || next > 1e18 {
			break
		}
		fib = append(fib, next)
	}
	sumFib = make([]int64, len(fib))
	for i, v := range fib {
		if i == 0 {
			sumFib[i] = v
		} else {
			sumFib[i] = sumFib[i-1] + v
		}
	}
}

func expected(tc testCase) string {
	var sb strings.Builder
	memo := make(map[[2]int64]*big.Int)
	for idx, n := range tc {
		k := len(fib) - 1
		for k >= 0 && fib[k] > n {
			k--
		}
		ans := f(n, k, memo)
		sb.WriteString(ans.String())
		if idx+1 < len(tc) {
			sb.WriteByte('\n')
		}
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	prepareFib()
	tests, err := parseTests()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	for i, tc := range tests {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", len(tc)))
		for _, v := range tc {
			sb.WriteString(strconv.FormatInt(v, 10))
			sb.WriteByte('\n')
		}
		input := sb.String()
		want := expected(tc)
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("Test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != want {
			fmt.Printf("Test %d failed: expected %s got %s\n", i+1, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

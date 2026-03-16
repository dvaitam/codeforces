package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func buildOracle() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	oracle := filepath.Join(dir, "oracleD")
	cmd := exec.Command("go", "build", "-o", oracle, "909D.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return oracle, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	const testcasesRaw = `sreltpus
tap
rhgwprrpm
hueqm
v
fys
jy
i
txmwznmxzsoeldbe
givnyujnqmslrsns
kvaitvwf
rssdwugusij
cpup
lzc
eajnyndbttybmw
kriqhbjacdtrbgnjtie
kk
emmoqmutvrdt
inuxwhjniqjrkazns
amtsuebuuko
vltwixpasbva
iuojstkflfky
tijzmdyasvxe
qhuzihkfvn
dtkk
ozfckxug
oihzdbqgkzsfikzuczt
senjqziolunj
snbnegaptqnrwh
xo
jrkhcsjdzhbbzwqgn
bapxdfqjhvaqrnbtdke
rpzzblhgd
dfhzizeapusmbyihit
qnbpkyabyebdbcpbw
qqp
fkclmumsjli
knderaw
zcsfblotuzrmu
tn
lu
ywknwnoahgriwscz
hneaklrzidowdx
zmvdxksrdzswapehy
bqcsdvmfakdad
wjsjzcbysqqwhdrx
rbrk
fchfuhotwymiltmlrnc
qhnxfnwsysvqv
eumefdpxpwqosxfe
ygesqkhwr
vwntssigja
pzmgfslhk
yenwpwtgosurapxc
zxbohhuwyvcgi
gyieftwv
if
kf
cxzcdcijblosxv
aakknmpcgus
merkdicvndoqidqw
vylyojvvvuzi
ykvs
qdvpqlbwjvxsxfuuxu
fluod
reku
nrjufopjzfwcdwfyrrsx
ldiimbebpqihw
lkmorzyclpdeisdvd
yxdfwgsnvxmxestemzg
qfsfgilzjazonmkrsj
qvwjvpatgxuadyyv
pfquoggz
gbquodsjveeozctba
thqcprakkklw
ctybw
xkz
cgnwyhp
dzbncgwfmpp
rng
jaooywmofobxillo
ltmhazgizleorgfga
smqfua
dtfop
bamokn
ww
hm
mp
h
dmpgfktd
dtbzxjizozjp
riakulkc
vn
sta
avca
qbpbgu
kgypkzplvbmjytumc
fndqmrkrvy
fxxmrlflznoho
lifqxytwwmpbefwy`

	scanner := bufio.NewScanner(strings.NewReader(testcasesRaw))
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		input := "1\n" + line + "\n"
		cmdO := exec.Command(oracle)
		cmdO.Stdin = strings.NewReader(input)
		var outO bytes.Buffer
		cmdO.Stdout = &outO
		if err := cmdO.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "oracle run error: %v\n", err)
			os.Exit(1)
		}
		expected := strings.TrimSpace(outO.String())

		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\nstderr: %s\n", idx, err, stderr.String())
			os.Exit(1)
		}
		got := strings.TrimSpace(out.String())
		if got != expected {
			fmt.Printf("test %d failed. Expected %s got %s\n", idx, expected, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}

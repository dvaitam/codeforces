#include <bits/stdc++.h>
#define rep(i, l, r) for (int i = l; i <= r; i++)
#define per(i, r, l) for (int i = r; i >= l; i--)
#define srep(i, l, r) for (int i = l; i < r; i++)
#define sper(i, r, l) for (int i = r; i > l; i--)
#define erep(i, x) for (int i = h[x]; i; i = e[i].next)
#define erep2(i, x) for (int& i = cur[x]; i; i = e[i].next)
#define pii pair<int, int>
#define pll pair<ll, ll>
#define pdd pair<double, double>
#define fi first
#define se second
#define ui unsigned int
#define ld long double
#define pb push_back
#define pc putchar
#define lowbit(x) (x & -x)
#define maxbuf 2000020
#define gc() ((p1 == p2 && (p2 = (p1 = buffer) + fread(buffer, 1, maxbuf, stdin), p1 == p2)) ? EOF : *p1++)
using namespace std;

namespace Fast_Read{
    char buffer[maxbuf], *p1, *p2;
    template<class T> void read_signed(T& x){
        char ch = gc(); x = 0; bool f = 1;
        while (!isdigit(ch) && ch != '-') ch = gc();
        if (ch == '-') f = 0, ch = gc();
        while ('0' <= ch && ch <= '9') x = (x << 1) + (x << 3) + ch - '0', ch = gc();
        x = (f) ? x : -x;
    }
    template<class T, class... Args> void read_signed(T& x, Args&... args){
        read_signed(x), read_signed(args...);
    }
    template<class T> void read_unsigned(T& x){
        char ch = gc(); x = 0;
        while (!isdigit(ch)) ch = gc(); 
        while (isdigit(ch)) x = (x << 1) + (x << 3) + ch - '0', ch = gc(); 
    }
    template<class T, class... Args> void read_unsigned(T& x, Args&... args){
        read_unsigned(x), read_unsigned(args...);
    }
    #define isletter(ch) ('a' <= ch && ch <= 'z')
    int read_string(char* s){
        char ch = gc(); int l = 0;
        while (!isletter(ch)) ch = gc();
        while (isletter(ch)) s[l++] = ch, ch = gc();
        s[l] = '\0'; return l;
    }
}using namespace Fast_Read; 

int _num[20];
template <class T> void write(T x, char sep = '\n'){	
	if (!x) {putchar('0'), putchar(sep); return;}
	if (x < 0) putchar('-'), x = -x;
	int c = 0;
	while (x) _num[++c] = x % 10, x /= 10;
	while (c) putchar('0' + _num[c--]); 
	putchar(sep);
}

#define read read_signed
#define reads read_string 
#define writes puts

#define maxn 400020
#define maxm
#define maxs
#define maxb
#define inf 
#define eps
#define M 
#define ll long long int 


int n, a[maxn], p[maxn], cp = 0;
int main(){
    int T; read(T);
    while (T--){
        read(n);
        int num1 = 0;
        rep(i, 1, n) read(a[i]), num1 += abs(a[i]);
        if (num1 & 1) {
            puts("-1");
            continue;
        }
        int i = 1;
        int rev = 0;
        cp = 0;
        p[++cp] = 1;
        while(1){
            while (a[i] == 0 && i <= n) i++;
            if (i > n) break;
            int l = i; i += 1;
            while (a[i] == 0 && i <= n) i++;
            assert(i <= n);
            int r = i; i += 1;
            if ((r - l) & 1) {
                if (a[l] == a[r]) continue;
                int sgnl = (l & 1) ^ rev;
                if (sgnl){
                    if (p[cp] != r) p[++cp] = r;
                    p[++cp] = r + 1;
                }
                else{   
                    if (p[cp] != l) p[++cp] = l;
                    p[++cp] = l + 1;
                }
            }
            else {
                if (a[l] != a[r]) continue;
                int sgnl = (l & 1) ^ rev;
                if (sgnl){
                    if (p[cp] != r) p[++cp] = r - 1;
                    p[++cp] = r + 1;
                    rev ^= 1;
                }
                else{   
                    if (p[cp] != l) p[++cp] = l;
                    p[++cp] = l + 1;
                }
            }
        }
        if (p[cp] != n + 1) p[++cp] = n + 1;
        printf("%d\n", cp - 1);
        srep(i, 1, cp) {
            printf("%d %d\n", p[i], p[i + 1] - 1);
        }
    }
	return 0;
}
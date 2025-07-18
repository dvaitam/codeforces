#include <bits/stdc++.h>
using namespace std;

typedef double lf;
typedef long long ll;
typedef vector<ll> vll;
typedef long double llf;
typedef pair<ll, ll> pll;
typedef vector<int> vint;
typedef unsigned int uint;
typedef pair<int, int> pii;
typedef unsigned long long ull;

bool Dbg;
const int mod = 998244353, MAXN = 5e4 + 7, inft = 1e9 + 7; const ll infl = 1ll << 59; const llf eps = 1e-7;

#define xx first
#define yy second
#define pb push_back
#define mp make_pair
#define mid ((l+r)>>1)
//l = mid+1, ans = mid
#define Debug(x) (Dbg ? void(cerr << #x << " = " << x << ' ') : void())
#define debug(...) (Dbg ? (fprintf(stderr, __VA_ARGS__)) : 0)
#define lop(i,a,b) for(int i = (a), i##end = (b); i <= i##end; ++i)
//toposort
#define dlop(i,a,b) for(int i = (a), i##end = (b); i >= i##end; --i)
#define ergo(a) for(__typeof((a).end())it = (a).begin(), it##end = (a).end(); it != it##end; ++it)

template<typename T> inline T sqr(T x) {return x * x;}
template<typename T> inline T abs(T x) {return x >= 0 ? x : -x;}
template<typename A, typename B> inline A max(A a, B b) {return a > b ? a : b;}
template<typename A, typename B> inline A min(A a, B b) {return a < b ? a : b;}
template<typename A, typename B> inline void chmax(A &x, B y) {if (x < y) x = y;}
template<typename A, typename B> inline void chmin(A &x, B y) {if (x > y) x = y;}
template<typename A, typename B> inline void amod(A &x, B y) {x += y; if (x >= mod) x -= mod;}
template<typename A, typename B> inline A gcd(A a, B b) {if (a < b) swap(a, b); if (!b) return a; while (A t = a % b) a = b, b = t; return b;}
template<typename A, typename B> inline A lcm(A a, B b) {return a / gcd(a, b) * b;}
template<typename A, typename B> inline A Pow(A a, B b, int mod) {A ret; for (ret = 1; b; b >>= 1) {if (b & 1) ret = ret * 1ll * a % mod; a = a * 1ll * a % mod;} return ret % mod;}

struct IO {
#define IS(x) (x == 10 || x == 13 || x == ' ')
#define OP operator
#define RT return *this
#define RX x=0;int t=P();while((t<'0'||t>'9')&&t!='-')t=P();f=1;\
if(t=='-')t=P(),f=-1;x=t-'0';for(t=P();t>='0'&&t<='9';t=P())x=x*10+t-'0'
#define RL if(t=='.'){lf u=0.1;for(t=P();t>='0'&&t<='9';t=P(),u*=0.1)x+=u*(t-'0');}if(f==-1)x=-x;
#define RU x=0;int t=P();while(t<'0'||t>'9')t=P();x=t-'0';for(t=P();t>='0'&&t<='9';t=P())x=x*10+t-'0'
#define TR *this,x;return x
#define WI if(x){if(x<0)P('-'),x=-x;c=0;while(x)s[c++]=x%10+'0',x/=10;while(c--)P(s[c]);}else P('0')
#define WL if(y){lf t=0.5;for(int i=y;i--;)t*=0.1;if(x>=0)x+=t;else x-=t,P('-');*this,(ll)(abs(x));P('.');if(x<0)\
x=-x;while(y--){x*=10;x-=floor(x*0.1)*10;P(((int)x)%10+'0');}}else if(x>=0)*this,(ll)(x+0.5);else *this,(ll)(x-0.5);
#define WU if(x){c=0;while(x)s[c++]=x%10+'0',x/=10;while(c--)P(s[c]);}else P('0')
  struct Cg {inline int operator()() {return getchar();}}; struct Cp {inline void operator()(int x) {putchar(x);}}; template<typename T>struct Fr {int f; T P; inline Fr&OP, (int&x) {RX; x *= f; RT;} inline OP int() {int x; TR;} inline Fr&OP, (ll &x) {RX; x *= f; RT;} inline OP ll() {ll x; TR;} inline Fr&OP, (char&x) {for (x = P(); IS(x); x = P()); RT;} inline OP char() {char x; TR;} inline Fr&OP, (char*x) {char t = P(); for (; IS(t); t = P()); if (~t) {for (; !IS(t) && ~t; t = P()) * x++ = t;}*x++ = 0; RT;} inline Fr&OP, (lf&x) {RX; RL; RT;} inline OP lf() {lf x; TR;} inline Fr&OP, (llf&x) {RX; RL; RT;} inline OP llf() {llf x; TR;} inline Fr&OP, (uint&x) {RU; RT;} inline OP uint() {uint x; TR;} inline Fr&OP, (ull&x) {RU; RT;} inline OP ull() {ull x; TR;}}; Fr<Cg>in; template<typename T>struct Fw {int c, s[24]; T P; inline Fw&OP, (int x) {WI; RT;} inline Fw&OP()(int x) {WI; RT;} inline Fw&OP, (uint x) {WU; RT;} inline Fw&OP()(uint x) {WU; RT;} inline Fw&OP, (ll x) {WI; RT;} inline Fw&OP()(ll x) {WI; RT;} inline Fw&OP, (ull x) {WU; RT;} inline Fw&OP()(ull x) {WU; RT;} inline Fw&OP, (char x) {P(x); RT;} inline Fw&OP()(char x) {P(x); RT;} inline Fw&OP, (const char*x) {while (*x)P(*x++); RT;} inline Fw&OP()(const char*x) {while (*x)P(*x++); RT;} inline Fw&OP()(lf x, int y) {WL; RT;} inline Fw&OP()(llf x, int y) {WL; RT;}}; Fw<Cp>out;
} io;
#define in io.in
#define out io.out


int q, l, r, Number, a[MAXN << 2], visit[MAXN << 2];
char s[35];

int main() {
  in, q;
  l = r = 3e5;
  lop(i, 1, q) {
    in, s, Number;
    if (i == 1) {
      if (s[0] == 'L') 
        visit[Number] = l;
       else if (s[0] == 'R') 
        visit[Number] = r;
      
    } else 
      if (s[0] == 'L') visit[Number] = --l; else if (s[0] == 'R') visit[Number] = ++r; else if (s[0] == '?') out, min(visit[Number] - l, r - visit[Number]), '\n';
    
  }
  return 0;
}
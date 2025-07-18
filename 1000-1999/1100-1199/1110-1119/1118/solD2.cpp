#include <bits/stdc++.h>
using namespace std;
bool Dbg;

typedef double lf;
typedef long long ll;
typedef long double llf;
typedef vector<int> vint;
typedef pair<ll, ll> pll;
typedef unsigned int uint;
typedef pair<int, int> pii;
typedef unsigned long long ull;

#define fi first
#define se second
#define pb push_back
#define mp make_pair
#define mid ((l+r)>>1)
#define all(x) x.begin(), x.end()
#define dlop1(i,b) for (int i = (b); i >= 1; --i)
#define dlop0(i,b) for (int i = (b)-1; i >= 0; --i)
#define lop0(i,b) for (int i = 0, i##end = (b); i < i##end; ++i)
#define lop1(i,b) for (int i = 1, i##end = (b); i <= i##end; ++i)
#define lop(i,a,b) for (int i = (a), i##end = (b); i <= i##end; ++i)
#define dlop(i,a,b) for (int i = (a), i##end = (b); i >= i##end; --i)
#define debug(...) (Dbg ? void(fprintf(stderr, __VA_ARGS__)) : void())
#if __cplusplus >= 201103L
mt19937 Rand(time(0) ^ (ull)(new char));
#else
inline uint Rand() {static uint x = time(0) ^ (ull)(new char); x ^= x << 13; x ^= x >> 17; x ^= x << 5; return x;}
#endif
#define Debug(x) (Dbg ? void(cerr << #x << " = " << x << '\n') : void())
#define trav(it, a) for (__typeof((a).end())it = (a).begin(); it != (a).end(); ++it)
#define dtrav(it, a) for (__typeof((a).rend())it = (a).rbegin(); it != (a).rend(); ++it)
#define IS(x) (x == 10 || x == 13 || x == ' ')
#define OP operator
#define RT return *this
#define TR *this,x;return x;
#define RX x=0;t=P();while((t<'0'||t>'9')&&t!='-')t=P();f=1;\
  if(t=='-')t=P(),f=-1;x=t-'0';for(t=P();t>='0'&&t<='9';t=P())x=x*10+t-'0'
#define RU x=0;t=P();while(t<'0'||t>'9')t=P();x=t-'0';for(t=P();t>='0'&&t<='9';t=P())x=x*10+t-'0'
#define WI if(x){if(x<0)P('-'),x=-x;c=0;while(x)s[c++]=x%10+'0',x/=10;while(c--)P(s[c]);}else P('0')
#define WU if(x){c=0;while(x)s[c++]=x%10+'0',x/=10;while(c--)P(s[c]);}else P('0')
struct Cg {int operator()() {return getchar(); } }; struct Cp {void operator()(int x) {putchar(x); } };
struct Fr {
  int f, t; Cg P;
#ifdef __SIZEOF_INT128__
  Fr &OP, (__int128 &x) {RX; x *= f; RT; }
  OP __int128() {__int128 x; TR; }
#endif
  Fr &OP, (int &x) {RX; x *= f; RT; } OP int() {int x; TR; } Fr &OP, (ll &x) {RX; x *= f; RT; } OP ll() {ll x; TR; } Fr &OP, (char &x) {for (x = P(); IS(x); x = P()); RT; } OP char() {char x; TR; } Fr &OP, (string &x) {cin >> x; RT; } OP string() {string x; TR; } Fr &OP, (char *x) {char t = P(); for (; IS(t); t = P()); if (~t) {for (; !IS(t) && ~t; t = P()) * x++ = t; } *x++ = 0; RT; } Fr &OP, (lf &x) {scanf("%lf", &x); RT; } OP lf() {lf x; TR; } Fr &OP, (llf &x) {lf y; scanf("%lf", &y); x = y; RT; } OP llf() {llf x; TR; } Fr &OP, (uint &x) {RU; RT; } Fr &OP, (ull &x) {RU; RT; } OP uint() {uint x; TR; }
} in;
struct Fw {
  int c, s[50]; Cp P;
#ifdef __SIZEOF_INT128__
  Fw &OP, (__int128 x) {WI; RT; }
#endif
  Fw &OP, (int x) {WI; RT; } Fw &OP, (uint x) {WU; RT; } Fw &OP, (ll x) {WI; RT; } Fw &OP, (ull x) {WU; RT; } Fw &OP, (char x) {P(x); RT; }  Fw &OP, (const string &x) {cout << x; RT; } Fw &OP, (const char *x) {while (*x)P(*x++); RT; }
  Fw &OP, (const lf &x) {printf("%.5f", x); RT; } Fw &OP, (const llf &x) {printf("%.5f", lf(x)); RT; }
} out;
const int mod = 998244353, MAXN = 1e5 + 7, inft = 1e9 + 7; const ll infl = llf(1e18) + 1;
const lf eps = 1e-7;
template<class T> inline T sqr(T x) {return x * x; }
template<class A, class B> inline A _gcd(A a, B b) {A t; if (a < b) swap(a, b); if (!b) return a; while ((t = a % b)) a = b, b = t; return b; }
template<class A, class B> inline ll _lcm(A a, B b) {return a / gcd(a, b) * 1ll * b; } template<class T> inline T abs(T x) {return x >= 0 ? x : -x; }
template<class A, class B> inline ll mul(A a, B b, ll mod) {if (b < 0) b = -b, a = -a; ll ret; for (ret = 0; b; b >>= 1) {if (b & 1) ret = (ret + a) % mod; a = (a + a) % mod;} return ret % mod; } template<class A, class B> inline A Pow1(A a, B b, int mod) {A ret; for (ret = 1; b; b >>= 1) {if (b & 1) ret = ret * 1ll * a % mod; a = a * 1ll * a % mod; } return ret % mod; } template<class A, class B> inline ll Pow(A a, B b, ll mod) {assert(b >= 0); a %= mod; if (mod <= 2e9) return Pow1(a, b, mod); ll ret; for (ret = 1; b; b >>= 1) {if (b & 1) ret = mul(ret, a, mod); a = mul(a, a, mod); } return ret % mod; }
template<class A, class B> inline A max(A a, B b) {return a > b ? a : b; } template<class A, class B> inline A min(A a, B b) {return a < b ? a : b; }
template<class A, class B> inline bool chmax(A &x, B y) {return x < y ? x = y, 1 : 0;}
template<class A, class B> inline bool chmin(A &x, B y) {return x > y ? x = y, 1 : 0;}

void init() {
#ifdef QvvQ
  // freopen("data.in", "r", stdin);
  // freopen("data.out", "w", stdout);
  Dbg = 1;
#endif
}
ll n, a[MAXN << 1], m, sum;
inline bool check(int X) {
  ll tmp = 0; int cnt = 0, cur = 0;
  dlop1(i, n) {
    tmp += max(0, a[i] - cur);
    if (++cnt == X) ++cur, cnt = 0;
  }
  return tmp >= m;
}
void solve() {
  in, n, m;
  lop1(i, n) sum += (a[i] = in);
  sort(a + 1, a + 1 + n);
  int l = 1, r = n, ans = -1;
  while (l <= r) {
    if (check(mid)) ans = mid, r = mid - 1;
    else l = mid + 1;
  }
  out, ans;
}

void quit() {
#ifdef QvvQ
  fprintf(stderr, "\ntime:%.5fms", clock() * 1000.0 / CLOCKS_PER_SEC);
#endif
}

int main() {
  init();
  int T = 1;
  while (T--) solve();
  quit();
  return 0;
}
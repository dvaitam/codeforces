#include <bits/stdc++.h>
using namespace std;
bool Dbg;
typedef double lf; typedef long long ll; typedef long double llf; typedef vector<int> vint; typedef unsigned int uint; typedef pair<int, int> pii; typedef unsigned long long ull;

#define xx first
#define yy second
#define pb push_back
#define mp make_pair
#define mid ((l+r)>>1)
#define all(x) x.begin(), x.end()
#define debug(...) (Dbg ? void(fprintf(stderr, __VA_ARGS__)) : void())
#if __cplusplus <= 201103L
#define lop0(i,b) for (register int i = 0, i##end = (b); i < i##end; ++i)
#define lop1(i,b) for (register int i = 1, i##end = (b); i <= i##end; ++i)
#define dlop(i,a,b) for (register int i = (a), i##end = (b); i >= i##end; --i)
#define lop(i,a,b) for (register int i = (a), i##end = (b); i <= i##end; ++i)
#define dlop0(i,b) for (register int i = (b)-1; i >= 0; --i)
#define dlop1(i,b) for (register int i = (b); i >= 1; --i)
#else
#define lop0(i,b) for (int i = 0, i##end = (b); i < i##end; ++i)
#define lop1(i,b) for (int i = 1, i##end = (b); i <= i##end; ++i)
#define dlop(i,a,b) for (int i = (a), i##end = (b); i >= i##end; --i)
#define lop(i,a,b) for (int i = (a), i##end = (b); i <= i##end; ++i)
#define dlop0(i,b) for (int i = (b)-1; i >= 0; --i)
#define dlop1(i,b) for (int i = (b); i >= 1; --i)
#endif
#if __cplusplus >= 201103L
mt19937 Rand(time(0) ^ (ull)(new char));
#define mt make_tuple
#else
uint Rand() {static uint x = time(0) ^ (ull)(new char); x ^= x << 13; x ^= x >> 17; x ^= x << 5; return x;}
#endif
#define Debug(x) (Dbg ? void(cerr << #x << " = " << x << '\n') : void())
#define ergo(it, a) for (__typeof((a).end())it = (a).begin(); it != (a).end(); ++it)
#define dergo(it, a) for (__typeof((a).rend())it = (a).rbegin(); it != (a).rend(); ++it)
#define ergo1(it, a) for (__typeof((a).end())it = (a).begin(), it##1; it != (a).end(); it = it1)
#define dergo1(it, a) for (__typeof((a).rend())it = (a).rbegin(), it##1; it != (a).rend(); it = it1)
#define IS(x) (x == 10 || x == 13 || x == ' ')
#define OP operator
#define RT return *this
#define RX x=0;t=P();while((t<'0'||t>'9')&&t!='-')t=P();f=1;\
  if(t=='-')t=P(),f=-1;x=t-'0';for(t=P();t>='0'&&t<='9';t=P())x=x*10+t-'0'
#define RU x=0;t=P();while(t<'0'||t>'9')t=P();x=t-'0';for(t=P();t>='0'&&t<='9';t=P())x=x*10+t-'0'
#define TR *this,x;return x
#define WI if(x){if(x<0)P('-'),x=-x;c=0;while(x)s[c++]=x%10+'0',x/=10;while(c--)P(s[c]);}else P('0')
#define WU if(x){c=0;while(x)s[c++]=x%10+'0',x/=10;while(c--)P(s[c]);}else P('0')
struct Cg {inline int operator()() {return getchar(); } }; struct Cp {inline void operator()(int x) {putchar(x); } };
struct Fr {
	int f, t; Cg P;
#ifdef __SIZEOF_INT128__
	inline Fr&OP, (__int128 &x) {RX; x *= f; RT; }
#endif
	inline Fr&OP, (int &x) {RX; x *= f; RT; } inline Fr&OP, (ll &x) {RX; x *= f; RT; } inline Fr&OP, (char &x) {for (x = P(); IS(x); x = P()); RT; } inline Fr&OP, (string &x) {cin >> x; RT; } inline Fr&OP, (char *x) {char t = P(); for (; IS(t); t = P()); if (~t) {for (; !IS(t) && ~t; t = P()) * x++ = t; }*x++ = 0; RT; } inline Fr&OP, (lf &x) {scanf("%lf", &x); RT; } inline Fr&OP, (llf &x) {lf y; scanf("%lf", &y); x = y; RT; } inline Fr&OP, (uint &x) {RU; RT; } inline Fr&OP, (ull &x) {RU; RT; }
} in;
struct Fw {
	int c, s[50]; Cp P;
#ifdef __SIZEOF_INT128__
	inline Fw&OP, (__int128 x) {WI; RT; }
#endif
	inline Fw&OP, (int x) {WI; RT; } inline Fw&OP, (uint x) {WU; RT; } inline Fw&OP, (ll x) {WI; RT; } inline Fw&OP, (ull x) {WU; RT; } inline Fw&OP, (char x) {P(x); RT; } inline Fw&OP, (lf x) {printf("%.3f", x); RT; } inline Fw&OP, (const llf &x) {printf("%.5lf", lf(x)); RT; } inline Fw&OP, (const string &x) {cout << x; RT; } inline Fw&OP, (const char *x) {while (*x)P(*x++); RT; }
} out;
const int mod = 998244353, MAXN = 1e5 + 7, inft = 1e9 + 7; const ll infl = llf(1e18) + 1;
const lf eps = 1e-7;
template<typename T> inline T sqr(T x) {return x * x; }
template<typename A, typename B> inline A gcd(A a, B b) {A t; if (a < b) swap(a, b); if (!b) return a; while (t = a % b) a = b, b = t; return b; }
template<typename A, typename B> inline A lcm(A a, B b) {return a / gcd(a, b) * b; } template<typename T> inline T abs(T x) {return x >= 0 ? x : -x; }
template<typename A, typename B> inline ll mul(A a, B b, ll mod) {if (b < 0) b = -b, a = -a; ll ret; for (ret = 0; b; b >>= 1) {if (b & 1) ret = (ret + a) % mod; a = (a + a) % mod;} return ret % mod; }
template<typename A, typename B> inline A Pow1(A a, B b, int mod) {A ret; for (ret = 1; b; b >>= 1) {if (b & 1) ret = ret * 1ll * a % mod; a = a * 1ll * a % mod; } return ret % mod; }
template<typename A, typename B> inline ll Pow(A a, B b, ll mod) {assert(b >= 0); a %= mod; if (mod <= 2e9) return Pow1(a, b, mod); ll ret; for (ret = 1; b; b >>= 1) {if (b & 1) ret = mul(ret, a, mod); a = mul(a, a, mod); } return ret % mod; }
template<typename A, typename B> inline A max(A a, B b) {return a > b ? a : b; } template<typename A, typename B> inline A min(A a, B b) {return a < b ? a : b; } template<typename A, typename B> inline void chmax(A &x, B y) {if (x < y) x = y; } template<typename A, typename B> inline void chmin(A &x, B y) {if (x > y) x = y; } template<typename A, typename B> inline void amod(A &x, B y, int mod) {x += y; while (x < 0) x += mod; while (x > mod) x -= mod; }


pii a[MAXN << 1];
// vint b[MAXN << 1];
bool vis[MAXN << 1];
int main() {
#ifdef LOCAL_DEBUG
	// freopen("data.in", "r", stdin), freopen("data.out", "w", stdout);
	Dbg = 1;
#endif
	int n;
	in, n;
	lop1(i, n) in, a[i].xx, a[i].yy;
	vint ans;
	ans.pb(1);
	vis[1] = 1;
	for (int i = 0; i < n; i += 2) {
		int u = a[ans[i]].xx, v = a[ans[i]].yy;
		if (a[u].xx == v || a[u].yy == v) {
			if (!vis[u]) ans.pb(u);
			if (!vis[v]) ans.pb(v);
		}
		else {
			if (!vis[v]) ans.pb(v);
			if (!vis[u]) ans.pb(u);
		}
		vis[u] = vis[v] = 1;
		// cerr << ans.size() << endl;
	}
	ergo(it, ans) out, *it, ' ';

#ifdef LOCAL_DEBUG
	fprintf(stderr, "\ntime:%.5lfms", clock() * 1.0 / CLOCKS_PER_SEC * 1000);
#endif
	return 0;
}
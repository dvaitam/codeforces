#include <iostream>
#include <cstdio>
#include <cstdlib>
#include <cstring>
#include <vector>
#include <map>
#include <complex>
#include <queue>
#include <algorithm>
#include <string>
#include <stack>
#include <bitset>
#include <cmath>
#include <set>

int N = 1e6, SZ = 320, INF = 1 << 29;
long long LINF = (1LL << 61), mod = 1e9 + 7;
const long double eps = 1e-9, PI = acos(-1.0);

#define lowbit(x) (x & (-(x)))
#define MAX(a, b) ((a) < (b) ? (b) : (a))
#define MIN(a, b) ((a) < (b) ? (a) : (b))
#define rp(a, b, c) for (int a = b; a <= c; ++a)
#define RP(a, b, c) for (int a = b; a < c; ++a)
#define lp(a, b, c) for (int a = b; a >= c; --a)
#define LP(a, b, c) for (int a = b; a > c; --a)
#define rps(i, s) for (int i = 0; s[i]; i++)
#define fson(u) for (int i = g[u]; ~i; i = edg[i].nxt)
#define adde(u, v) edg[++ecnt] = Edge(u, v, 0, s[u]), s[u] = ecnt
#define addew(u, v, w) edg[++ecnt] = Edge(u, v, w, s[u]), s[u] = ecnt
#define MID (l + r >> 1)
#define mst(a, v) memset(a, v, sizeof(a))
#define bg(x)            \
	Edge edg[maxn << x]; \
	int g[maxn], ecnt
#define ex(v)  \
	cout << v; \
	return 0
#define debug(x) cout << "debug: " << x << endl;
#define sqr(x) ((x) * (x))

using namespace std;
typedef long long ll;
typedef unsigned long long ull;
typedef double db;
typedef long double ld;
typedef pair<int, int> pii;
typedef pair<ll, ll> pll;
typedef complex<double> cpx;
typedef vector<int> vi;
typedef vector<ll> vll;
typedef map<int, int> mii;
typedef map<ll, ll> mll;

char READ_DATA;
int SIGNAL_INPUT;
template <typename Type>
inline Type ru(Type & v)
{
	SIGNAL_INPUT = 1;
	while ((READ_DATA = getchar()) < '0' || READ_DATA > '9')
		if (READ_DATA == '-')
			SIGNAL_INPUT = -1;
		else if (READ_DATA == EOF)
			return EOF;
	v = READ_DATA - '0';
	while ((READ_DATA = getchar()) >= '0' && READ_DATA <= '9')
		v = v * 10 + READ_DATA - '0';
	v *= SIGNAL_INPUT;
	return v;
}
inline ll modru(ll & v)
{
	ll p = 0;
	SIGNAL_INPUT = 1;
	while ((READ_DATA = getchar()) < '0' || READ_DATA > '9')
		if (READ_DATA == '-')
			SIGNAL_INPUT = -1;
		else if (READ_DATA == EOF)
			return EOF;
	p = v = READ_DATA - '0';
	while ((READ_DATA = getchar()) >= '0' && READ_DATA <= '9')
	{
		v = (v * 10 + READ_DATA - '0') % mod;
		p = (p * 10 + READ_DATA - '0') % (mod - 1);
	}
	v *= SIGNAL_INPUT;
	return p;
}
template <typename A, typename B>
inline int ru(A & x, B & y)
{
	if (ru(x) == EOF)
		return EOF;
	ru(y);
	return 2;
}
template <typename A, typename B, typename C>
inline int ru(A & x, B & y, C & z)
{
	if (ru(x) == EOF)
		return EOF;
	ru(y);
	ru(z);
	return 3;
}
template <typename A, typename B, typename C, typename D>
inline int ru(A & x, B & y, C & z, D & w)
{
	if (ru(x) == EOF)
		return EOF;
	ru(y);
	ru(z);
	ru(w);
	return 4;
}
inline ll gcd(ll a, ll b)
{
	while (b)
	{
		a %= b;
		swap(a, b);
	}
	return a;
}

inline ll fastmul(ll a, ll b, ll mod = 1e9 + 7)
{
	return (a * b - (ll)((long double)a * b / mod) * mod + mod) % mod;
}

inline ll dirmul(ll a, ll b, ll mod = 1e9 + 7)
{
	return a * b% mod;
}

inline ll ss(ll a, ll b, ll mod = 1e9 + 7, ll(*mul)(ll, ll, ll) = dirmul)
{
	if (b < 0)
	{
		b = -b;
		a = ss(a, mod - 2, mod);
	}
	ll ans = 1;
	while (b)
	{
		if (b & 1)
			ans = mul(ans, a, mod);
		a = mul(a, a, mod);
		b >>= 1;
	}
	return ans;
}

inline int isprime(ll n)
{
	if (n == 1)
		return 0;

	for (ll d = 2; d * d <= n; ++d)
	{
		if (n % d == 0)
			return 0;
	}

	return 1;
}

template <typename Type>
void brc(Type * a, int n)
{
	int k;
	for (int i = 1, j = n / 2; i < n - 1; i++)
	{
		if (i < j)
			swap(a[i], a[j]);

		k = n >> 1;
		while (j >= k)
		{
			j ^= k;
			k >>= 1;
		}
		if (j < k)
			j ^= k;
	}
}
void fft(cpx * a, int n, int inv = 1)
{
	cpx u, t;
	brc(a, n);
	for (int h = 2; h <= n; h <<= 1)
	{
		cpx wn(cos(inv * 2.0 * PI / h), sin(inv * 2.0 * PI / h));
		for (int j = 0; j < n; j += h)
		{
			cpx w(1, 0);
			for (int k = j; k < j + (h >> 1); k++)
			{
				u = a[k];
				t = w * a[k + (h >> 1)];
				a[k] = u + t;
				a[k + (h >> 1)] = u - t;
				w *= wn;
			}
		}
	}
	if (inv == -1)
		RP(i, 0, n)
		a[i] /= n;
}
void ntt(ll * a, int n, int inv = 1)
{
	ll u, t;
	brc(a, n);
	for (int h = 2; h <= n; h <<= 1)
	{
		ll wn = ss(3, (mod - 1) / h);
		if (inv == -1)
			wn = ss(wn, mod - 2);
		for (int j = 0; j < n; j += h)
		{
			ll w = 1;
			for (int k = j; k < j + (h >> 1); k++)
			{
				u = a[k];
				t = w * a[k + (h >> 1)] % mod;
				a[k] = (u + t) % mod;
				a[k + (h >> 1)] = (u - t + mod) % mod;
				(w *= wn) %= mod;
			}
		}
	}
	if (inv == -1)
	{
		ll tmp = ss(n, mod - 2);
		RP(i, 0, n)
			(a[i] *= tmp) %= mod;
	}
}
struct Edge
{
	int u, v, nxt;
	int w;
	Edge(int _u = 0, int _v = 0, int _w = 0, int _nxt = 0)
	{
		u = _u;
		v = _v;
		w = _w;
		nxt = _nxt;
	}

	int operator<(const Edge &b) const
	{
		return w < b.w;
	}
};


const int maxn = 2e5 + 5;
/*------------------------------------------------------------------------yah01------------------------------------------------------------------------*/

int n, a[maxn], f[maxn], g[maxn];

int ans;
string opt;
int main()
{
	ru(n);
	rp(i, 1, n) 
		ru(a[i]);
	f[n] = g[1] = 1;

	rp(i, 2, n)
	{
		if (a[i - 1] > a[i]) g[i] = g[i - 1] + 1;
		else g[i] = 1;
	}

	lp(i, n - 1, 1)
	{
		if (a[i + 1] > a[i]) f[i] = f[i + 1] + 1;
		else f[i] = 1;
	}


	int l = 1, r = n, now = 0;
	while (a[l] > now || a[r] > now)
	{
		int j;
		if (a[l] > now && a[r] > now)
		{
			if (a[l] == a[r])
			{
				j = (f[l] > g[r] ? l : r);
			}
			else 
				j = (a[l] < a[r] ? l : r);
		}
		else if (a[l] > now)
			j = l;
		else if (a[r] > now)
			j = r;

		now = a[j];
		++ans;
		opt.push_back((j == l ? 'L' : 'R'));
		if (j == l) ++l;
		else if (j == r) --r;
	}

	cout << ans << endl << opt;
	
}
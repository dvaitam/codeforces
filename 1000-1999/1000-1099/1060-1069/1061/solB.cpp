#include <cstdio>
#include <cstring>
#include <algorithm>
#include <cmath>

#ifdef WIN32
	#define LL "%I64d"
#else
	#define LL "%lld"
#endif

#ifdef CT
	#define debug(...) printf(__VA_ARGS__)
	#define setfile() 
#else
	#define debug(...)
	#define filename ""
	#define setfile() freopen(filename".in", "r", stdin); freopen(filename".out", "w", stdout)
#endif

#define R register
#define getc() (S == T && (T = (S = B) + fread(B, 1, 1 << 15, stdin), S == T) ? EOF : *S++)
#define dmax(_a, _b) ((_a) > (_b) ? (_a) : (_b))
#define dmin(_a, _b) ((_a) < (_b) ? (_a) : (_b))
#define cmax(_a, _b) (_a < (_b) ? _a = (_b) : 0)
#define cmin(_a, _b) (_a > (_b) ? _a = (_b) : 0)
#define cabs(_x) ((_x) < 0 ? (- (_x)) : (_x))
char B[1 << 15], *S = B, *T = B;
inline int F()
{
	R char ch; R int cnt = 0; R bool minus = 0;
	while (ch = getc(), (ch < '0' || ch > '9') && ch != '-') ;
	ch == '-' ? minus = 1 : cnt = ch - '0';
	while (ch = getc(), ch >= '0' && ch <= '9') cnt = cnt * 10 + ch - '0';
	return minus ? -cnt : cnt;
}

#define maxn 100010
typedef long long ll;
int a[maxn];
int main()
{
//	setfile();
	int n = F(), m = F(), maxx = 0; ll sum = 0;
	for (int i = 1; i <= n; ++i) a[i] = F(), sum += a[i];
	std::sort(a + 1, a + n + 1);
	for (int i = n; i > 1; --i)
		if (a[i - 1] >= a[i] - 1)
		{
			if (a[i] <= 1) a[i - 1] = 1, ++maxx;
			else a[i - 1] = a[i] - 1, ++maxx;
		}
		else maxx += a[i] - a[i - 1];
	maxx += a[1];
	printf("%lld\n", sum - maxx);
	return 0;
}
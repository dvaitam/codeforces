#include <iostream>
#include <cstdio>
#include <algorithm>
#include <queue>
#define gc getchar
#define il inline
#define re register
#define LL long long
using namespace std;
template <typename T>
void rd(T &s)
{
	s = 0;
	bool p = 0;
	char ch;
	while (ch = gc(), p |= ch == '-', ch < '0' || ch > '9');
	while (s = s * 10 + ch - '0', ch = gc(), ch >= '0' && ch <= '9');
	s *= (p ? -1 : 1);
}
template <typename T, typename... Args>
void rd(T &s, Args&... args)
{
	rd(s);
	rd(args...);
}
LL check(int k)
{
	LL p = 1;
	while (k)
		p = p * (k % 10),
		k /= 10;
	return p;
}
const int MAXN = 200000;
int f[MAXN], c[MAXN];
int du[MAXN];
bool p = 0;
int main()
{
	int n;
	rd(n);
	for (int i = 1; i <= n; ++i)
	{
		rd(f[i], c[i]);
		if (f[i] == -1)
			f[i] = 0;
		if (c[i] == 0)
			++du[f[i]];
	}
	for (int i = 1; i <= n; ++i)
	{
		if (du[i] == 0 && c[i] == 1 && f[i])
		{
			printf("%d ", i);
			p = 1;
		}
	}
	if (!p)
		printf("-1");
}
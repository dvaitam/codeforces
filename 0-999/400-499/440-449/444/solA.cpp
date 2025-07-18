#include <cstdio>
#include <algorithm>
#define V_MAX 500
inline int I()
{
	register int c, x = 0;
	do c = getchar(); while (c < '0' || c > '9');
	do x = x * 10 - '0' + c, c = getchar(); while (c >= '0' && c <= '9');
	return x;
}
template <typename _nt>
inline _nt max(_nt x, _nt y)
{
	return x > y ? x : y;
}
int V, E, u, v, w, x[V_MAX + 1];
double ans = 0;
int main()
{
	V = I(), E = I();
	for (u = 1; u <= V; ++u)
		x[u] = I();
	while (E--)
		u = I(), v = I(), w = I(), ans = std::max(ans, (x[u] + x[v]) / double(w));
	printf("%.16f\n", ans);
	return 0;
}
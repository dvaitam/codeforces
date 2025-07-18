#include <cstdio>
#include <iostream>
#include <cmath>
#include <algorithm>
#include <vector>
#include <string>
#include <map>
#include <iterator>
#include <list>
#include <set>
#include <queue>
#include <numeric>
#include <cstdlib>
#include <ctime>
#include <limits>
#include <valarray>
#include <cassert>

#define all(c) (c).begin(), (c).end()

using namespace std;

typedef long long lli;
typedef int li;

template <class T>
bool Maximize (T &v, T nv) { if (nv > v) return v = nv, 1; return 0; }

template <class T>
bool Minimize (T &v, T nv) { if (nv < v) return v = nv, 1; return 0; }

template <class T>
T Mod (T &v, T mod) { return v = (v % mod + mod) % mod; }

const lli INFLL = numeric_limits<lli>::max();
const li INF = numeric_limits<li>::max(), N = 2e5 + 1;
const li di[4][2] = {{1, 0}, {0, 1}, {-1, 0}, {0, -1}};

double real[N], approx[N];
li a[N];

void solve ()
{
	li n, t;
	double c;
	scanf("%d %d %lf", &n, &t, &c);

	for (li i = 1; i <= n; ++ i)
		scanf("%d", a + i);

	lli sum = 0;
	for (li i = 1; i <= n; ++ i)
	{
		sum += a[i];
		if (i >= t)
		{
			sum -= a[i - t];
			real[i] = (double)sum / t;
		}
		approx[i] = (approx[i - 1] + (double)a[i] / t) / c;
	}

	li m;
	scanf("%d", &m);

	for (li i = 0; i < m; ++ i)
	{
		li x;
		scanf("%d", &x);
		printf("%.5f %.5f %.5f\n", real[x], approx[x], abs(approx[x] - real[x]) / real[x]);
	}
}

void init ()
{
	//freopen("input.txt", "r", stdin);
	//freopen("output.txt", "w", stdout);
}

int main()
{
	init();
	solve();
	return 0;
}
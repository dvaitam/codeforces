/*
id: hamed_51
prog: ?
lang: c++
*/

#include <iostream>
#include <iomanip>
#include <fstream>
#include <sstream>
#include <cmath>
#include <cstdio>
#include <string>
#include <vector>
#include <algorithm>
#include <cstdlib>
#include <cstring>
#include <map>
#include <queue>
#include <set>
#include <queue>
#include <stack>
#include <list>
#include <deque>
#include <assert.h>
#include <ctime>
#include <bitset>
#include <numeric>
using namespace std;

#if (_win32 || __win32__)
#define lld "%i64d"
#else
#define lld "%lld"
#endif

#define FOREACH(i, c) for(__typeof((c).begin()) i = (c).begin(); i != (c).end(); ++i)
#define FOR(i, a, n) for (register int i = (a); i < (int)(n); ++i)
#define Size(n) ((int)(n).size())
#define all(n) (n).begin(), (n).end()
#define ll long long
#define pb push_back
#define error(x) cerr << #x << " = " << x << endl;
#define ull unsigned long long
#define point pair<int, int>

#define MAXN 111

bool g[MAXN][MAXN], mark[MAXN];
int n, deg[MAXN];

void dfs(int pos) {
	mark[pos] = true;
	FOR(i, 0, n) if (!mark[i] && g[pos][i]) dfs(i);
}

int main() {
	int m;
	cin >> n >> m;
	FOR(i, 0, m) {
		int u, v;
		cin >> u >> v;
		u--; v--;
		g[u][v] = g[v][u] = true;
		deg[u]++;
		deg[v]++;
	}
	dfs(0);
	FOR(i, 0, n) if (!mark[i] || n != m) {
		cout << "NO" << endl;
		return 0;
	}
	FOR(rep, 0, n) FOR(i, 0, n) if (deg[i] == 1) {
		int f = 0;
		while (!g[i][f]) f++;
		g[i][f] = g[f][i] = false;
		deg[i]--; deg[f]--;
		m--;
	}
	int tot = 0;
	FOR(i, 0, n) if(deg[i] != 0) tot++;
	if (tot == m)
		cout << "FHTAGN!" << endl;
	else
		cout << "NO" << endl;
	return 0;
}
#ifdef _MSC_VER
#define _CRT_SECURE_NO_DEPRECATE
#pragma comment(linker, "/STACK:66777216")
#include <cstdio>
#define GETS gets_s
#else
#pragma GCC optimize("O3")
#pragma GCC target("sse,sse2,sse3,ssse3,sse4,popcnt,abm,mmx,avx")
#include <stdio.h>
#define GETS gets
#endif
#include <algorithm>
#include <string>
#include <set>
#include <map>
#include <vector>
#include <queue>
#include <iostream>
#include <iterator>
#include <cmath>
#include <cstdlib>
#include <sstream>
#include <fstream>
#include <ctime>
#include <cstring>
#include <functional>
#include <chrono>
using namespace std;
#define pb push_back
#define ppb pop_back
#define pi 3.1415926535897932384626433832795028841971
#define mp make_pair
#define x first
#define y second
#define pii pair<int,int>
#define pdd pair<double,double>
#define INF 1000000000
#define FOR(i,a,b) for (int _n(b), i(a); i <= _n; i++)
#define FORD(i,a,b) for(int i=(a),_b=(b);i>=_b;i--)
#define all(c) (c).begin(), (c).end()
#define SORT(c) sort(all(c))
#define rep(i,n) FOR(i,1,(n))
#define rept(i,n) FOR(i,0,(n)-1)
#define L(s) (int)((s).size())
#define C(a) memset((a),0,sizeof(a))
#define VI vector <int>
#define ll long long


struct T {
	int val;
	int num;
	int nx;
};

int a, b, c, d, n, m, k;
int mas[1000002];

int last[1 << 20];
T vnum[1 << 20];
int deg[1 << 20];
pii st[1 << 20];
int path[1 << 20], cur[1 << 20];
bool used[1 << 20];


inline void add(int a, int b, int vv) {
	vnum[m].val = b; vnum[m].num = vv; vnum[m].nx = last[a]; last[a] = m++;
	vnum[m].val = a; vnum[m].num = vv; vnum[m].nx = last[b]; last[b] = m++;
	
	++deg[a]; ++deg[b];
}

inline int fnd(int v, int u) {
	int t = mas[v] & u;
	while (last[t] != -1) {
		int w = last[t];
		last[t] = vnum[w].nx;
		if (!used[vnum[w].val]) return vnum[w].val;
	}
	return -1;
}

void eul(int r) {
	int s = 0;
	st[s++] = mp(r, -1);
	while (s > 0) {
		int v = st[s - 1].x;
		int w = st[s - 1].y;
		if (deg[v] == 0) {
			path[k++] = w;
			--s;
			continue;
		}

		while (last[v] != -1 && used[vnum[last[v]].num]) {
			last[v] = vnum[last[v]].nx;
		}

		used[vnum[last[v]].num] = 1;
		--deg[v]; --deg[vnum[last[v]].val];
		st[s++] = mp(vnum[last[v]].val, vnum[last[v]].num);
		last[v] = vnum[last[v]].nx;
	}

	--k;
}
inline bool check(int len, bool need = 0) {
	if (len > 20) return 0;
	
	m = 0;
	const int u = (1 << len) - 1;

	memset(last, -1, (u + 1) * sizeof(int));
	memset(deg, 0, (u + 1) * sizeof(int));

	for (int i = 0; i < 2 * n; i += 2) {
		add(mas[i] & u, mas[i + 1] & u, i / 2);
	}

	rept(i, u + 1) {
		if (deg[i] % 2) return 0;
	}

	memset(used, 0, n);

	k = 0;
	eul(mas[0] & u);
	if (k < n) return 0;

	if (!need) return 1;

	int s = path[0] * 2;
	if ((mas[s] & u) != (mas[0] & u)) s ^= 1;

	k = 0;
	rept(iter, n) {
		cur[k++] = s;
		cur[k++] = s ^ 1;

		s = s ^ 1;
		if (iter == n - 1) break;
		int t = path[iter + 1];
		if ((mas[2 * t] & u) == (mas[s] & u)) {
			s = 2 * t;
		}
		else {
			s = 2 * t + 1;
		}
	}


	return 1;
}

char buf[8388609], *ch = buf;
inline char get_char() {
	if (*ch == 0) {
		ch = buf;
		fread(buf, sizeof(char), 8388608, stdin);
	}
	return *(ch++);
}
inline void next_int(int& ans) {
	ans = 0;
	char ch;
	while ((ch = get_char()) < '0' || ch > '9');
	do {
		ans = ans * 10 + ch - '0';
	} while ((ch = get_char()) >= '0' && ch <= '9');
}

inline void out_int(int a) {
	int k = 0;
	if (!a) {
		ch[k++] = '0';
	}
	while (a > 0) {
		ch[k++] = (char)(a % 10 + '0');
		a /= 10;
	}

	reverse(ch, ch + k);
	ch += k;
}
int main() {
	//freopen("input.txt", "r", stdin);
	//freopen("output.txt", "w", stdout);

	next_int(n);
	rept(i, n) {
		next_int(mas[2 * i]);
		next_int(mas[2 * i + 1]);
	}

	int l = 0, r = 21;
	while (r - l > 1) {
		int xx = (r + l) / 2;
		if (!check(xx)) r = xx; else l = xx;
	}

	check(l, 1);

	ch = buf;
	out_int(l);
	*ch = 0;
	puts(buf);

	ch = buf;
	rept(i, k) {
		out_int(cur[i] + 1);
		*ch = ' ';
		++ch;
	}
	*ch = 0;
	puts(buf);
}
#include <bits/stdc++.h>

using namespace std;

inline int Read() {
	int x = 0; char y; bool z = false;
	do y = getchar(), z |= y == '-'; while (y < '0' || y > '9');
	do x = x * 10 + y - '0', y = getchar(); while (y >= '0' && y <= '9');
	return z ? -x : x;
}

#define PII pair <int, int>
#define fr first
#define sc second
#define mp make_pair
#define ll long long
#define ld long double
#define INF 0x7f7f7f7f
#define N 200050

int fi[N], c[N * 2][3], ss = 1;
int ln[N][2], n, m, p, Ans[N], ans, d[N];

inline void Line(int x, int y) {
	c[++ss][0] = y; c[ss][1] = fi[x];
	c[fi[x]][2] = ss; fi[x] = ss;
	c[++ss][0] = x; c[ss][1] = fi[y];
	c[fi[y]][2] = ss; fi[y] = ss;
	d[x]++; d[y]++;
	return;
}

queue <int> li;

inline void _Erase(int x) {
	int k = c[x ^ 1][0];
	if (fi[k] == x)
		fi[k] = c[x][1]; else {
			c[c[x][1]][2] = c[x][2];
			c[c[x][2]][1] = c[x][1];
		}
	return;
}

inline void Erase(int x) {
	_Erase(x * 2);
	_Erase(x * 2 + 1);
	return;
}

void Delete(int x) {
	ans--;
	for (int i = fi[x]; i; i = c[i][1])
	 if (d[c[i][0]]-- == p)
	 	li.push(c[i][0]);
	for (int i = fi[x]; i; i = c[i][1])
	 if (d[c[i][0]] >= p)
	 	Erase(i >> 1);
	return;
}

void Pretreat() {
	ans = n;
	for (int i = 1; i <= n; i++)
	 if (d[i] < p)
	 	li.push(i);
	while (!li.empty()) {
		int k = li.front();
		li.pop();
		Delete(k);
	}
	return;
}

void Solve() {
	for (int i = m; i >= 1; i--) {
		Ans[i] = ans;
		int k = ln[i][0], l = ln[i][1];
		if (d[k] < p || d[l] < p)
			continue;
		d[k]--; d[l]--;
		Erase(i);
		if (d[k] >= p && d[l] >= p)
			continue;
		if (d[k] < p)
			Delete(k);
		if (d[l] < p)
			Delete(l);
		while (!li.empty()) {
			int k = li.front();
			li.pop();
			Delete(k);
		}
	}
	return;
}

int main() {
	//freopen("input.txt", "r", stdin);
	n = Read(); m = Read(); p = Read();
	for (int i = 1; i <= m; i++) {
		ln[i][0] = Read();
		ln[i][1] = Read();
		Line(ln[i][0], ln[i][1]);
	}

	Pretreat();
	Solve();
	for (int i = 1; i <= m; i++)
		printf("%d\n", Ans[i]);
	return 0;
}
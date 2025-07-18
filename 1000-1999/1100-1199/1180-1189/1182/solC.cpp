#include <cstdio>
#include <cctype>
#include <vector>
#include <cassert>
#include <cstring>
#include <iostream>
#include <algorithm>
#define R register
#define ll long long
using namespace std;
const int N = 110000, M = N * 10;

int n, tab[200], st[N], num[N], bel[N], x, y, ans[N][2][2], used[N], buc[M], m;
char s[M];
vector<int> app[10];

inline void read(int &x) {
	x = 0;
	char ch = getchar(), w = 0;
	while (!isdigit(ch))
		w = (ch == '-'), ch = getchar();
	while (isdigit(ch))
		x = (x << 1) + (x << 3) + (ch ^ 48), ch = getchar();
	x = w ? -x : x;
	return;
}

inline int read(char *s) {
	int len = 0;
	char ch = getchar();
	while (!isalpha(ch)) ch = getchar();
	while (isalpha(ch))
		s[++len] = ch, ch = getchar();
	s[len + 1] = 0;
	return len;
}

inline bool cmp(int a, int b) {
	return num[a] < num[b];
}

int main() {
	tab['a'] = 1, tab['e'] = 2, tab['i'] = 3, tab['o'] = 4, tab['u'] = 5;
	read(n);
	for (R int i = 1; i <= n; ++i) {
		st[i + 1] = st[i] + read(s + st[i]);
		for (R int j = st[i] + 1; j <= st[i + 1]; ++j)
			if (tab[s[j]])
				++num[i], bel[i] = tab[s[j]];
		buc[num[i]] ^= 1, x += 1 - buc[num[i]];
		if (bel[i])
			app[bel[i]].push_back(i);
	}
	for (R int i = 1; i <= n; ++i)
		buc[num[i]] = 0;
	for (R int i = 1; i <= 5; ++i)
		sort(app[i].begin(), app[i].end(), cmp);
	for (R int i = 1; i <= 5; ++i) {
		int sz = app[i].size();
		for (R int j = 0; j < sz; ++j)
			buc[num[app[i][j]]] ^= 1, y += 1 - buc[num[app[i][j]]];
		for (R int j = 0; j < sz; ++j)
			buc[num[app[i][j]]] = 0;
	}
	printf("%d\n", m = min(y, x >> 1));
	int ind = 0;
	for (R int i = 1; i <= 5; ++i) {
		int sz = app[i].size();
		for (R int j = 0, k; j < sz; ++j) {
			if (k = buc[num[app[i][j]]])
				ans[++ind][0][1] = k, ans[ind][1][1] = app[i][j], used[k] = used[app[i][j]] = 1, buc[num[app[i][j]]] = 0;
			else
				buc[num[app[i][j]]] = app[i][j];
			if (ind >= (x >> 1))
				break;
		}
		for (R int j = 0; j < sz; ++j)
			buc[num[app[i][j]]] = 0;
		if (ind >= (x >> 1))
			break;
	}
	ind = 0;
	for (R int i = 1, j; i <= n; ++i) {
		if (used[i])
			continue;
		if (j = buc[num[i]])
			ans[++ind][0][0] = j, ans[ind][1][0] = i, buc[num[i]] = 0;
		else
			buc[num[i]] = i;
	}
	for (R int i = 1; i <= m; ++i) {
		for (R int j = st[ans[i][0][0]] + 1; j <= st[ans[i][0][0] + 1]; ++j)
			putchar(s[j]);
		putchar(' ');
		for (R int j = st[ans[i][0][1]] + 1; j <= st[ans[i][0][1] + 1]; ++j)
			putchar(s[j]);
		putchar('\n');
		for (R int j = st[ans[i][1][0]] + 1; j <= st[ans[i][1][0] + 1]; ++j)
			putchar(s[j]);
		putchar(' ');
		for (R int j = st[ans[i][1][1]] + 1; j <= st[ans[i][1][1] + 1]; ++j)
			putchar(s[j]);
		putchar('\n');
	}
	return 0;
}
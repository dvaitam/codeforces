#define _CRT_SECURE_NO_WARNINGS

#include <map>
#include <ctime>
#include <cmath>
#include <stack>
#include <cstdio>
#include <cctype>
#include <vector>
#include <bitset>
#include <cstdlib>
#include <cstring>
#include <cassert>
#include <fstream>
#include <iostream>
#include <algorithm>

using namespace std;

typedef long long LL;

inline char gc() {
	static const LL L = 233333;
	static char sxd[L], *sss = sxd, *ttt = sxd;
	if (sss == ttt) {
		ttt = (sss = sxd) + fread(sxd, 1, L, stdin);
		if (sss == ttt) {
			return EOF;
		}
	}
	return *sss++;
}

#define dd c = getchar()
#ifndef dd
#define dd c = gc()
#endif
inline char readalpha() {
	char dd;
	for (; !isalpha(c); dd);
	return c;
}

inline char readchar() {
	char dd;
	for (; c == ' '; dd);
	return c;
}

template <class T>
inline bool read(T& x) {
	bool flg = false;
	char dd;
	x = 0;
	for (; !isdigit(c); dd) {
		if (c == '-') {
			flg = true;
		} else if(c == EOF) {
			return false;
		}
	}
	for (; isdigit(c); dd) {
		x = (x << 1) + (x << 3) + (c ^ 48);
	}
	if (flg) {
		x = -x;
	}
	return true;
}
#undef dd

template <class T>
inline void write(T x) {
	if (x < 0) {
		putchar('-');
		x = -x;
	}
	if (x < 10) {
		putchar(x | 48);
		return;
	}
	write(x / 10);
	putchar((x % 10) | 48);
}

typedef long long LL;

const int maxn = 2005;

int xx[maxn][maxn];

int main() {
	int n;
	read(n);
	int cnt = 0;
	for (int i = 1; i <= n >> 1; ++i) {
		for (int j = 1; j <= n >> 1; ++j) {
			xx[i][j] = 3;
			xx[i][j] |= cnt << 2;
			cnt++;
		}

	}
	cnt = 0;
	for (int i = 1; i <= n >> 1; ++i) {
		for (int j = (n >> 1) + 1; j <= n; ++j) {
			xx[i][j] = 2;
			xx[i][j] |= cnt << 2;
			cnt++;
//			cout << i << ' ' << j << ' ' << cnt << endl;
		}
	}
	cnt = 0;
	for (int i = (n >> 1) + 1; i <= n; ++i) {
		for (int j = 1; j <= n >> 1; ++j) {
			xx[i][j] = 0;
			xx[i][j] |= cnt << 2;
			cnt++;
		}
	}
	cnt = 0;
	for (int i = (n >> 1) + 1; i <= n; ++i) {
		for (int j = (n >> 1) + 1; j <= n; ++j) {
			xx[i][j] = 1;
			xx[i][j] |= cnt << 2;
			cnt++;
		}
	}
	for (int i = 1; i <= n; ++i) {
//		int ans = 0;
		for (int j = 1; j <= n; ++j) {
			write(xx[i][j]);
			putchar(' ');
//			ans ^= xx[i][j];
		}
//		write(ans);
		puts("");
	}
	return 0;
}
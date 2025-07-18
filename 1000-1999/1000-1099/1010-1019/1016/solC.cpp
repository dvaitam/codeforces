#include <bits/stdc++.h>
using namespace std;

inline int read(int u = 0, char c = getchar(), bool f = false) {
	for (;!isdigit(c); c = getchar()) f |= c == '-';
	for (; isdigit(c); c = getchar()) u = u * 10 + c - '0';
	return f ? -u : u;
}

const int maxn = 3e5 + 10;

typedef long long ll;

int a[maxn], b[maxn];

ll s[maxn];

int main() {
	int n = read();
	ll nowL = 0, nowR1 = 0, nowR2 = 0, k = 0, ans = 0;
	for (int i = 1; i <= n; i++) a[i] = read();
	for (int i = 1; i <= n; i++) b[i] = read();
	
	for (int i = n; i >= 1; i--) s[i] = a[i] + b[i] + s[i + 1];
	
	for (int i = 1; i <= n; i++) nowR1 += a[i] * (k++);
	for (int i = n; i >= 1; i--) nowR1 += b[i] * (k++);
	
	k = 0;
	
	for (int i = 1; i <= n; i++) nowR2 += b[i] * (k++);
	for (int i = n; i >= 1; i--) nowR2 += a[i] * (k++);
	
	k = 0;
	
	for (int i = 1; i <= n; i++) {
		if (i & 1) {
			ans = max(ans, nowL + nowR1);
			nowL += a[i] * (k++);
			nowL += b[i] * (k++);
		} else {
			ans = max(ans, nowL + nowR2);
			nowL += b[i] * (k++);
			nowL += a[i] * (k++);
		}
		
		nowR1 -= (n * 2 - 1ll) * b[i] + (i - 1ll) * 2 * a[i];
		nowR1 += s[i + 1];
		
		nowR2 -= (n * 2 - 1ll) * a[i] + (i - 1ll) * 2 * b[i];
		nowR2 += s[i + 1];
	}
	
	cout << ans << endl;
	
}

/*
10 3 4
codeforces
for
1 3
3 10
5 6
5 7
*/
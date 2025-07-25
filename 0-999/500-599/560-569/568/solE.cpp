#include <bits/stdc++.h>
using namespace std;

const int MaxN = 200010, inf = 1 << 30;
int n, m;
int a[MaxN], b[MaxN], nxt[MaxN];
int value[MaxN], vtot;

void input() {
	scanf("%d", &n);
	for (int i = 1; i <= n; ++i) {
		scanf("%d", a + i);
		if (a[i] != -1) value[++vtot] = a[i];
	}
	scanf("%d", &m);
	for (int i = 1; i <= m; ++i) {
		scanf("%d", &b[i]);
		value[++vtot] = b[i];
	}
	sort(b + 1, b + m + 1);
	sort(value + 1, value + vtot + 1);
	vtot = unique(value + 1, value + vtot + 1) - value - 1;
	for (int i = 1; i <= m; ++i)
		b[i] = lower_bound(value + 1, value + vtot + 1, b[i]) - value;
	for (int i = 1; i <= n; ++i)
		if (a[i] != -1) a[i] = lower_bound(value + 1, value + vtot + 1, a[i]) - value;
	for (int i = 0; i <= vtot; ++i) nxt[i] = upper_bound(b + 1, b + m + 1, i) - b;
}

int c[MaxN], pos[MaxN], pre[MaxN];
int val[MaxN];
bool mark[MaxN], used[MaxN];

int main() {
	input();
	memset(c, 63, sizeof(c));
	c[0] = 0;
	int len = 0;
	for (int i = 1; i <= n; ++i)
		if (a[i] != -1) {
			int l = 0, r = len;
			while (l != r) {
				int mid = (l + r + 1) >> 1;
				if (c[mid] < a[i]) l = mid;
				else r = mid - 1;
			}
			if (c[l + 1] > a[i]) {
				c[l + 1] = a[i];
				pos[l + 1] = i;
			}
			pre[i] = pos[l];
			len = max(l + 1, len);
		} else {
			 for (int j = len; j >= 0; --j) {
			 	if (nxt[c[j]] > m) continue;
				if (b[nxt[c[j]]] < c[j + 1]) {
					c[j + 1] = b[nxt[c[j]]];
					pos[j + 1] = pos[j];
					len = max(len, j + 1);
				}
			 }
		}
	int t = 0;
	for (int i = pos[len]; i; i = pre[i]) {
		val[++t] = a[i];
		mark[i] = 1;
	}
	reverse(val + 1, val + t + 1);
	val[t + 1] = inf;
	int now = 0, last = 0, node = 1;
	for (int i = 1; i <= n; ++i) {
		now += mark[i];
		if (a[i] == -1) {
			while ((node <= m) && (b[node] <= last)) ++node;
			if (node > m) break;
			if (b[node] < val[now + 1]) {
				a[i] = b[node];
				used[node] = 1;
				++node;
				last = max(last, a[i]);
			}
		} else if (mark[i]) last = max(last, a[i]);
	}
	node = 1;
	for (int i = 1; i <= n; ++i)
		if (a[i] == -1) {
			while (used[node]) ++node;
			a[i] = b[node];
			++node;
		}
	for (int i = 1; i <= n; ++i) printf("%d ", value[a[i]]); puts("");
	return 0;
}
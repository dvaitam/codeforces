#include <bits/stdc++.h>

typedef long long LL;

#define FOR(i, a, b) for (int i = (a), i##_END_ = (b); i <= i##_END_; ++i)
#define DNF(i, a, b) for (int i = (a), i##_END_ = (b); i >= i##_END_; --i)

template <typename Tp> void in(Tp &x) {
	char ch = getchar(), f = 1; x = 0;
	while (ch != '-' && (ch < '0' || ch > '9')) ch = getchar();
	if (ch == '-') f = -1, ch = getchar();
	while (ch >= '0' && ch <= '9') x = x * 10 + ch - '0', ch = getchar();
	x *= f;
}

template <typename Tp> bool chkmax(Tp &x, Tp y) {return x >= y ? 0 : (x=y, 1);}
template <typename Tp> bool chkmin(Tp &x, Tp y) {return x <= y ? 0 : (x=y, 1);}
template <typename Tp> Tp Max(const Tp &x, const Tp &y) {return x > y ? x : y;}
template <typename Tp> Tp Min(const Tp &x, const Tp &y) {return x < y ? x : y;}

const int MAXN = 300010;

int n, m, s;
int head[MAXN], data[MAXN << 1], nxt[MAXN << 1], id[MAXN << 1], cnt;

int ans1[MAXN], ans2[MAXN];

bool vis[MAXN];

void add(int x, int y, int z, bool tp)
{
	if (tp == 0) {
		nxt[cnt] = head[x]; data[cnt] = y; id[cnt] = 0; head[x] = cnt++;
	}
	else {
		nxt[cnt] = head[x]; data[cnt] = y; id[cnt] = z; head[x] = cnt++;
		nxt[cnt] = head[y]; data[cnt] = x; id[cnt] = -z; head[y] = cnt++;
	}
}

void dfs(int now)
{
	vis[now] = true;
	for (int i = head[now]; i != -1; i = nxt[i]) {
		if (!id[i] && !vis[data[i]]) dfs(data[i]);
	}
}

using std::queue;

queue<int>q;

bool isss[MAXN];

int main()
{
	in(n); in(m); in(s);
	memset(head, -1, sizeof head);

	FOR(i, 1, m) {
		int t, u, v; in(t); in(u); in(v);
		if (t == 1) add(u, v, i, 0);
		else {isss[i] = true; add(u, v, i, 1);}
	}

	dfs(s);

	FOR(i, 1, n) if (vis[i]) q.push(i);

	FOR(i, 1, n) if (vis[i]) {
		for (int j = head[i]; j != -1; j = nxt[j]) {
			if (!vis[data[j]]) {
				if (id[j] > 0)
					ans2[id[j]] = 0;
				else ans2[-id[j]] = 1;
			}
		}
	}

	int A2 = q.size(), A1 = q.size();
	
	while (!q.empty()) {
		int now = q.front(); q.pop();
		for (int i = head[now]; i != -1; i = nxt[i]) {
			if (!vis[data[i]]) {
				vis[data[i]] = true;
				if (id[i]) {
					if (id[i] > 0) {
						ans1[id[i]] = 1;
						//ans2[id[i]] = 0;
					}
					else {
						ans1[-id[i]] = 0;
						//ans2[-id[i]] = 1;
					}
				}
				q.push(data[i]);
				A1++;
			}
		}
	}

	printf("%d\n", A1);
	FOR(i, 1, m) if (isss[i]) {
		printf("%c", ans1[i] ? '+' : '-');
	}

	putchar(10);

	printf("%d\n", A2);
	FOR(i, 1, m) if (isss[i]) {
		printf("%c", ans2[i] ? '+' : '-');
	}

	putchar(10);
	
	return 0;
}
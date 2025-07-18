// LUOGU_RID: 102246479
#include <bits/stdc++.h>

using namespace std;



namespace vbzIO {

	char ibuf[(1 << 20) + 1], *iS, *iT;

	#if ONLINE_JUDGE

	#define gh() (iS == iT ? iT = (iS = ibuf) + fread(ibuf, 1, (1 << 20) + 1, stdin), (iS == iT ? EOF : *iS++) : *iS++)

	#else

	#define gh() getchar()

	#endif

	#define pi pair<int, int>

	#define mp make_pair

	#define fi first

	#define se second

	#define pb push_back

	#define ins insert

	#define era erase

	inline int read () {

		char ch = gh();

		int x = 0;

		bool t = 0;

		while (ch < '0' || ch > '9') t |= ch == '-', ch = gh();

		while (ch >= '0' && ch <= '9') x = (x << 1) + (x << 3) + (ch ^ 48), ch = gh();

		return t ? ~(x - 1) : x;

	}

	inline void write(int x) {

		if (x < 0) {

			x = ~(x - 1);

			putchar('-');

		}

		if (x > 9)

			write(x / 10);

		putchar(x % 10 + '0');

	}

}

using vbzIO::read;

using vbzIO::write;



const int maxn = 4e5 + 400;

int n, m, b, top, st[maxn], vis[maxn], dep[maxn];

vector<int> ans, g[maxn];



void dfs(int u, int fa) {

    dep[u] = dep[fa] + 1, st[++top] = u;

    for (int v : g[u]) {

        if (dep[v]) {

            if (dep[v] > dep[u] - b + 1) continue;

            write(2), puts("");

            vector<int> tp;

            while (st[top] != v) {

                tp.pb(st[top]);

                top--;

            }

            write(tp.size() + 1), puts("");

            for (int i : tp) write(i), putchar(' ');

            write(v);

            exit(0);

        } else dfs(v, u);

    }

    if (!vis[u]) {

        ans.pb(u);

        for (int v : g[u]) vis[v] = 1;

    }

    top--;

}



int main() {

	n = read(), m = read(), b = sqrt(n - 1) + 1;

    for (int i = 1, u, v; i <= m; i++) {

        u = read(), v = read(); 

        g[u].pb(v), g[v].pb(u);

    }

    dfs(1, 0);

    write(1), puts("");

    for (int i = 0; i < b; i++) write(ans[i]), putchar(' ');

	return 0;

}
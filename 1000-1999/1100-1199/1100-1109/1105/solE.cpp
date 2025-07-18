#include <cstdio>
#include <cstring>
#include <algorithm>
#include <string>
#include <map>
using namespace std;

typedef long long ll;

int n, m, cnt, ans, temp;
string A;
char s[50];
int f[50];
ll adj[50];
map <string, int> id;

void dfs(int x, int sz, ll cur) {
    if (x > m) {
	if (sz > temp) temp = sz;
	return;
    }
    if (sz + f[x] <= temp) return;
    
    if (!(adj[x] & cur)) dfs(x + 1, sz + 1, cur | (1LL << (x - 1)));
    dfs(x + 1, sz, cur);
}

int main() {
    scanf("%d %d", &n, &m);
    for (int o = 1; o <= n; ) {
	int tp;
	scanf("%d", &tp);
	if (tp == 1) { ++o; continue; }
	else if (tp == 2) {
	    ll aj = 0;
	    scanf("%s", s);
	    A = s;
	    if (!id[A]) id[A] = ++cnt;
	    aj |= (1LL << (id[A] - 1));
	    
	    for (++o; o <= n; ++o) {
		scanf("%d", &tp);
		if (tp == 1) { ++o; break; }
		scanf("%s", s);
		A = s;
		if (!id[A]) id[A] = ++cnt;
		aj |= (1LL << (id[A] - 1));
	    }

	    for (int i = 1; i <= m; ++i)
		if (aj & (1LL << (i - 1))) adj[i] |= aj;
	}	
    }

    for (int i = m; i >= 1; --i) {
	temp = 0;
	dfs(i + 1, 1, 1LL << (i - 1));
	f[i] = max(f[i + 1], temp);
	//printf("%d ", f[i]);
    }

    printf("%d\n", f[1]);
    return 0;
}
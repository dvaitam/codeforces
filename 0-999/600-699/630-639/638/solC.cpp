// Claris you are our red sun, I can not live without you.
// NanoApe you are our red sun, I can not live without you.
// YJQ you are our red sun, I can not live without you.
#include <bits/stdc++.h>
using namespace std;
typedef long long ll;
typedef long double ld;
typedef pair<int, int> pii;
typedef pair<ll, int> pli;
typedef pair<int, ll> pil;
typedef pair<ll, ll> pll;
inline int getint() {
	int x=0, f=1, c=getchar();
	for(; c<48||c>57; f=c=='-'?-1:f, c=getchar());
	for(; c>47&&c<58; x=x*10+c-48, c=getchar());
	return x*f;
}
const int N=200005;
int mx, ihead[N], cnt;
struct E {
	int next, to, i;
}e[N<<1];
vector<int> ans[N];
void add(int x, int y, int i) {
	e[++cnt]=(E){ihead[x], y, i}; ihead[x]=cnt;
	e[++cnt]=(E){ihead[y], x, i}; ihead[y]=cnt;
}
void dfs(int x, int last=0, int f=0) {
	int tot=0;
	for(int i=ihead[x]; i; i=e[i].next) {
		int y=e[i].to;
		if(y!=f) {
			++tot;
			if(tot==last) {
				++tot;
			}
			ans[tot].push_back(e[i].i);
			dfs(y, tot, x);
		}
	}
	mx=max(mx, tot);
}
int main() {
	int n=getint();
	for(int i=1; i<n; ++i) {
		add(getint(), getint(), i);
	}
	dfs(1);
	printf("%d\n", mx);
	for(int i=1; i<=mx; ++i) {
		int g=ans[i].size();
		printf("%d", g);
		for(int j=0; j<g; ++j) {
			printf(" %d", ans[i][j]);
		}
		puts("");
	}
	return 0;
}
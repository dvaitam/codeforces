//#define __USE_MINGW_ANSI_STDIO 0
#include <iostream>
#include <bits/stdc++.h>
using namespace std;

#define debug(x) cout<<#x<<" :: "<<x<<endl;
#define debug2(x,y) cout<<#x<<" :: "<<x<<"\t"<<#y<<" :: "<<y<<endl;
#define debug3(x,y,z) cout<<#x<<" :: "<<x<<"\t"<<#y<<" :: "<<y<<"\t"<<#z<<" :: "<<z<<endl;

#define boost ios::sync_with_stdio(0); cin.tie(0)

#define fi first
#define se second
#define pb(x) push_back(x)
#define mp(x,y) make_pair(x,y)

typedef long long ll;
typedef double ld;
typedef pair<int, int> pii;

const int N = 1e5 + 5;
const int X = 18;
const int MOD = 1e9 + 7;

/***************************************************************************/

int tim;
int st[N], en[N];
int par[N][X], dep[N];
vector<int> adj[N];

void dfs(int v, int p) {
	st[v] = en[v] = ++tim;
	dep[v] = dep[p] + 1;
	par[v][0] = p;
	for(int i=1; i<X; i++) {
		par[v][i] = par[par[v][i-1]][i-1];
	}
	for(auto it : adj[v]) {
		if(it == p) continue;
		dfs(it, v);
	}
	en[v] = tim;
}

int __lca(int a, int b) {
	if(dep[b] > dep[a]) swap(a, b);
	for(int i=X-1; i>=0; i--) {
		if(dep[a]-(1<<i) >= dep[b]) {
			a = par[a][i];
		}
	}
	if(a == b) return a;

	for(int i=X-1; i>=0; i--) {
		if(par[a][i] != par[b][i]) {
			a = par[a][i], b = par[b][i];
		}
	}
	return par[a][0];
}

int isThere[N];
int A[N], L[N];
vector<pii> grp[N];

vector<int> radj[N];

void make_tree(vector<pii> &V, int l, int r, int p) {
	if(l > r) return;
	int v = V[l].se; if(p != v) radj[p].pb(v);
	int p1 = l+1, p2 = upper_bound(V.begin()+l, V.begin()+r+1, mp(en[v]+1, 0)) - V.begin() - 1;
	make_tree(V, p1, p2, v);
	make_tree(V, p2+1, r, p);
}

queue<pii> Q;
ll dp[305], tdp[305];

int main() {

	boost;
	int n, q; cin>>n>>q;
	for(int i=1; i<n; i++) {
		int a, b; cin>>a>>b;
		adj[a].pb(b), adj[b].pb(a);
	}
	dfs(1, 0);
	while(q--) {
		int k, m, r; cin>>k>>m>>r;
		vector<pii> V;

		for(int i=1; i<=k; i++) {
			cin>>A[i];
			isThere[A[i]] = 1;

			int l = __lca(A[i], r);
			V.pb(mp(dep[l], l));

			grp[l].pb(mp(st[A[i]], A[i]));
		}

		sort(V.begin(), V.end());
		V.resize(unique(V.begin(), V.end()) - V.begin());

		for(auto it : V) {
			sort(grp[it.se].begin(), grp[it.se].end());
			make_tree(grp[it.se], 0, grp[it.se].size()-1, it.se);
		}

		for(int i=1; i<(int)V.size(); i++) {
			radj[V[i].se].pb(V[i-1].se);
		}
		radj[0].pb(V.back().se);

		memset(dp, 0, sizeof dp);
		dp[0] = 1;

		Q.push(mp(0, 0));
		while(Q.size()) {
			int v = Q.front().fi, d = Q.front().se;
			Q.pop();

			for(auto it : radj[v]) {
				Q.push(mp(it, d+isThere[v]));
			}
			if(isThere[v] == 0) continue;

			for(int i=1; i<=m; i++) {
				tdp[i] = dp[i-1];
				if(d < i) {
					tdp[i] = (tdp[i] + dp[i] * (i-d)) % MOD;
				}
			}

			for(int i=0; i<=m; i++) {
				dp[i] = tdp[i]; tdp[i] = 0;
			}
		}

		ll ans = 0;
		for(int i=0; i<=m; i++) {
			ans += dp[i];
		}
		ans %= MOD;
		cout<<ans<<"\n";

		for(int i=1; i<=k; i++) {
			isThere[A[i]] = 0;
			radj[A[i]].clear();
		}
		for(auto it : V) {
			grp[it.se].clear();
			radj[it.se].clear();
		}
		radj[0].clear();
	}
	return 0;
}
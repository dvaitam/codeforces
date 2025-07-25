#include <cstdio>
#include <cstdlib>
#include <cstring>
#include <cmath>
#include <iostream>
#include <string>
#include <vector>
#include <map>
#include <set>
#include <queue>
#include <stack>
#include <algorithm>

using namespace std;

#define REP(i,n) for(int(i)=0;(i)<(int)(n);(i)++)
#define SZ(c) ((int)(c).size())
#define ITER(it,c) for(__typeof((c).begin()) it=(c).begin();it!=(c).end();it++)
#define FIND(x,c) ((c).find((x))!=(c).end())
#define MP(x,y) make_pair((x),(y))

typedef long long ll;
const int MAXN=1005;
const int MAXM=5005;
int n,m,u[MAXM],v[MAXM],s[MAXN];
bool vis[MAXN],reach[MAXN];
vector <int> adj[MAXN];

void dfs(int u)
{
    vis[u]=true;
    if (u==n-1) reach[u]=true;
    ITER(it,adj[u]) {
	if (!vis[*it]) dfs(*it);
	if (reach[*it]) reach[u]=true;
    }
}

int main()
{
    scanf("%d%d",&n,&m);
    REP(i,m) {
	scanf("%d%d",u+i,v+i);
	u[i]--,v[i]--;
	adj[u[i]].push_back(v[i]);
    }
    dfs(0);
    bool updated=true;
    for (int l=1; l<=n && updated; l++) {
	updated=false;
	REP(j,m) {
	    if (reach[u[j]] && reach[v[j]]) {
		if (s[u[j]]-s[v[j]]<1) {
		    s[v[j]]=s[u[j]]-1;
		    updated=true;
		} else if (s[u[j]]-s[v[j]]>2) {
		    s[u[j]]=s[v[j]]+2;
		    updated=true;
		}
	    }
	}
    }
    if (updated) puts("No");
    else {
	puts("Yes");
	REP(i,m) printf("%d\n",reach[u[i]] && reach[v[i]] ? s[u[i]]-s[v[i]]:1);
    }
    return 0;
}
#include<bits/stdc++.h>

#define ll long long

using namespace std;

const int N=3e3+5;



bool vis[N];

vector<int>e[N];

int n,p[N],a[N],b[N],x[N],y[N];

bool dfs(int now) {

	if (vis[now]) {

		return 0;

	}

	vis[now]=1;

	for (int to:e[now]) {

		if (!y[to]||dfs(y[to])) {

			x[now]=to,y[to]=now;

			return 1;

		}

	}

	return 0;

}



signed main() {

	ios::sync_with_stdio(false),cin.tie(0);

	cout.precision(10),cout.setf(ios::fixed);

	

	cin>>n;

	for (int i=1;i<=n;i++) {

		cin>>p[i];

	}

	int n1=0;

	for (int i=1;i<=n;i++) {

		if (vis[i]) {

			continue;

		}

		n1++;

		for (int j=i;!vis[j];j=p[j]) {

			a[j]=n1;

			vis[j]=1;

		}

	}

	for (int i=1;i<=n;i++) {

		cin>>p[i];

	}

	int n2=0;

	memset(vis,0,sizeof(vis));

	for (int i=1;i<=n;i++) {

		if (vis[i]) {

			continue;

		}

		n2++;

		for (int j=i;!vis[j];j=p[j]) {

			b[j]=n2;

			vis[j]=1;

		}

	}

	for (int i=1;i<=n;i++) {

		e[a[i]].emplace_back(b[i]);

	}

	int ans=n;

	for (int i=1;i<=n1;i++) {

		memset(vis,0,sizeof(vis));

		ans-=dfs(i);

	}

	cout<<ans<<"\n";

	for (int i=1;i<=n;i++) {

		if (x[a[i]]==b[i]) {

			x[a[i]]=0;

		} else {

			cout<<i<<' ';

		}

	}

	

	return 0;

}
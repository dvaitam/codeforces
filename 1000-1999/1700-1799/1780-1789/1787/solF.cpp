#include<bits/stdc++.h>

using namespace std;

#define IOS ios_base::sync_with_stdio(false);cin.tie(0);cout.tie(0);

#define pb push_back

#define SZ(x) (int)x.size()

typedef vector<int> VI;

int power(int a, int b, int p) {

    int res = 1 % p;

    for (; b > 0; b /= 2, a = 1LL * a * a % p) {

        if (b % 2 == 1) {

            res = 1LL * res * a % p;

        }

    }

    return res;

}

void solve(){

	int n,k;cin>>n>>k;

	VI a(n);

	for (int i=0;i<n;++i){

		cin>>a[i];a[i]--;

	}

	

	VI vis(n),ans(n);

	vector<vector<VI>> f(n+1);

	for (int i=0;i<n;++i){

		if (vis[i]) continue;

		int j=i;

		VI t;

		while (!vis[j]) t.pb(j),vis[j]=1,j=a[j];

		f[SZ(t)].pb(t);

	}

	

	for (int i=1;i<=n;++i)

	 if (i%2==0 && SZ(f[i])){

	 	if (k>20 || SZ(f[i])%(1<<k)){

	 		cout<<"NO\n";

	 		return;

		 }

    }

    

	cout<<"YES\n";	 

	for (int i=1;i<=n;++i){

		while (!f[i].empty()){

			int t=min(k,__lg(SZ(f[i])));

			int pw=power(2,k-t,i);

			

			vector tmp(f[i].end() - (1<<t) , f[i].end());

			f[i].erase(f[i].end() - (1<<t) , f[i].end());;

			

			for (auto &g:tmp){

				VI ng(i);

				for (int j=0;j<i;++j)

				 ng[1LL*pw*j%i]=g[j];

				g=ng;

			}

			

			int all=i<<t;

			for (int j=0;j<all;++j)

			 ans[tmp[j%(1<<t)][j/(1<<t)]  ] = tmp [(j+1)%all%(1<<t)] [(j+1)%all/(1<<t)];

			

		}

	}

	for (int i=0;i<n;++i) cout<<ans[i]+1<<" \n"[i==n-1];

	

}



signed main(){

	IOS

	int T;cin>>T;

	while (T--) solve();

	return 0-0;

}
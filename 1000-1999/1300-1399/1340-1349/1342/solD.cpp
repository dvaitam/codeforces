#include<bits/stdc++.h>

using namespace std;

const int N=1e6+10,R=1e6;

const int mod=1e9+7;





void solve() {

	int n,k;

	cin>>n>>k;

	vector<int>a(k+1),b(k+1);

	for(int i=1;i<=n;i++){

		int x;

		cin>>x;

		a[x]++;

	}

	for(int i=1;i<=k;i++)cin>>b[i];

	long long suf=0,num=-1;

	for(int i=k;i>=1;i--){

		suf+=a[i];

		num=max(num,suf/b[i]+(suf%b[i]==0?0:1));

	}

	vector<vector<int> >ans(num);

	for(int i=1,id=-1;i<=k;i++){

		for(int j=1;j<=a[i];j++){

			id++;

			ans[id%num].push_back(i);

		}

	}

	cout<<num<<'\n';

	for(int i=0;i<num;i++){

		cout<<ans[i].size()<<" ";

		for(int v:ans[i])cout<<v<<" ";

		cout<<'\n';

	}

}

int main() {

	ios::sync_with_stdio(false);

	cin.tie(0);

	cout.tie(0);

	int T=1;

	while(T--) {

		solve();

	}

	return 0;

}

/*



*/
#pragma G++ optimize("Ofast")

#pragma G++ optimize("unroll-loops")

#include<iostream>

#include<algorithm>

#include<vector>

#include<cstring>

#include<functional>

#include<queue>

#include<unordered_map>

#include<map>

#include<set>

#include<stack>

#include<cmath>

#include<bitset>

#include<iomanip>

#include<numeric>



using namespace std;

using ll=long long;

using ld=long double;

using P=pair<ll,ll>;

const int INF=1e9;

const ll inf=1e18;



void solve()

{

	int n,k; cin>>n>>k;

	vector<int>nxt(k);

	cout<<'a';  int now=0;

	for(int i=2;i<=n;i++)

	{

		cout<<(char)((nxt[now])+'a');

		int p=nxt[now]; nxt[now]++;

		if(nxt[now]>=k)

		{

			for(int j=now+1;j<k;j++)nxt[j]=max(nxt[j],now+1);

		}

		now=p;

		if(nxt[k-1]>=k)nxt.assign(k,0);

	}

	cout<<"\n";

}	



int main()

{

	ios::sync_with_stdio(false);

	cin.tie(0); cout.tie(0);

	int t=1; //cin>>t;

	while(t--)solve();

	return 0;

}
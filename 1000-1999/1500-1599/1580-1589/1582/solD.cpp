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

	int n; cin>>n;

	vector<int>a(n+1);



	for(int i=1;i<=n;i++)

	{

		cin>>a[i];

	}



	if(n&1^1)

	{

		for(int i=1;i<=n;i+=2)

		{

			cout<<-a[i+1]<<" "<<a[i]<<" ";

		}

		cout<<"\n";

	}

	else

	{

		for(int i=1;i<=n-3;i+=2)

		{

			cout<<-a[i+1]<<" "<<a[i]<<" ";

		}

		if(a[n-2]+a[n-1]!=0)

		{

			cout<<-a[n]<<" "<<-a[n]<<" "<<(a[n-2]+a[n-1]);

		}

		else if(a[n-2]+a[n]!=0)

		{

			cout<<-a[n-1]<<" "<<a[n-2]+a[n]<<" "<<-a[n-1];

		}

		else if(a[n-1]+a[n]!=0)

		{

			cout<<a[n-1]+a[n]<<" "<<-a[n-2]<<" "<<-a[n-2];

		}

		cout<<"\n";

	}



}	



int main()

{

	ios::sync_with_stdio(false);

	cin.tie(0); cout.tie(0);

	int t=1; cin>>t;

	while(t--)solve();

	return 0;

}
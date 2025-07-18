#include<bits/stdc++.h>

typedef long long ll;

typedef unsigned long long ull;

using namespace std;

const ll mod = 998244353;

const int mm = 1e5 + 10;

int main(){

	std::ios::sync_with_stdio(false);

	std::cin.tie(0);

	std::cout.tie(0);

	ll x;

	cin>>x;

	pair<ll,ll>ans={1,x};

	vector<ll>v;

	for(int i=2;i<=x/i;i++){

		if(x%i==0){

			ll js=1;

			while(x%i==0){

				x/=i;

				js*=i;

			}

			v.push_back(js);

		}

	}if(x>1)v.push_back(x);

	for(ll i=1;i<=(1ll<<v.size())-1;i++){

		ll bj=i;

		vector<int>zy;

		while(bj||zy.size()<v.size()){

			zy.push_back(bj%2);

			bj/=2;

		}

		ll bjl=1,bjr=1;

		for(int j=0;j<v.size();j++){

			if(zy[j])bjl*=v[j];

			else bjr*=v[j];

		}

		if(max(bjl,bjr)<max(ans.first,ans.second)){

			ans={bjl,bjr};

		}

	}

	ll xx=ans.first,yy=ans.second;

	cout<<min(xx,yy)<<' '<<max(xx,yy)<<'\n';

}
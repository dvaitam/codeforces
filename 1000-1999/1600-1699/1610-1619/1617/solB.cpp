#include<iostream>

#include<cstdio>

#include<algorithm>

#include<cstring>

#include<cassert>

#define siz(x) int((x).size())

#define cauto const auto

#define all(x) (x).begin(),(x).end()

#define fi first

#define se second

using std::cin;using std::cout;

using loli=long long;

using uloli=unsigned long long;

using lodb=long double;

using venti=__int128_t;

using pii=std::pair<int,int>;

using inlsi=const std::initializer_list<int>&;

constexpr venti operator""_vt(uloli x){return venti(x);}

signed main(){

//	freopen(".in","r",stdin);

//	freopen(".out","w",stdout);

	std::ios::sync_with_stdio(false);cin.tie(nullptr);

	int T;cin>>T;for(int n;T--;){

		cin>>n;

		if(~n&1)cout<<n/2-1<<' '<<n/2<<' '<<1<<'\n';

		else if((n-1)%4==0)cout<<(n-1)/2-1<<' '<<(n-1)/2+1<<' '<<1<<'\n';

		else cout<<(n-1)/2-2<<' '<<(n-1)/2+2<<' '<<1<<'\n';

	}

	return 0;

}
#include<bits/stdc++.h>
#define vi vector<int>
#define vll vector<long long>
#define ll long long
#define ull unsigned long long
#define FAST_IO ios_base::sync_with_stdio(false); cin.tie(NULL)
#define nl cout<<'\n'
#define pb push_back
#define mp make_pair
#define pi pair<int,int>
#define pll pair<long long,long long>
#define vpii vector<pair<int,int> >
#define vpll vector<pair<ll,ll> >
#define si set<int>
#define sll set<long long>
#define all(v) (v).begin(),(v).end()
#define Unique(x)           (x).erase(unique(all(x)), x.end())
#define v_present(vec,x)  (find(all(vec),x)!=(vec).end())
#define present(myset,x)    ((myset).find(x)!=(myset).end())
#define Cprint(c)  for(auto i:(c)) { cout<<i<<" "; } nl
#define trace2(x,y) cout<<x<<" "<<y<<endl;
#define trace3(x,y,z) cout<<x<<" "<<y<<" "<<z<<endl;
#define trace4(x,y,z,a) cout<<x<<" "<<y<<" "<<z<<" "<<a<<endl;
ll mod = 1e9+7;
const int MAX = 1e5+5;
using namespace std;

int main()
{
FAST_IO;
int q;
cin>>q;
for(int i=0;i<q;i++){
	ll l,r;
	cin>>l>>r;
	
		ll a=ceil((long double)(l-1)/2.0);
		ll a1=a;
		ll b=ceil((long double)r/2.0);
		ll b1=b;
		if((l-1)%2!=0)
		a1=-a;
		if(r%2!=0)
		b1=-b;
		cout<<b1-a1;
		nl;
		
	
}
return 0;
}
#include<bits/stdc++.h>
using namespace std;
#define fastio ios_base::sync_with_stdio(false); cin.tie(NULL); cout.tie(NULL)
using ll=long long;
using ld=long double;
#define INF 1000000007
const ll INFll=1ll*INF*INF;
#define MOD 998244353
#define MODE 1000000007
const long double PI=3.141592653589793238462643383279502884197169399375105820974944;
#define pb push_back
#define mp make_pair
#define rep(i,n) for (ll i = 0; i < n; i++)
#define rev(i,n) for (ll i = n-1; i >= 0; i--)
#define repe(i,n)   for(ll (i)=1;(i)<=(n);(i)++)
#define reve(i,n)   for(ll (i)=n;(i)>=1;(i)--)
#define all(v) v.begin(),v.end()
#define f first
#define s second
#define mii map<int,int>
#define vi vector<int>
#define vl vector<ll>
#define vli vector<pair<ll,int>>
#define pll pair<ll,ll>
#define vll vector<pll >
#define vvl vector<vector<ll> >
#define pii pair<int,int>
#define vii vector< pii >
#define vvi vector< vector<int > >
#define vvii vector< vector<pii > >
#define W(t) while(t --)
#define print(arr) for (auto it = arr.begin(); it != arr.end(); ++it) cout << *it << ' '; cout << endl;
#define printii(arr) for (auto it = arr.begin(); it != arr.end(); ++it) cout << it->f<<' '<<it->s << endl; cout << endl;
#define MID(i,j) ((i)+(j))/2
#define nl '\n' 
#define lcm(a,b) ((a)*(b))/gcd((a),(b))
#define gcd(a,b){ return (b==0)? a:gcd(b,a%b);} 

int main()
{
	fastio;

	ll n,k;
	cin >>n >>k;

	ll sig=(k*(k+1))/2;
	if(k==1)
	{
			cout << "YES\n";
			cout << n << nl;
			return 0;
	}
	if(n<sig)
	{
		cout << "NO\n";
		return 0;
	}
	n-=sig;
	ll q=n/k,r=n%k;
	if(q>0 ||(q==0 && r!=k-1))
	{
			cout << "YES\n";
			repe(i,k-1)
			{
				cout << i+q<<' ';
			}
			cout << k+q+r <<nl;
			return 0;
	}
	if(k>=4)
	{
			cout << "YES\n";
			repe(i,k-2)
			{
				cout << i+q<<' ';
			}
			cout << k+q <<' ';
			cout << k+q+r-1 <<nl;
		
	}
	else
		cout << "NO\n";

	

	
}
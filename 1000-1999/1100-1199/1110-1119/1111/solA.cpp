#include <bits/stdc++.h>
using namespace std;

#define pb push_back
#define mp make_pair
#define fr first
#define sc second
#define clr(a) memset(a, 0, sizeof(a))
#define rep(n) for(int i=0; i<n; i++)
#define repc(i, n) for(int i=0; i<n; i++)
#define rep1(n) for(int i=1; i<=n; i++)
#define repc1(i, n) for(int i=1; i<=n; i++)
#define all(v) v.begin(), v.end()
#define alla(a,n) a, a+n
//#define endl "\n" 
typedef long long ll;
typedef unsigned long long ull;
typedef vector<ll> vi;
#define PI 3.14159265

ll divn(ll a, ll b)
{
	ll x = a/b;
	if(a%b)
	{
		x++;
	}
	return x;
}
void input(){ios_base::sync_with_stdio(false);cin.tie(NULL);}
//.............................................................
bool vow(char c)
{
	if(c=='a' || c=='u' || c=='e' || c=='i' || c=='o')
		return true;
	return 0;
	
}

int main(int argc, char const *argv[])
{
	input();
	string s, t;
	cin>>s>>t;
	if(s.length()!=t.length())
	{
		cout<<"No";
		return 0;
	}
	for(int i=0; s[i]!='\0' && t[i]!='\0'; i++)
	{
		//~ cout<<i;
		if(vow(s[i])==vow(t[i]))
		{
			continue;
		}
		
		cout<<"No";
		return 0;
	}
	cout<<"Yes";
	return 0;
}
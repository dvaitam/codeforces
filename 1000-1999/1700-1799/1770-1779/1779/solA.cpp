#include<bits/stdc++.h>

using namespace std;

#define f_(i,a,b) for(int i=a;i>=b;i--)

#define f(i,a,b) for(int i=a;i<b;i++)

#define all(v) v.begin(),v.end()

#define maxm(a,b) a=max(a,b)

#define minm(a,b) a=min(a,b)

#define lst(x) x[x.size()-1]

#define sz(x) int(x.size())

#define mp make_pair

#define pb push_back

#define S second

#define F first

#define ps ios_base::sync_with_stdio(0),cin.tie(0),cout.tie(0);

typedef long long ll;

typedef long double ld;

typedef double dbl;

const int maxn=1e5+5;

int n,m,ans;

bool fl;

string s;

int main(){

	ps

	int t;

	cin>>t;

	while(t--){

		cin>>n>>s;

		fl=true,ans=maxn;

		f(i,0,n-1){

			if(s[i]=='R'&&s[i+1]=='L'){fl=0,ans=0;break;}

			if(s[i]=='L'&&s[i+1]=='R'&&fl)fl=0,ans=i+1;

		}

		(fl)?cout<<-1<<'\n':cout<<ans<<'\n';

	}

	return 0;

}
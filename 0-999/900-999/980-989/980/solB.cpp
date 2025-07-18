/* ****hellojarvis**** */

#include<bits/stdc++.h>
#define ll          long long
#define pb          push_back
#define pii         pair<ll int,ll int>
#define vi          vector<ll int>
#define all(a)      (a).begin(),(a).end()
#define F           first
#define S           second
#define sz(x)       (ll int)x.size()
#define hell        1000000007
#define rep(i,a,b)	for(ll int i=a;i<b;i++)
#define lbnd        lower_bound
#define ubnd        upper_bound
#define bs          binary_search
#define mp          make_pair
#define io	ios_base::sync_with_stdio(false);cin.tie(0);cout.tie(0);
using namespace std;

int main()
{
	io
	ll i,j,n,k;
	cin>>n>>k;
	string s;
	s="";
	if(k%2==0)
	{
		cout<<"YES\n";
		
		cout<<string(n,'.')<<"\n";
		
		s+=".";
		s+=string(k/2,'#');
		s+=string(n-1-(k/2),'.');
		s+="\n";
		cout<<s<<s;
		
		cout<<string(n,'.')<<"\n";

	}
	else if(k>=(n-2))
	{
		cout<<"YES\n";
		
		cout<<string(n,'.')<<"\n";
		s+=".";
		s+=string(n-2,'#');
		s+='.';
		cout<<s<<"\n";

		s="";
		s+='.';
		s+=string( (k-(n-2))/2, '#' );
		s+=string(n-2-(k-(n-2)),'.');
		s+=string( (k-(n-2))/2, '#' );
		s+='.';
		cout<<s<<"\n";
		
		cout<<string(n,'.')<<"\n";



	}
	else
	{
		cout<<"YES\n";

		cout<<string(n,'.')<<"\n";

		s+=string((n-k)/2,'.');
		s+=string(k, '#');
		s+=string((n-k)/2,'.');
		s+='\n';
		cout<<s;
		cout<<string(n,'.')<<"\n";

		cout<<string(n,'.')<<"\n";


	}




	return 0;
}
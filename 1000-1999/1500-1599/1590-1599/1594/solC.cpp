#include<bits/stdc++.h>

using namespace std;

typedef long long L;

typedef pair<int,int> P;

const int N=3e5+2;

void so()

{

	int n;char a;string s;

	cin>>n>>a>>s;

	int c=0,x;

	for(int i=0;i<n;i++)

	{

		if(s[i]==a)

		{

			c++;

			x=i+1;

		}

	}//cout<<c<<'*';

	if(c==n) cout<<0<<'\n';

	else if(x>n/2||c==n-1) cout<<1<<'\n'<<x<<'\n';

	else cout<<2<<'\n'<<n<<' '<<n-1<<'\n';

}

int main()

{

	ios::sync_with_stdio(0);

	cin.tie(0);cout.tie(0);

	int q=1;

	cin>>q;

	while(q--) so();

	return 0;

}
#include <bits/stdc++.h>

#define lli long long int

using namespace std;

int up2 (int n)

{

	return (n>>1)+(n&1);

}

int main () 

{

    ios_base::sync_with_stdio(false);

    cin.tie(NULL); cout.tie(NULL);

    int t; cin>>t;

    while (t--)

    {

    	int n; cin>>n;

    	cout<<"2\n";

    	int cur=n;

    	for (int i=n-1;i>=1;i--)

    	{

    		cout<<i<<" "<<cur<<"\n";

    		cur=up2(i+cur);

		}

	}

}
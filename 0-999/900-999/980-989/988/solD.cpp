#include <bits/stdc++.h>
using namespace std;
typedef vector<int> vi;
typedef long long int lli;
typedef pair<lli, int> pli;
typedef vector<lli> vli;
#define arr_input(a,n) \
	for(int i=0;i<n;i++) {cin>>a[i];}

int main()
{
	ios_base::sync_with_stdio(false);
	cin.tie(nullptr);
	cout.tie(nullptr);
	
	int n;
	cin>>n;
	vli a(n);
	arr_input(a,n);
	
	sort(a.begin(),a.end());	
	
	bool ans2 = false;
	pair<int,int> ans= make_pair(-1,-1);
	for(auto pos = a.begin(); pos != a.end();pos++)
	{
		lli x = *pos;
		lli diff = 1;
		lli prevdiff = -1;
		
		while((x+diff) <= a.back())
		{
			auto searchpos = lower_bound(pos,a.end(),x+diff);
			if((searchpos != a.end()) && (*searchpos == (x+diff)))//found
			{
				if(prevdiff == diff/2)
				{
					cout<<"3\n";
					cout<<x<<" "<<(x+prevdiff)<<" "<<(x+diff)<<endl;
					return 0;
				}
				else
				{
					prevdiff = diff;
				}					
			}
			
			diff <<= 1;
		}
		if(prevdiff != -1)
		{
			ans2 = true;
			ans.first = x;
			ans.second = x+prevdiff;
		}				
	}
	
	if(ans2)
	{
		cout<<"2\n";
		cout<<ans.first<<" "<<ans.second;
	}
	else
	{
		cout<<"1\n";
		cout<<a[0]<<endl;
	}
}
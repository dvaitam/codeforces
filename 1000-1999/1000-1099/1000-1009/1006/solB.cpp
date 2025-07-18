#include <bits/stdc++.h>
using namespace std;

int main()
{
	int n,k;
	cin>>n>>k;
	vector<int> v;
	vector<int> temp;
	for(int i=0; i<n; i++)
	{
		int x;
		cin>>x;
		v.push_back(x);
	}

	temp = v;
	sort(temp.begin(),temp.end(),greater<int>());

	int sum=0;
	vector<int> maxv;
	for(int i=0; i<k; i++)
	{
		sum+=temp[i];
		maxv.push_back(temp[i]);
	}

	temp.clear();

	int i;
	int prev = -1;
	for(i=0; i<n; i++)
	{
		if(maxv.empty())
		{
			break;
		}

		for(int j=0; j<maxv.size(); j++)
		{
			if(v[i]==maxv[j])
			{
				temp.push_back(i-prev);
				prev = i;
				maxv.erase(maxv.begin()+j);
				break;
			}
		}
	}
	temp[temp.size()-1]+=n-i;
	cout<<sum<<endl;
	for(int i=0; i<temp.size(); i++)
	{
		cout<<temp[i]<<" ";
	}


}
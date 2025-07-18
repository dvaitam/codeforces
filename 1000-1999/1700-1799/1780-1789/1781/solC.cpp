#include <bits/stdc++.h>

#define ll long long

using namespace std;

struct fairy

{

	int pos,val;

};

bool cmp(fairy x,fairy y)

{

	return x.val>y.val;

}

void solve()

{

	int n;

	cin>>n;

	vector<int> sum(26);

	vector<char> s(n+1);

	for (int i=1;i<=n;i++)

	{

		cin>>s[i];

		sum[s[i]-'a']++;

	}

	

	int minn=1e9,ans;

	for (int num=1;num<=n;num++)

	{

		if (n%num!=0||n/num>26)

			continue;

		int grp=n/num;

		

		vector<fairy> temp;

		for (int i=0;i<26;i++)

		{

			if (sum[i])

				temp.push_back({i,sum[i]});

		}

		sort(temp.begin(),temp.end(),cmp);

		

		int cnt=0;

		if (temp.size()<grp)

		{

			for (auto i:temp)

			{

				if (i.val>num)

					cnt+=i.val-num;

				else

					break;

			}

		}

		else

		{

			for (int i=0;i<temp.size();i++)

			{

				if (i<grp&&temp[i].val>num)

					cnt+=temp[i].val-num;

				if (i>=grp)

					cnt+=temp[i].val;

			}

		}

		if (cnt<minn)

			minn=cnt,ans=num;

	}

	

	cout<<minn<<"\n";

	int grp=n/ans;

	

	vector<fairy> temp;

	for (int i=0;i<26;i++)

	{

		if (sum[i])

			temp.push_back({i,sum[i]});

	}

	sort(temp.begin(),temp.end(),cmp);

	

	vector<int> good(26);

	for (int i=0;i<temp.size();i++)

	{

		if (i<grp&&temp[i].val<=ans)

			good[temp[i].pos]=1;

		else

			good[temp[i].pos]=2;

	}

	

	for (int i=1;i<=n;i++)

	{

		int ch=s[i]-'a';

		if (good[ch]==2&&((sum[ch]>0&&sum[ch]<ans)||sum[ch]>ans))

		{

			for (int j=0;j<26;j++)

			{

				if (j==ch)

					continue;

				if (good[j]==1&&sum[j]<ans)

				{

					s[i]=j+'a';

					sum[j]++;

					sum[ch]--;

					break;

				}

			}

		}

	}

	for (int i=1;i<=n;i++)

	{

		int ch=s[i]-'a';

		if (good[ch]==2&&sum[ch]>ans)

		{

			for (int j=0;j<26;j++)

			{

				if (good[j]==0&&sum[j]<ans)

				{

					sum[j]++;

					sum[ch]--;

					s[i]=j+'a';

					break;

				}

			}

		}

	}

	

	

	for (int i=1;i<=n;i++)

		cout<<s[i];

	cout<<"\n";

}

int main()

{

	ios::sync_with_stdio(0);

    cin.tie(0);

    cout.tie(0);

    

    int t;

    cin>>t;

	while (t--)

	{

		solve();

	}

	return 0;

}

/*

1

20

sasaasalssedsalaqalq

*/
#include<bits/stdc++.h>

#define endl '\n'

using namespace std;

typedef long long i64;

typedef int i32;

typedef long double ld;



void like(void)

{

	i64 a,b,c;

	cin>>a>>b>>c;

	if(a*b<=c)

	{

		if(a==c&&b==1)

		{

			cout<<"-1 -1"<<endl;

			return ;

		}

		cout<<"1 -1"<<endl;

		return;

	}

	if(a>=c)

	{

		cout<<-1<<' '<<b<<endl;

		return;

	}

	cout<<1<<' '<<b<<endl;

}



int main()

{

	ios::sync_with_stdio(false);

	cin.tie(nullptr);

	cout.tie(nullptr);

	int t;

	cin>>t;

	while(t--)

   		like();

	return 0;

}
#include<bits/stdc++.h>

using namespace std;

int main()

{

	int n,k,a,b,c,d;

	cin>>n>>k>>a>>b>>c>>d;

	if(k<=n || n<=4 ) return cout<<-1,0;

	for(int i=0;i<2;i++)

	{

		cout<<a<<" "<<c<<" ";

		for(int j=1;j<=n;j++)

		{

			if(j!=a && j!=b && j!= c && j!=d)

			{

				cout<<j<<" ";

			}

		}

		cout<<d<<" "<<b<<endl;

		swap(a,c);

		swap(b,d);

	}

}
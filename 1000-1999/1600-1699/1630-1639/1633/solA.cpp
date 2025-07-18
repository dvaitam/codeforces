#include<bits/stdc++.h>

using namespace std;

int main()

{

	int t;

	cin>>t;

	while(t--)

	{

		int n;

		cin>>n;

		if(n%7==0)

		cout<<n<<endl;

		else

		{

			for(int i=1;i<10;i++)

			{

				if((n-(n%10)+i)%7==0)

				{

					cout<<n-n%10+i<<endl;

					break;

				}

			}

		}

	}

}
#include<bits/stdc++.h>
#include<math.h>
using namespace std;
int main()
{
	long n,s;
	long q=1e10,i,j,t;
	cin>>n;
	int a[n];
	cin>>a[0];
	long min=a[0],max=a[0];
	for(long  i=1;i<n;i++)
	{
	  cin>>a[i];
	  if(a[i]<min)
	  min=a[i];
	  if(max<a[i])
	  max=a[i];
	  }

//	cout<<"min="<<min<<" max"<<max<<endl;
	for( i=min;i<=max;i++)
	{
		long l=0;
		for( j=0;j<n;j++)
		{
			if(abs(i-a[j])<=1)
			;
			else
			l=l-1+abs(i-a[j]);
		}//cout<<"l="<<l<<" i="<<i<<endl;
		if(l<q)
		{
			q=l;
			t=i;
		}
	}
    cout<<t<<" "<<q;
	return 0;	
}
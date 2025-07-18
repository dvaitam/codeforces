#include<iostream>
#include<math.h>
#include<stdio.h>
using namespace std;
int main()
{
	int n,h;
	while(cin>>n>>h)
	{
		double *out=new double[n-1];
		double *out1=new double[n-1];
		out[0]=sqrt(2);
		for(int i=1;i<n-1;i++)
		{
			out[i]=sqrt(2-(1.0/(out[i-1]*out[i-1])));
		}
	
		double ou;
	
		out1[n-2]=h/out[n-2];
		for(int i=n-3;i>=0;i--)
		{
		    
		     out1[i]=out1[i+1]/out[i];
		}
		for(int i=0;i<n-1;i++)
		{
			printf("%.12f ",out1[i]); 
		}
		
		
	} 
	
	return 0;
}
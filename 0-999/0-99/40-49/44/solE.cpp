#include <iostream>
#include <cstring>
#include <cstdio>
#include <cctype>
#include <string> 
using namespace std; 
int main()
{
	int k,a,b,len; 
	string str; 
	scanf("%d%d%d",&k,&a,&b); 
	cin>>str; 
	len=str.size(); 
	if (!(len>=a*k&&len<=b*k))
	{
		printf("No solution\n"); 
		return 0;
	}
	for (int ix=0;ix<len;ix++)
	{
		double M=(len-ix)/(double)k;
		for (int z=1;z<=(int)M;z++,ix++)
			cout<<str[ix]; 
		ix--; 
		cout<<endl;
		k--;
	}
}
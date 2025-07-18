#include<cstdio>
#include<cstring>
#include<cmath>
using namespace std;
int has[100001];
bool stay[100001];
int num[100001];

#define youhua __attribute__((optimize("O2")))
char c;
youhua inline  void read(int &a)
{
	a=0;
	do
	 c=getchar();
	while(c<'0'||c>'9');
	while(c>='0'&&c<='9')
	 {
	 	a*=10;
	 	a+=c-'0';
	 	c=getchar();
	 }
}

 inline int find_nextj(int x)
{
	do 
	 x++;
	while(has[x]);
	return x;
}

youhua int main()
{
	
	int n,i,j,k;
	read(n);
	memset(has,0,sizeof(has));
	for(i=1;i<=n;i++)
	 read(num[i]),has[num[i]]=i;
	for(i=1;i<=n;i++)
         stay[has[i]]=true;
	j=0;
	j=find_nextj(j);
	for(i=1;i<=n;i++)
	 if(stay[i])
	    printf("%d ",num[i]);
	else
	   printf("%d ",j),j=find_nextj(j);
	
	return 0;
}
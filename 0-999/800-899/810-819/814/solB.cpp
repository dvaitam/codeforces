#include<stdlib.h>
#include<string.h>
#include<stdio.h>
using namespace std;
int main()
{
    int n;
    while(scanf("%d",&n)!=EOF)
    {
		int a[1005],b[1005];
        for(int i=0;i<n;i++)
		    scanf("%d",&a[i]);
        for(int i=0;i<n;i++)
		    scanf("%d",&b[i]);
        int x[5],y[5],k=0;
		int num[1100];
		memset(num,0,sizeof(num));		
		for(int i=0;i<n;i++)
		{
			if(a[i]!=b[i])
			{
				x[k]=a[i];
				y[k]=b[i];
				a[i]=-1;
				k++;
			}
			else
			{
				num[a[i]]=1;
			}
		}
		int f[5],key=0;
		for(int i=1;i<=n;i++)
		{
			if(num[i]==0)
			f[key++]=i;
		}
		if(key==1)
		{
             for(int i=0;i<n;i++)
			 {
				if(a[i]!=-1)
				{
					printf("%d%c",a[i],i==n-1?'\n':' ');
				}
				else
				printf("%d%c",f[0],i==n-1?'\n':' ');
			  }
		}
		else
		{
			if((f[0]==x[0]&&f[0]!=y[0]&&f[1]==y[1]&&f[1]!=x[1])||(f[0]!=x[0]&&f[0]==y[0]&&f[1]!=y[1]&&f[1]==x[1]))
		    {
				int s=0;
				for(int i=0;i<n;i++)
				{
					if(a[i]!=-1)
					{
				    	printf("%d%c",a[i],i==n-1?'\n':' ');
					}
					else
					{
						printf("%d%c",f[s],i==n-1?'\n':' ');
						s++;
					}
				}
			}
			else  if((f[1]==x[0]&&f[1]!=y[0]&&f[0]==y[1]&&f[0]!=x[1])||(f[1]!=x[0]&&f[1]==y[0]&&f[0]!=y[1]&&f[0]==x[1]))
			{
				int s=1;
				for(int i=0;i<n;i++)
				{
					if(a[i]!=-1)
					{
						printf("%d%c",a[i],i==n-1?'\n':' ');
					}
					else
					{
						printf("%d%c",f[s],i==n-1?'\n':' ');
						s--;
					}
				}
			}
		}
    }
    return 0;
}
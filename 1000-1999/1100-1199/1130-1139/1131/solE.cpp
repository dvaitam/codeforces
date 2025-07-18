#include <bits/stdc++.h>
using namespace std;
char s[100010];
int maxlen[26];
void chkmax(int &a,int b){a=max(a,b);}
void chkmin(int &a,int b){a=min(a,b);}
inline int read()
{
    int x=0,t=1;
    char ch=getchar();
    while((ch<'0'||ch>'9')&&ch!='-')
    {
        ch=getchar();
    } 
    if(ch=='-')
    {
        t=-1;
        ch=getchar();
    }
    while(ch>='0'&&ch<='9')
    {
        x=x*10+ch-48;
        ch=getchar();
    }
    return x*t;
}
int main()
{
	int myst;
	myst=read();
	for(int k=1;k<=myst;++k)
	{
		scanf("%s",s+1);
		int n=strlen(s+1);
		if(k!=1)
		{
			int l=1;
			while(l<n&&s[l+1]==s[1])
			{
				l++;
			}
			int r=n;
			while(r>1&&s[r-1]==s[n])
			{
				r--;
			}
			r=n-r+1;
			if(l==n)
			{
				for(int i=0;i<26;++i)
				{
					if(i!=s[1]-'a')
					{
						chkmin(maxlen[i],1);
					}
				}
				if(maxlen[s[1]-'a'])
				{
					maxlen[s[1]-'a']=(maxlen[s[1]-'a']+1)*(n+1)-1;
				}
				else 
				{
					maxlen[s[1]-'a']=n;
				}
			}
			else
			{
				for(int i=0;i<26;++i)
				{
					chkmin(maxlen[i],1);
				}
				maxlen[s[1]-'a']+=l;
				maxlen[s[n]-'a']+=r;
			}
		}
		for(int i=1,j;i<=n;i=j+1)
		{
			j=i;
			while(j<n&&s[j+1]==s[i])
			{
				j++;	
			}
			chkmax(maxlen[s[i]-'a'],j-i+1);
		}
	}
	int ans=0;
	for(int i=0;i<26;++i)
	{
		chkmax(ans,maxlen[i]);
	}
	printf("%d\n",ans);
	return 0;
}
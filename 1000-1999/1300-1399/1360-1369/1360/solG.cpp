#include<stdio.h>
#include<string.h>
inline void read(int&x)
{
	register char c=getchar();for(;c<'0'||'9'<c;c=getchar());
	for(x=0;'0'<=c&&c<='9';x=(x<<3)+(x<<1)+(c^48),c=getchar());
}
int t,n,m,a,b;bool ans[55][55],ok;
main()
{
	for(read(t);t--;)
	{
		read(n);read(m);read(a);read(b);ok=1;memset(ans,0,sizeof(ans));
		for(register int i=0,j=0;i<n;++i)
			for(register int k=0;k<a;++k)
			{
				ans[i][j++]=1;
				if(j==m)j=0;
			}
		for(register int i=0,cnt;i<m;++i)
		{
			cnt=0;
			for(register int j=0;j<n;++j)if(ans[j][i])++cnt;
			if(cnt!=b){ok=0;break;}
		}
		if(ok)
		{
			puts("YES");
			for(register int i=0;i<n;++i)
			{
				for(register int j=0;j<m;++j)putchar(ans[i][j]?'1':'0');
				putchar('\n');
			}
		}else puts("NO");
	}
}
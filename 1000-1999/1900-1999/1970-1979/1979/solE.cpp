#include<stdio.h>
#include<map>
#include<vector>
#include<algorithm>
#define pr pair<int,int> 
using namespace std;
// inline char nc()
// {
// 	static char buf[99999],*l,*r;
// 	return l==r&&(r=(l=buf)+fread(buf,1,99999,stdin),l==r)?EOF:*l++;
// }
#define nc getchar
inline void read(int&x)
{
	bool t=0;char c=nc();for(;c<'0'||'9'<c;t|=c=='-',c=nc());
	for(x=0;'0'<=c&&c<='9';x=(x<<3)+(x<<1)+(c^48),c=nc());if(t)x=-x;
}
int t,n,m;map<pr,int>a;map<int,vector<pr> >mmp1,mmp2;
struct __readt__{inline __readt__(){read(t);}}_readt___;
main()
{
	read(n);read(m);m>>=1;a.clear();mmp1.clear();mmp2.clear();
	for(int i=1,x,y;i<=n;++i)
	{
		read(x),read(y);a[(pr){x,y}]=i;
		mmp1[x+y].emplace_back(x,i);
		mmp2[x-y].emplace_back(x,i);
	}
	for(map<int,vector<pr> >::iterator it=mmp1.begin();it!=mmp1.end();++it)
		sort(it->second.begin(),it->second.end());
	for(map<int,vector<pr> >::iterator it=mmp2.begin();it!=mmp2.end();++it)
		sort(it->second.begin(),it->second.end());
	for(map<pr,int>::iterator it=a.begin();it!=a.end();++it)
	{
		int x=it->first.first,y=it->first.second;
		if(a.count((pr){x+m,y+m}))
		{
			vector<pr>&tmp=mmp2[(x+m)-(y-m)];
			int i=lower_bound(tmp.begin(),tmp.end(),(pr){x+m,0})-tmp.begin();
			if(i<tmp.size())if(tmp[i].first<=x+m+m)
			{
				printf("%d %d %d\n",it->second,a[(pr){x+m,y+m}],tmp[i].second);
				goto nxt;
			}
		}
		if(a.count((pr){x+m,y+m}))
		{
			vector<pr>&tmp=mmp2[(x-m)-(y+m)];
			int i=lower_bound(tmp.begin(),tmp.end(),(pr){x-m,0})-tmp.begin();
			if(i<tmp.size())if(tmp[i].first<=x)
			{
				printf("%d %d %d\n",it->second,a[(pr){x+m,y+m}],tmp[i].second);
				goto nxt;
			}
		}
		if(a.count((pr){x+m,y-m}))
		{
			vector<pr>&tmp=mmp1[(x+m)+(y+m)];
			int i=lower_bound(tmp.begin(),tmp.end(),(pr){x+m,0})-tmp.begin();
			if(i<tmp.size())if(tmp[i].first<=x+m+m)
			{
				printf("%d %d %d\n",it->second,a[(pr){x+m,y-m}],tmp[i].second);
				goto nxt;
			}
		}
		if(a.count((pr){x+m,y-m}))
		{
			vector<pr>&tmp=mmp1[(x-m)+(y-m)];
			int i=lower_bound(tmp.begin(),tmp.end(),(pr){x-m,0})-tmp.begin();
			if(i<tmp.size())if(tmp[i].first<=x)
			{
				printf("%d %d %d\n",it->second,a[(pr){x+m,y-m}],tmp[i].second);
				goto nxt;
			}
		}
	}
	printf("0 0 0\n");
	nxt:if(--t)main();
}
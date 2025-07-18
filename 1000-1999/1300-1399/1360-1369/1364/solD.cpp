#include <bits/stdc++.h>

using namespace std;

char buf[1<<23],*p1=buf,*p2=buf;

#define getchar() (p1==p2&&(p2=(p1=buf)+fread(buf,1,1<<21,stdin),p1==p2)?EOF:*p1++)

template <typename T>

inline void read(T &f)

{

    f=0;T fu=1;char c=getchar();

    while(c<'0'||c>'9') {if(c=='-'){fu=-1;}c=getchar();}

    while(c>='0'&&c<='9') {f=(f<<3)+(f<<1)+(c&15);c=getchar();}

    f*=fu;

}

template <typename T> 

void print(T x,char c=0)

{

    if(x<0) putchar('-'),x=-x;

    if(x<10) putchar(x+48);

    else print(x/10),putchar(x%10+48);

    if(c) putchar(c);

}

inline void reads(string &f)

{

    string str="";char ch=getchar();

    while(ch<'!'||ch>'~') ch=getchar();

    while((ch>='!')&&(ch<= '~')) str+=ch,ch=getchar();

    f=str;

}

void prints(string s)

{

    for(int i=0;s[i];++i) 

    putchar(s[i]);

}

typedef long long ll;

const int multicase=0,debug=0,maxn=2e5+50;

vector<int> G[maxn];

int n,m,k;

vector<int> st;

int dfn[maxn];

void dfs2(int u,int fa)

{

	st.push_back(u);

	dfn[u]=(int)st.size();

	for(auto v:G[u])

    {

		if(v!=fa)

		{

			if(!dfn[v])

			{

				dfs2(v,u);

				continue;

			}

			if(dfn[u]-dfn[v]+1<=k&&dfn[u]+1>=dfn[v])

			{

				print(2,'\n');

                print(dfn[u]-dfn[v]+1,'\n');

				for(int j=dfn[v];j<=dfn[u];++j)

				print(st[j-1],' ');

				exit(0);

			}

        }

    }

	st.pop_back();

}

int col[maxn];

vector<int> col1,col2;

void dfs1(int u,int fa)

{

    if(col[u]==1) col1.push_back(u);

    else col2.push_back(u);

    if(col1.size()==(k+1)/2)

    {

        for(auto x:col1)

        print(x,' ');

        exit(0);

    }

    if(col2.size()==(k+1)/2)

    {

        for(auto x:col2)

        print(x,' ');

        exit(0);

    }

    for(auto v:G[u])

    {

        if(v==fa) continue;

        if(col[v]==0)

        {

            if(col[u]==1) col[v]=2;

            else col[v]=1;

            dfs1(v,u);

        }

    }

}

void solve()

{

    read(n),read(m),read(k);

    for(int i=1;i<=m;++i)

    {

        int u,v;

        read(u),read(v);

        G[u].push_back(v),G[v].push_back(u);

    }

    dfs2(1,0);

    print(1,'\n');

    col[1]=1;

    dfs1(1,0);

}

int main()

{

    #ifdef AC

    freopen("in.txt","r",stdin);

    freopen("out.txt","w",stdout);

    #endif

    clock_t program_start_clock=clock();

    int _=1;

    if(multicase) read(_);

    while(_--) solve();

    fprintf(stderr,"\nTotal Time: %lf ms",double(clock()-program_start_clock)/1000);

    return 0;

}
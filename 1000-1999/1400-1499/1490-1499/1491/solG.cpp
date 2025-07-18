#include<bits/stdc++.h>

using namespace std;

typedef int i_t;

typedef int o_t;

constexpr size_t iSIZE=1<<23,oSIZE=1<<23;

char buf1[iSIZE],buf2[oSIZE],*p1=buf1,*p2=buf2;

inline i_t get(){

	i_t x=0,f=1;

	while(*p1<'0'||*p1>'9')if(*(p1++)=='-')f=-1;

	while(*p1>='0'&&*p1<='9')x=x*10+(*(p1++)^48);

	return f*x;

}

inline void write(const char *s){

	while(*s)*(p2++)=*(s++);

}

inline void write(o_t x,char c){

	static int sta[39],*top=sta;

	if(x<0)x=-x,*(p2++)='-';

	do *(top++)=int(x%10),x/=10; while(x);

	while(top!=sta)*(p2++)=char(*(--top)^48);

	*(p2++)=c;

}

struct file{

	file(){

		fread(buf1,1,iSIZE,stdin);

	}

	~file(){

		fwrite(buf2,1,p2-buf2,stdout);

	}

}file_;

constexpr int N=200001;

int c[N];

vector<vector<int>> g,G;

bitset<N> vis;

int pos[N];

vector<pair<int,int>> ops;

inline void Swap(int x,int y){

	ops.emplace_back(pos[x],pos[y]);

	swap(c[pos[x]],c[pos[y]]);

	swap(pos[x],pos[y]);

}

inline void two_cycles(int p1,int p2){

	

	while(!vis[c[p1]]&&c[p1]!=p2){

		int np1=c[p1];

		vis[p1]=1;

		Swap(p1,c[p1]);

		p1=np1;

	}

	while(!vis[c[p2]]&&c[p2]!=p1){

		int np2=c[p2];

		vis[p2]=1;

		Swap(p2,c[p2]);

		p2=np2;

	}

	Swap(p1,p2);

	vis[p1]=vis[p2]=1;

}

int main(){

	int n=get();

	for(int i=1;i<=n;++i){

		pos[c[i]=get()]=i;

	}

	for(int i=1;i<=n;++i){

		if(!vis[i]){

			G.emplace_back();

			for(int u=i;!vis[u];u=c[u]){

				G.back().emplace_back(u);

				vis[u]=1;

			}

		}

	}

	for(auto &v:G){

		if(v.size()>1){

			g.emplace_back(v);

		}

	}

	vis.reset();

	while(g.size()>=2){

		//cerr<<"OK "<<g.size()<<'\n';

		auto &g1=g.back(),&g2=g[g.size()-2];

		Swap(g1.back(),g2.back());

		int p1=g1.back(),p2=g2.back();

		two_cycles(p1,p2);

		g.pop_back(),g.pop_back();

	}

	if(g.size()==1){

		auto &g1=g.back();

		if(g1.size()==n){

			int a1=g1[0],a2=g1[1],a3=g1[2];

			Swap(a1,a2),Swap(a1,a3);

			two_cycles(a2,a3);

		}

		else{

			int node=0;

			for(int i=1;i<=n;++i){

				if(c[i]==i){

					node=i;

					break;

				}

			}

			vis[node]=0;

			Swap(g1.back(),node);

			two_cycles(g1.back(),node);

		}

	}

	write(ops.size(),'\n');

	for(auto [x,y]:ops){

		write(x,' ');

		write(y,'\n');

	}

}
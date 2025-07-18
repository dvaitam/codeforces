#include <bits/stdc++.h>

using namespace std;

typedef long long LL;



#define clr(a,x) memset(a,x,sizeof a)

const int MAXN=50005;

const int MAXE=100005;

namespace Fast_Read{

	const int SIZE=1<<16;

	char buf[SIZE],c;

	int l=0,r=0;

	inline char readchar(){

		if(l==r){l=0,r=fread(buf,1,SIZE,stdin);}

		return buf[l++];

	}

	inline void scanf(int&x){

		while((c=readchar())<48);

		x=c-48;

		while((c=readchar())>47)x=x*10+c-48;

	}

	inline void scanfchar(char&x){

		while((x=readchar())<'A');

	}

}

struct Node *null;

struct Node{

	Node *c[2],*f;

	LL cnt,sum;

	int val,siz,lazy,tree_siz;

	inline void new_node(){sum=val=lazy=0,c[0]=c[1]=f=null;}

	inline void add_siz(int SIZE){

		if(this==null)return;

		tree_siz+=SIZE;

		cnt+=2LL*siz*SIZE;

		lazy+=SIZE;

	}

	inline void up(){sum=c[0]->sum+(LL)siz*val+c[1]->sum;}

	inline void down(){if(lazy)c[0]->add_siz(lazy),c[1]->add_siz(lazy),lazy=0;}

	inline bool is_root(){return f==null||f->c[0]!=this&&f->c[1]!=this;}

	inline void sign_down(){if(!is_root())f->sign_down();down();}

	inline void setc(Node *o,bool d){c[d]=o,o->f=this;}

	inline void rot(bool d){

		Node *p=f,*g=f->f;

		p->setc(c[d],!d);

		if(!p->is_root())g->setc(this,p==g->c[1]);

		else f=g;

		setc(p,d);

		p->up();

	}

	inline void splay(){

		for(sign_down();!is_root();){

			if(f->is_root())rot(this==f->c[0]);

			else{

				if(f==f->f->c[0])this==f->c[0]?f->rot(1):rot(0),rot(1);

				else this==f->c[1]?f->rot(0):rot(1),rot(0);

			}

		}

		up();

	}

	inline void access(){

		for(Node*o=this,*x=null;o!=null;x=o,o=o->f){

			if(x!=null){

				while(x->c[0]!=null)x=x->c[0];

				x->splay();

			}

			o->splay();

			o->siz=o->tree_siz-x->tree_siz;//withoutsplay_treeright_sontree_siz

			o->setc(x,1);

			o->up();

		}

		splay();

	}

}pool[MAXN],*node[MAXN],*curN;;

struct Edge{

	int v;

	Edge *next;

}edge[MAXE],*H[MAXN],*curE;

LL ans;

int n,m;

void init(){

	ans=0;

	curE=edge;

	curN=pool;

	null=curN++;

	null->new_node();

	null->cnt=null->tree_siz=null->siz=0;

	clr(H,0);

}

void addedge(int u,int v){

	curE->v=v;

	curE->next=H[u];

	H[u]=curE++;

}

void dfs(int u){

	node[u]->cnt=node[u]->tree_siz=1;

	for(Edge *e=H[u];e;e=e->next){

		int v=e->v;

		dfs(v);

		node[u]->cnt+=2LL*node[u]->tree_siz*node[v]->tree_siz;

		node[u]->tree_siz+=node[v]->tree_siz;

	}

	ans+=node[u]->cnt*node[u]->val;//calculateinitans

	node[u]->siz=node[u]->tree_siz;

}

void solve(){

	int x,y,i;

	char op;

	double tot=(double)n*n;

	init();

	for(i=1;i<=n;++i){

		node[i]=curN++;

		node[i]->new_node();

	}

	for(i=2;i<=n;++i){

		Fast_Read::scanf(x);

		addedge(x,i);

		node[i]->f=node[x];

	}

	for(i=1;i<=n;++i)Fast_Read::scanf(node[i]->val);

	dfs(1);

	printf("%.10f\n",ans/tot);

	Fast_Read::scanf(m);

	while(m--){

		Fast_Read::scanfchar(op);

		Fast_Read::scanf(x);

		Fast_Read::scanf(y);

		if(op=='P'){

			node[y]->access();

			node[x]->splay();

			if(node[x]->f==null)swap(x,y);

			//cut

			node[x]->access();

			Node *o=node[x]->c[0];

			ans-=2LL*o->sum*node[x]->tree_siz;

			o->add_siz(-node[x]->tree_siz);

			node[x]->c[0]->f=null;

			node[x]->c[0]=null;

			node[x]->up();

			o->up();

			//link

			node[y]->access();

			ans+=2LL*node[y]->sum*node[x]->tree_siz;

			node[y]->add_siz(node[x]->tree_siz);

			node[y]->up();

			node[x]->f=node[y];

		}else{

			node[x]->access();

			ans+=node[x]->cnt*(y-node[x]->val);

			node[x]->val=y;

			node[x]->up();

		}

		printf("%.10f\n",ans/tot);

	}

}

int main(){

	Fast_Read::scanf(n);

	solve();

	return 0;

}
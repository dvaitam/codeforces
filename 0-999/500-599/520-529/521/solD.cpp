#include <bits/stdc++.h>
#include<cstdio>

#include<queue>

#include<vector>

#include<algorithm>

using namespace std;

const int N=1e5+5;

int num[N],assign_mx[N],assign_id[N];

typedef int64_t ll;

struct Add{

	int type,val,id;

	ll a,b;

	// ratio:a/b

	inline bool operator <(const Add &tmp)const{

		return val>tmp.val;

	}

};

vector<Add>add[N];

struct Mul{

	int type,id,add_id1,add_id2;

	ll a,b;

	inline bool operator <(const Mul &tmp)const{

		return a*tmp.b<b*tmp.a;

	}

};

priority_queue<Mul>pque;

struct Ans{

	int type,id;

	inline bool operator <(const Ans &tmp)const{

		return type<tmp.type;

	}

}ans[N];

void rd(int &res){

	res=0;

	char c;

	while(c=getchar(),c<48);

	do res=(res<<3)+(res<<1)+(c^48);

		while(c=getchar(),c>47);

}

int main(){

	int n,m,K;

	rd(n);rd(m);rd(K);

	for(int i=1;i<=n;++i){

		rd(num[i]);

		assign_mx[i]=num[i];

		add[i].clear();

	}

	while(!pque.empty())pque.pop();

	for(int i=1,type,tar,val;i<=m;++i){

		rd(type);rd(tar);rd(val);

		if(type==1){

			if(val>assign_mx[tar]){

				assign_mx[tar]=val;

				assign_id[tar]=i;

			}

		}

		else if(type==2)add[tar].push_back((Add){2,val,i});

		else pque.push((Mul){3,i,-1,-1,val-1,1});

	}

	for(int i=1;i<=n;++i){

		if(assign_mx[i]>num[i])

			add[i].push_back((Add){1,assign_mx[i]-num[i],assign_id[i]});

		if(!add[i].size())continue;

		sort(add[i].begin(),add[i].end());

		ll sum=num[i];

		for(int j=0;j<add[i].size();++j){

			add[i][j].b=sum;

			sum+=add[i][j].val;

			add[i][j].a=add[i][j].val;

		}

		pque.push((Mul){add[i][0].type,add[i][0].id,i,0,add[i][0].a,add[i][0].b});

	}

	int tot=0;

	while(K--){

		if(pque.empty())break;

		Mul cur=pque.top();

		if(cur.a<=0)break;

		pque.pop();

		ans[tot++]=(Ans){cur.type,cur.id};

		if(cur.type<3){

			int id1=cur.add_id1,id2=cur.add_id2;

			if(id2<add[id1].size()-1)pque.push((Mul){add[id1][id2+1].type,add[id1][id2+1].id,id1,id2+1,add[id1][id2+1].a,add[id1][id2+1].b});

		}

	}

	printf("%d\n",tot);

	sort(ans,ans+tot);

	for(int i=0;i<tot;++i)

		printf("%d%c",ans[i].id,i==tot-1?'\n':' ');

	return 0;

}
#include<cstdio>

#include<cstring>

#include<algorithm>

#include<queue>

#define rep(i,a,b) for(i=a;i<=b;i++)

using namespace std;

const int N=3010,INF=1e9;

struct node{

	int w,pos;

	friend bool operator < (node c,node d){

		return c.w<d.w;

	}

};

priority_queue<node> h1,h2;

priority_queue<node> q;

struct Node{

	int a,b,pos;

};Node x[N];

int b[N];

int ans;

int t1[N],t2[N];

int sum[N];



bool cmp(Node c,Node d){

	return c.a>d.a;

}



int main(){

	//freopen("a.in","r",stdin);

	//freopen("a.out","w",stdout);

	int i;

	int n,p,s;

	node tmp,tmp1,tmp2,tmp3;

	scanf("%d%d%d",&n,&p,&s);

	//printf("%d %d %d\n",n,p,s);

	rep(i,1,n){

		scanf("%d",&x[i].a);

	}

	rep(i,1,n){

		scanf("%d",&x[i].b);

		x[i].pos=i;

		sum[i]=x[i].b-x[i].a;

	}

	sort(x+1,x+1+n,cmp);

	/*rep(i,1,n){

		printf("%d %d %d\n",x[i].a,x[i].b,x[i].pos);

	}*/

	rep(i,1,p){

		tmp.w=x[i].b-x[i].a;tmp.pos=x[i].pos;

		q.push(tmp);

		ans+=x[i].a;

		t1[x[i].pos]=1;

	}

	//printf("%d\n",ans);

	rep(i,p+1,n){

		tmp.w=x[i].a;tmp.pos=x[i].pos;

		h1.push(tmp);

		tmp.w=x[i].b;tmp.pos=x[i].pos;

		h2.push(tmp);

	}

	int w1,w2,w3,pos1,pos2,pos3;

	rep(i,1,s){

		do{

			tmp1=h1.top();w1=tmp1.w;pos1=tmp1.pos;

			if(b[pos1]){

				h1.pop();

			}

		}while(b[pos1]);		

		do{

			tmp2=h2.top();w2=tmp2.w;pos2=tmp2.pos;

			if(b[pos2]){

				h2.pop();

			}

		}while(b[pos2]);

		tmp3=q.top();w3=tmp3.w;pos3=tmp3.pos;

		//printf("w1=%d w2=%d w3=%d\n",w1,w2,w3);

		w1+=w3;

		if(w1>w2){

			q.pop();t1[pos3]=0;t2[pos3]=1;

			tmp3.w=sum[pos1];tmp3.pos=pos1;

			q.push(tmp3);

			b[pos1]=1;t1[pos1]=1;

			h1.pop();

			ans+=w1;

		}

		else{

			h2.pop();t2[pos2]=1;

			b[pos2]=1;

			ans+=w2;

		}

	}

	printf("%d\n",ans);

	rep(i,1,n){

		if(t1[i])	printf("%d ",i);

	}

	printf("\n");

	rep(i,1,n){

		if(t2[i])	printf("%d ",i);

	}

	printf("\n");

}
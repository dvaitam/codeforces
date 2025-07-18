//08:54

#include<bits/stdc++.h>

#define N 5100

using namespace std;

int ta=0,tid;

struct st{

    int w,h,id;

};

struct tr{

    int s,e,id,num;

};

int jilu[N];

int from[N];

st all[N];

int en;

bool cmp(st a,st b){

    return a.h>b.h;

}

tr tree[N<<2];

void build(int s,int e,int k){

    tr &a=tree[k];

    a.s=s;a.e=e;a.num=-1;a.id=0;

    if(s==e)return;

    int mid=(s+e)>>1;

    build(s,mid,k<<1);

    build(mid+1,e,k<<1|1);

}

void fi(int s,int e,int k){

    if(s>e)return;

    tr &a=tree[k];

    if(a.s==s&&a.e==e){

        if(ta<a.num){

            ta=a.num;

            tid=a.id;

        }

        return;

    }

    int mid=(a.s+a.e)>>1;

    if(s>mid)fi(s,e,k<<1|1);

    else if(e<=mid)fi(s,e,k<<1);

    else{

        fi(s,mid,k<<1);

        fi(mid+1,e,k<<1|1);

    }

}

void ins(int tar,int k,int num,int fid){

    tr &a=tree[k];

    if(a.s==a.e){

        if(a.num<num){

            a.num=num;

            a.id=fid;

        }

        return;

    }

    int mid=(a.s+a.e)>>1;

    if(tar<=mid)ins(tar,k<<1,num,fid);

    else ins(tar,k<<1|1,num,fid);

    if(tree[k<<1].num>tree[k<<1|1].num){

        a.num=tree[k<<1].num;

        a.id=tree[k<<1].id;

    }

    else{

        a.num=tree[k<<1|1].num;

        a.id=tree[k<<1|1].id;

    }

}

struct sst{

    int fid,ta,id;

};

int main(){

    int n,w,h;

    stack<sst>ms;

    scanf("%d%d%d",&n,&w,&h);

    jilu[en++]=w;

    for(int i=1;i<=n;i++){

        st &a=all[i];

        scanf("%d%d",&a.w,&a.h);

        a.id=i;

        jilu[en++]=a.w;

    }

    all[n+1].id=n+1;

    all[n+1].h=h;

    all[n+1].w=w;

    sort(all+1,all+2+n,cmp);

    sort(jilu,jilu+en);

    en=unique(jilu,jilu+en)-jilu;

    build(1,en,1);

    bool ok1=0,ok2=0;

    for(int i=1;i<=n+1;i++){

        sst ttmp;

        int id=upper_bound(jilu,jilu+en,all[i].w)-jilu;

        ttmp.fid=id;

        ttmp.id=all[i].id;

        if(i!=1&&all[i].h!=all[i-1].h){

            while(!ms.empty()){

                sst tmp=ms.top();

                ms.pop();

                ins(tmp.fid,1,tmp.ta,tmp.id);

            }

        }

        ta=-1;

        tid=0;

        fi(id+1,en,1);

        ta++;

        ttmp.ta=ta;

        from[all[i].id]=tid;

        if(all[i].id==n+1){

            printf("%d\n",ta);

            int tp=all[i].id;

            while(from[tp]){

                printf("%d ",from[tp]);

                tp=from[tp];

            }

            break;

        }

        ms.push(ttmp);

    }

}

//09:43 WA 16

//09:55 WA 16

//10:02 WA 16

//10:13 WA 16

//10:18 WA 16
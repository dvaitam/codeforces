// LUOGU_RID: 155733018
#include<bits/stdc++.h>
#define ll unsigned long long
#define N 5005
using namespace std;
ll read(){
    ll x=0,f=1;char ch=getchar();
    while(ch<'0' || ch>'9')f=(ch=='-'?-1:f),ch=getchar();
    while(ch>='0' && ch<='9')x=(x<<1)+(x<<3)+(ch^48),ch=getchar();
    return x*f;
}
void write(ll x){
    if(x<0)x=-x,putchar('-');
    if(x/10)write(x/10);
    putchar(x%10+'0');
}
char s[N][N];
ll id[N],has[N],val[N];
ll flag=0;
struct edg{
    ll u,v;
    edg(ll u=0,ll v=0):u(u),v(v){}
};
struct info{
    ll deci,val;
    info(ll deci=0,ll val=0):deci(deci),val(val){}
}p[N],tem[N];
vector<edg> ans;
mt19937_64 myrand(20080623);
bool cmp(info x,info y){
    return x.val<y.val;
}
bool cmp2(ll x,ll y){
    return p[x].val<p[y].val;
}
ll sta[N],n,tothas[N],all;
void solve(ll l,ll r){
    if(r==l)return;
    ll pos=0,tot=0;
    for(ll i=l;i<=r;i++)tot^=has[id[i]];
    for(ll i=l;i<=r;i++){
        ll u=id[i];
        if(val[u]!=tot && val[u]!=0){pos=i;break;}
    }
    if(!pos){
    	//if(n==5000)cout<<l<<' '<<r<<' '<<ans.size()<<endl;
        ll cnt[3];cnt[1]=cnt[2]=cnt[0]=0;
        for(ll i=l;i<=r;i++){
            ll u=id[i];
            if(p[u].deci)cnt[p[u].val]++;
            else if(val[u])cnt[1]++,p[u].val=1;
            else cnt[0]++,p[u].val=0;
        }
        if(!cnt[1] && cnt[0]){flag=1;return;}
        if(cnt[0]<cnt[1]+1 && cnt[1]!=0){flag=1;return;}
        if(!cnt[1]){
            for(ll i=l+1;i<=r;i++)ans.push_back(edg(id[i-1],id[i]));
            return;
        }
        for(ll i=l;i<=r;i++)tem[i]=p[id[i]];
        sort(tem+l,tem+r+1,cmp);ll pos=l,pos1=l+1;
        while(tem[pos].val==0)pos++;ll tmp=pos;
        //cout<<pos1<<' '<<pos<<endl;
        //if(n==5000)cout<<l<<' '<<r<<' '<<ans.size()<<' '<<pos1<<' '<<pos<<endl;
        sort(id+l,id+r+1,cmp2);
        while(pos<=r-cnt[2]){
            ans.push_back(edg(id[pos],id[pos1]));
            ans.push_back(edg(id[pos],id[pos1-1]));
            pos++;pos1++;
        }
        while(pos1<tmp)ans.push_back(edg(id[pos1],id[tmp])),pos1++;
        for(ll i=r-cnt[2]+1;i<=r;i++)ans.push_back(edg(id[tmp],id[i]));
    	//if(n==5000)cout<<l<<' '<<r<<' '<<ans.size()<<endl;
    }
    else{
        ll tp=0,U=id[pos];
        for(ll i=l;i<=r;i++){
            ll u=id[i];
            if(s[u][id[pos]]=='1')sta[++tp]=i;
        }
        for(ll i=1;i<=tp;i++)if(sta[i]!=l+i-1)swap(id[sta[i]],id[l+i-1]);
        for(ll i=l;i<=l+tp-1;i++){
            for(ll j=l+tp;j<=r;j++){
                ll u=id[i],v=id[j];
                if(s[u][v]=='1')val[v]^=has[u];
            }
        }
        for(ll i=l+tp;i<=r;i++){
            for(ll j=l;j<=l+tp-1;j++){
                ll u=id[i],v=id[j];
                if(s[u][v]=='1')val[v]^=has[u];
            }
        }
        ll pos2=0;
        for(ll i=l+tp;i<=r;i++){
        	ll fl=0;
        	//cerr<<(tothas[id[i]]^tothas[U])<<' '<<all<<' '<<p[id[i]].deci<<()<<endl;
        	if(s[id[i]][id[i]]=='0' || p[id[i]].deci || (tothas[id[i]]^tothas[U])!=all)continue;
        	for(ll j=l;j<=l+tp-1;j++){
        		if(s[id[j]][id[i]]!='0'){fl=1;break;}
        	}
        	if(!fl)pos2=i;
        }
        if(!pos2){flag=1;return;}
        ll V=id[pos2];
        ans.push_back(edg(U,V));
        p[V]=info(1,2);p[U]=info(1,2);
        //cerr<<l<<' '<<r<<' '<<U<<' '<<V<<endl;
        for(ll i=l;i<=l+tp-1;i++){
            for(ll j=l+tp;j<=r;j++){
                ll u=id[i],v=id[j];
                if(V==v)continue;
                //cerr<<u<<' '<<v<<' '<<s[u][v]<<' '<<s[V][v]<<endl;
                if(s[u][v]!=s[V][v]){flag=1;return;}
            }
        }
        for(ll i=l+tp;i<=r;i++){
            for(ll j=l;j<=l+tp-1;j++){
                ll u=id[i],v=id[j];
                if(U==v)continue;
                if(s[u][v]!=s[U][v]){flag=1;return;}
            }
        }
        solve(l,l+tp-1);if(flag)return;
        solve(l+tp,r);
    }
}
int main(){
    #ifdef EAST_CLOUD
    freopen("a.in","r",stdin);
    freopen("a.out","w",stdout);
    #endif
    n=read();
    for(ll i=1;i<=n;i++)scanf("%s",s[i]+1);
    if(n==1){printf((s[1][1]=='1'?"NO\n":"YES\n"));return 0;}
    for(ll i=1;i<=n;i++)id[i]=i,has[i]=myrand(),all^=has[i];
    for(ll i=1;i<=n;i++){
        for(ll j=1;j<=n;j++)if(s[i][j]=='1')val[j]^=has[i];
    }
    for(ll i=1;i<=n;i++)tothas[i]=val[i];
    solve(1,n);
    if(flag)printf("NO");
    else{
        printf("YES\n");
        for(auto [u,v]:ans)write(u),putchar(' '),write(v),putchar('\n');
    }
}
#include<bits/stdc++.h>
using namespace std;

#define rep(i,a,b) for(int i=(a);i<=(b);++i)
template<typename T=int>
inline T Read(){
    T r=0;
    int f=0,c=getchar();
    while((c<'0' || c>'9') && ~c) f|=c=='-',c=getchar();
    while(c>='0' && c<='9') r=(r<<1)+(r<<3)+(c^48),c=getchar();
    return f?-r:r;
}

constexpr int AwA=1e6+10;

int n;
char s[AwA];

int tot;
int fail[AwA],trans[AwA][26],len[AwA];
int pos[AwA],id[AwA];

namespace SAM{
    inline void Init(){
        pos[0]=tot=1;
        id[1]=fail[1]=len[1]=0;
    }
    inline void Insert(int cnt,int c){
        int cur=++tot,lst=pos[cnt-1];
        pos[cnt]=cur,len[cur]=cnt,id[cur]=cnt;
        
        while(lst && !trans[lst][c]) trans[lst][c]=cur,lst=fail[lst];
        if(!lst){
            fail[cur]=1;
            return;
        }

        int p=trans[lst][c];
        if(len[p]!=len[lst]+1){
            int q=++tot;
            fail[q]=fail[p],len[q]=len[lst]+1,memcpy(trans[q],trans[p],sizeof(int)*26);
            fail[p]=q;
            while(lst && trans[lst][c]==p) trans[lst][c]=q,lst=fail[lst];
            fail[cur]=q;
        }else fail[cur]=p;
    }
}

int g[AwA],f[AwA];
int lstpos[AwA];
inline void Move(int& u,int length){
    while(u!=1 && len[fail[u]]>=length) u=fail[u];
}

namespace FD{
    inline void Main(){
        n=Read();
        scanf("%s",s+1);
        reverse(s+1,s+n+1);

        SAM::Init();
        rep(i,1,n) SAM::Insert(i,s[i]-'a');
    
        int lst=1;
        rep(i,1,n){
            int c=s[i]-'a';
            int cur=trans[lst][c],curfa=cur;
            assert(cur);
            
            f[i]=f[i-1]+1;
            Move(curfa,f[i]-1);
            while(f[i]!=1){
                if(g[lst]>=f[i]-1 || g[curfa]>=f[i]-1) break;

                f[i]--;
                Move(cur,f[i]);
                Move(lst,f[i]-1);
                Move(curfa,f[i]-1);

                int p=i-f[i];
                int u=lstpos[p];
                g[u]=max(g[u],f[p]);
                u=fail[u];
                while(u && g[u]<len[u]) g[u]=len[u],u=fail[u];
            }
            lst=cur;
            lstpos[i]=cur;
        }

        rep(i,1,n) f[0]=max(f[0],f[i]);
        printf("%d\n",f[0]);
    }
}

int main(){
    FD::Main();
    return 0;
}
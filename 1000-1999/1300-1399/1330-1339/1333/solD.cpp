//by Sshwy
//#define DEBUGGER
#include<algorithm>/*{{{*/
#include<cctype>
#include<cassert>
#include<cmath>
#include<cstdio>
#include<cstring>
#include<cstdlib>
#include<ctime>
#include<iostream>
#include<map>
#include<queue>
#include<set>
#include<vector>
using namespace std;
#define fi first
#define se second
#define pb push_back
#define FOR(i,a,b) for(int i=(a);i<=(b);++i)
#define ROF(i,a,b) for(int i=(a);i>=(b);--i)

#ifdef DEBUGGER

#define log_prefix() { fprintf(stderr, "\033[37mLine %-3d [%dms]:\033[0m  ", __LINE__,(int)clock()/1000);  }
#define ilog(...)    { fprintf(stderr, __VA_ARGS__); }
#define llog(...)    { log_prefix(); ilog(__VA_ARGS__); }
#define log(...)     { log_prefix(); ilog(__VA_ARGS__); fprintf(stderr,"\n");  }
#define red(...)     { log_prefix(); ilog("\033[31m"); ilog(__VA_ARGS__); ilog("\033[0m\n"); }
#define green(...)   { log_prefix(); ilog("\033[32m"); ilog(__VA_ARGS__); ilog("\033[0m\n"); }
#define blue(...)    { log_prefix(); ilog("\033[34m"); ilog(__VA_ARGS__); ilog("\033[0m\n"); }

#else

#define log_prefix() ;
#define ilog(...)    ;
#define llog(...)    ;
#define log(...)     ;
#define red(...)     ;
#define green(...)   ;
#define blue(...)    ;

#endif

namespace RA{
    int r(int p){return 1ll*rand()*rand()%p;}
    int r(int L,int R){return r(R-L+1)+L;}
}/*}}}*/
namespace IO{//require C++11
    const int _BS=200000;
    char _ob[_BS],*_z=_ob,_ib[_BS],*_x=_ib,*_y=_ib;
    // Input
    char _nc(){
        return _x==_y&&(_y=(_x=_ib)+fread(_ib,1,_BS,stdin),_x==_y)?EOF:*_x++;
    }
    inline void rd(char & c){ c=_nc(); }
    void rd(char * s){
        char c=_nc();
        while(!isblank(c)&&isprint(c)&&c!=EOF) *s++=c,c=_nc();
        *s=0;
    }
    template<class T>
    void rd(T & res){// unsigned/signed int/long long
        res=0; char c=_nc(),flag=0;
        while(!isdigit(c) && c!='-')c=_nc();
        if(c=='-')flag=1,c=_nc();
        while(isdigit(c))res=res*10+c-'0',c=_nc();
        if(flag)res=-res;
    }
    template <class T, class ...A> void rd(T & hd, A & ...rs) { rd(hd), rd(rs...); }
    // Output
    inline void flush(){ fwrite(_ob,1,_z-_ob,stdout),_z=_ob; }
    inline void wrch(char x){
        if(_z==_ob+_BS-1)flush();
        *_z=x,++_z;
    }
    void wr(char * s){ while(*s)wrch(*s),++s; }
    void wr(const char * s){ while(*s)wrch(*s),++s; }
    template<class T>
        void wr(T x){ // unsigned/signed int/long long
            if(x==0)return wrch('0'),void();
            if(x<0)wrch('-'),x=-x;
            int _t[30],_lt=0;
            while(x)_t[++_lt]=x%10,x/=10;
            while(_lt>0)wrch(_t[_lt]+'0'),--_lt;
        }
    inline void wr(char x){ wrch(x); }
    template <class T, class ...A> void wr(T hd, A ... rs) { wr(hd), wr(rs...); }
    template <class ...A> void wrln(A ... ar) { wr(ar...), wrch('\n'); }
    struct _Flusher{~_Flusher(){flush();}}_flusher;
}
/******************heading******************/
const int N=3005;
int n,k,tot;
char s[N];
vector< vector<int> > V,ans;
int main(){
    IO::rd(n,k);
    IO::rd(s+1);
    //cin>>n>>k;
    //cin>>(s+1);
    while(1){
        int cur=0,mn=n+1;
        vector<int> v;
        FOR(i,1,n){
            if(s[i]=='L')++cur;
            else --cur;
            if(s[i]=='L' && i>1 && s[i-1]=='R'){
                if(cur<mn){
                    v.clear();
                    mn=cur;
                    v.pb(i-1);
                }else if(cur==mn){
                    v.pb(i-1);
                }
            }
        }
        if(v.size()){
            //cout<<"> "; for(auto x:v)cout<<x<<" "; cout<<endl;
            for(int x:v)s[x]='L',s[x+1]='R';
            //cout<<(s+1)<<endl;
            V.pb(v);
            tot+=v.size();
        }else break;
    }
    //log("GG");
    if(V.size()>k)return IO::wrln(-1),0;//cout<<-1<<endl,0;
    if(tot<k)return IO::wrln(-1),0;//cout<<-1<<endl,0;
    int y=V.size();
    for(auto & v:V){
        --y;
        if(v.size()+y<=k){
            for(auto x:v){
                IO::wrln("1 ",x);
                //cout<<"1 "<<x<<endl;
            }
            y+=v.size();
        }else {
            int cnt=k-y-1;
            assert(cnt<v.size());
            FOR(i,0,cnt-1){
                IO::wrln("1 ",v[i]);
                //cout<<"1 "<<v[i]<<endl;
                ++y;
            }
            IO::wr(int(v.size()-cnt)," ");
            //cout<<v.size()-cnt<<" ";
            for(int i=cnt;i<v.size();++i)IO::wr(v[i]," ");//cout<<v[i]<<" ";
            IO::wr('\n');
            //for(auto x:v)cout<<x<<" ";
            //cout<<endl;
            ++y;
        }
    }
    return 0;
}
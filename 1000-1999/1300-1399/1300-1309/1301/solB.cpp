#include<cstdio>
const int N=1e5+5;
int T,n,Maxx,Minn,Num[N],Answer;
inline int Min(const int x,const int y){
    return x<y?x:y;
}
inline int Max(const int x,const int y){
    return x>y?x:y;
}
inline int Read(){
    char ch;
    int f=1;
    while((ch=getchar())<'0'||ch>'9')
        if(ch=='-') f=-1;
    int x=ch^48;
    while((ch=getchar())>='0'&&ch<='9')
        x=(x<<3)+(x<<1)+(ch^48);
    return x*f;
}
inline void Init(){
    Minn=1e9+1;
    Maxx=-1;
    Answer=0;
    n=Read();
    for(register int i=1;i<=n;i++){
        Num[i]=Read();
        if(Num[i]==-1&&i>1&&Num[i-1]!=-1){
            Maxx=Max(Maxx,Num[i-1]);
            Minn=Min(Minn,Num[i-1]);
        }
        if(Num[i-1]==-1&&Num[i]!=-1){
            Maxx=Max(Maxx,Num[i]);
            Minn=Min(Minn,Num[i]);
        }
    }
    return ;
}
inline int Abs(const int x){
    return x>0?x:-x;
}
inline void Query(){
    int k=Maxx+Minn>>1;
    for(register int i=1;i<=n;i++)
        if(Num[i]==-1) Num[i]=k;
    for(register int i=1;i<n;i++)
        Answer=Max(Answer,Abs(Num[i+1]-Num[i]));
    printf("%d %d\n",Answer,k);
    return ;
}
int main(){
    for(T=Read();T;T--){
        Init();
        Query();
    }
    return 0;
}
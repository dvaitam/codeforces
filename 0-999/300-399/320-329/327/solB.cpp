/* Miaki-Miha Solved On 4th July,2013 */
#include<cstdio>
#define Wsput putchar(' ')
#define Lnput putchar('\n')

inline void Iget(int &x){
    static char ch; x=0;
    while (((ch=getchar())<'0')||(ch>'9'));
    do{x=x*10+ch-'0';}while (((ch=getchar())>='0')&&(ch<='9'));
}

inline void Iput(int x){
    static char s[20],*p; p=s;
    do{*(++p)=x%10+'0',x/=10;}while (x);
    while (p!=s) putchar(*(p--));
}

int n;
int main(){
    Iget(n);
    for(int i=n;i;--i) Iput(10000000-i+1),i^1?Wsput:Lnput;
    return 0;
}
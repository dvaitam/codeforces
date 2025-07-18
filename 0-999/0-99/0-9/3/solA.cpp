#include <stdio.h>
char d[8][3]={"L", "R", "U", "D", "LU", "LD", "RU", "RD"};
int dx[]={-1,1,0,0,-1,-1,1,1};
int dy[]={0,0,1,-1,1,-1,1,-1};

int max(int a, int b){
    if(a>b) return a;
    return b;
}

int abs(int a){
    if(a>0) return a;
    return -a;
}

int main(){
    char s[10],t[10];
    scanf("%s%s",s,t);
    s[0]-='a';
    s[1]-='0';

    t[0]-='a';
    t[1]-='0';
    int ans=max(abs(s[0]-t[0]),abs(s[1]-t[1]));
    printf("%d\n", ans);
    while(s[0]!=t[0] || s[1]!=t[1]){
        for(int i=0;i<8;i++){
            s[0]+=dx[i];
            s[1]+=dy[i];
            if(max(abs(s[0]-t[0]),abs(s[1]-t[1])) == ans-1){
                printf("%s\n",d[i]);
                ans--;
                break;
            }else{
                s[0]-=dx[i];
                s[1]-=dy[i];
            }
        }
    }
}
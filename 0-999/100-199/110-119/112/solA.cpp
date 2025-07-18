#include<stdio.h>
#include<string.h>
#define FOR(_i , N) for(int _i = 0 ;_i < N ; _i ++)
char a[1000],b[1000];
int main(){
    scanf("%s",&a);
    scanf("%s",&b);

    int n=strlen(a);
    FOR(i,n){
        if(a[i]>=97) a[i]=a[i]-32;
        if(b[i]>=97) b[i]=b[i]-32;
    }
    FOR(i,n){
        if(a[i]>b[i]){
            printf("1");
            return 0;
        }
        else if(a[i]<b[i]){
            printf("-1");
            return 0;
        }
    }
    printf("0");
    return 0;
}
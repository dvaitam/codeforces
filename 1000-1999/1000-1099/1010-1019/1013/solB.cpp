#include<iostream>
#include<cstdio>
#include<string.h>
#include<limits.h>
using namespace std;
int main(){
    bool input[100005],dan[100005];
    int n,x,y;
    cin>>n>>x;
    memset(input,0,sizeof(input));
    memset(dan,0,sizeof(dan));
    bool p=false;
    int cnt=INT_MAX;
    for(int a=0;a<n;a++){
        scanf("%d",&y);
        if(input[y]){
            p=true;
            cnt=min(cnt,0);
        }
        if(dan[y]){
            p=true;
            cnt=min(cnt,1);
        }
        if(input[y & x]){
            p=true;
            cnt=min(cnt,1);
        }
        if(dan[y & x]){
            cnt=min(cnt,2);
            p=true;
            
        }
        input[y]=1;
        dan[y & x]=1;
    }
    if(p)cout<<cnt<<endl;
    else cout<<-1<<endl;
}
#include<bits/stdc++.h>
using namespace std;
const int maxn=1e5+10;
char a[maxn];
int main(){
    scanf("%s",a);
    int len=strlen(a);
    stack<char>q1,q2;
    for(int i=0;i<len;i++){
        q1.push(a[i]);
    }
    if(len==1){
        printf("NO\n");
        return 0;
    }
    int ans=0,cnt=0;
    while(!q1.empty()){
        char now1=q1.top();
        q1.pop();
        if(q2.empty()){
            q2.push(now1);
            if(q1.empty()) break;
            now1=q1.top();
            q1.pop();
        }
        char now2=q2.top();
        q2.pop();
        if(now2!=now1){
            q2.push(now2);
            q2.push(now1);
        }
        else ans++;
    }
    //printf("ans %d\n",ans);
    if(ans%2==0){
        printf("NO\n");
    }
    else printf("YES\n");
    return 0;
}
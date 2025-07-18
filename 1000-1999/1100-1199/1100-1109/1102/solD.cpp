#include<cstdio>
#include<vector>
#include<algorithm>
using namespace std;

int n;
char s[300001];

void solve(){
    int a=0, b=0, c=0;
    for(int i=0;i<n;++i){
        if(s[i]=='0') a++;
        else if(s[i]=='1') b++;
        else c++;
    }

    for(int i=0;i<n && a<n/3;++i){
        if(s[i] == '1' && b>n/3)
            s[i] = '0', a++, b--;
        else if(s[i] == '2' && c>n/3)
            s[i] = '0', a++, c--;
    }
    for(int i=n-1;i>=0 && a>n/3;--i)
        if(s[i] == '0') s[i] = '3', a--;

    for(int i=0;i<n && b<n/3;++i)
        if(s[i] =='3') s[i] = '1', b++;
    for(int i=0;i<n && b<n/3;++i)
        if(s[i] =='2') s[i] = '1', b++, c--;
    for(int i=n-1;i>=0 && b>n/3;--i)
        if(s[i] == '1') s[i] = '3', b--;

    for(int i=0;i<n;++i)
        if(s[i]=='3') s[i] = '2';

}

int main(){
    scanf("%d\n%s",&n,s);
    solve();
    printf("%s",s);
    return 0;
}
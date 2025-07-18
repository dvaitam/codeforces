#include<bits/stdc++.h>
using namespace std;
typedef long long int ll;
typedef unsigned long long int Ull;
typedef pair<int ,int> PII;
template <typename T>
inline T  read()
{
    char c=getchar();
    T x=0,f=1;
    while(!isdigit(c)){if(c=='-')f=-1;c=getchar();}
    while(isdigit(c))x=(x<<3)+(x<<1)+(c^48),c=getchar();
    return x*f;
}
template <typename T>
inline void write(T x) 
{
    if(x < 0) x = (~x) + 1, putchar('-');
    if(x / 10) write(x / 10);
    putchar(x % 10 | 48);
}
const int INF = 0x3f3f3f3f;
const ll N = 3e5+10;
int arr[N];
int main(){
    //ios_base::sync_with_stdio(0),cin.tie(0),cout.tie(0); 
    int  t = read<int>();
    while(t--){
        int n = read<int>();
        int mx = -INF, ii;
        for(int i = 1; i <= n; i++){
            arr[i] = read<int>();
            if(arr[i] > mx){
                mx = arr[i];
                ii = i;
            }
        }
        bool flag = 0;
        arr[0] = INF;
        arr[n+1] = INF;
        for(int i = 1; i <= n; i++){
            if(arr[i] == mx && (arr[i+1] < mx || arr[i-1] < mx)){
                flag = 1;
                ii = i;
                break;
            }
        }
        if(!flag) printf("-1\n");
        else printf("%d\n",ii);
    }
}
/*
 ID: tiendao1
 LANG: C++
 TASK: cowtour
 */
 /*
    
 */
 
#include<iostream>
#include<cstdio>
#include <sstream>
#include<string>
#include<cstring>
#include<map>
#include<set>
#include<deque>
#include<queue>
#include<stack>
#include<vector>
#include<algorithm>
#include<cstdlib>
#include<bitset>
#include<climits>
#include <list>
#include <functional>
//#include <cmath>

using namespace std;
typedef long long ll;
typedef unsigned long long ull;
#define INF 1e9
#define EPS 1e-7
#define PI acos(-1)
#define pii pair<int,int>
#define pll pair<unsigned long long, unsigned long long>
#define F first
#define S second
#define ROW first
#define COL second

long long gcd(long long x, long long y) {
    return y == 0 ? x : gcd(y, x % y);
}
int toInt(char xx) {
    return xx - '0';
}
char toChar(int xx) {
    return xx + '0';
}
int isDigit(char xx) {
    return ('0' <= xx && xx <= '9');
}
bool isCharacter(char xx) {
    return (('a' <= xx && xx <= 'z') || ('A' <= xx && xx <= 'Z'));
}
void swapInt(int &x, int &y) {
    x = x ^ y;
    y = x ^ y;
    x = x ^ y;
}

/*** IMPLEMENTATION ***/
const int dx[4] = { -1, 1, 0, 0 };
const int dy[4] = { 0, 0, -1, 1 };
const int Mod = 9901;
const int maxn = (1e6) + 10 ;
int exitInput = 0;

int n , k ; 
int ans[maxn] ; 

void read() 
{
#define endread { exitInput = 1 ; return ; }   
    cin >> n >> k ; 
}

void init() 
{   
    int i ; 
    for(i = 0; i < maxn; ++i) ans[i] = -1 ; 
}   

void solve() 
{
    if(n <= 3 || 3*k > n)
    {
        cout << -1 << endl ; 
        return ;
    }
    int id, i ; 
    
    i = 1 ; 
    
    int heso = 4 , lastval, min_not_ans ; 
    for(id = 1; id <= k-1 ; ++id)
    {
    //  cout << " i = " << i << endl ;
    //  cout << " heso = " << heso << endl ;
        lastval = i + 3 ; 
        if(heso == 4)
        {
            min_not_ans = i+2 ;
            ans[i] = id ; 
            ans[i+1] = id ; 
            ans[i+3] = id ; 
            heso = 2;
            i += 2 ;
        }
        else if(heso == 2)
        {
            min_not_ans = i+4 ;
            ans[i] = id ; 
            ans[i+2] = id ; 
            ans[i+3] = id ;
            heso = 4 ; 
            i += 4 ; 
        }
    }
    
    //cout << " lastval = " << lastval << endl;
    //cout << " min_not_ans =  " << min_not_ans << endl ;
    
    if(lastval < min_not_ans)
    {
        if(min_not_ans+2 == n)
        {
            // there are 3 number left : min_not_ans, min_not_ans+1,min_not_ans+2
            ans[min_not_ans] = k ; 
            ans[min_not_ans+2] = k ; 
            ans[min_not_ans+1] = k-1 ; 
            ans[min_not_ans-1] = k ;
            
        }
        else
        {
            for( i = min_not_ans; i <= n; ++i)
            {
                ans[i] = k ;
            }
            ans[min_not_ans+1] = 1 ; 
        }
        
    }
    else
    {
        for( i = min_not_ans; i <= n; ++i)
        {
            if(ans[i] == -1) ans[i] = k ;
        }
    }
    
    cout << ans[1] ;
    for(i = 2; i <= n; ++i)
    {
        cout <<  " " << ans[i] ;
    }
    cout << endl;
}

int main() {
#ifndef ONLINE_JUDGE
    freopen("input.txt", "r", stdin); freopen("output.txt", "w", stdout);

  // freopen("cowtour.in", "r", stdin);freopen("cowtour.out", "w", stdout);
#endif
    ios_base::sync_with_stdio(0);
    
    int ntest = 1;
    int itest = 1;
    
    
    //cin >> ntest;
        
    for (itest = 1; itest <= ntest ; ++itest) {
        read();
        if (exitInput) {
            break;
        }   
       // if(itest > 1)
        {
           // cout << endl;
            //printf("\n") ;
        }
        //cout << "Case " << itest << ": " ;
        init();
        solve();
        
    }
    return 0;
}
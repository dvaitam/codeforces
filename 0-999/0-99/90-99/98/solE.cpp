#include<cstdio>
#include<cstring>
using namespace std;

const int N = 1005;

#define fo(i , st , en) for (int i = st; i <= en; i++)
#define Me(x , y) memset(x , y , sizeof(x))

double f[N][N];
bool flag[N][N];
int n , m;

double dfs(int x , int y){
    if (!x || !y) return 1.0 / (y + 1);
    if (flag[x][y]) return f[x][y];
    flag[x][y] = 1;
    double a = 1 - dfs(y , x - 1) , b = (double)y / (y + 1) * (1 - dfs(y - 1 , x)) , c = b + 1.0 / (y + 1) , p = (c - b) / (1 - a - b + c);
    f[x][y] = p * (a - c) + c;
    return f[x][y];
}

int main(){
    scanf("%d%d" , &n , &m);
    Me(flag , 0); printf("%.10lf %.10lf\n" , dfs(n , m) , 1 - dfs(n , m));
    return 0;
}
# include <iostream>
# include <vector>
# include <queue>
# include <deque>
# include <stack>
# include <set>
# include <string>
# include <algorithm>
# include <cstdio>
# include <cstdlib>
# include <map>
# include <cmath>
# include <bitset>
# include <list>
# include <sstream>

using namespace std;

int main(){
 //   freopen("input.txt", "r", stdin);
   //  freopen("output.txt", "w", stdout);
    
    int N, M; cin >> N >> M;
    vector<int> price(N);
    for(int i=0; i<price.size(); ++i) scanf("%d", &price[i]);
    
    vector<vector<bool>> con(N, vector<bool>(N, false));
    for(int i=0; i<M; ++i)
    {
        static int a,b;
        scanf("%d %d", &a, &b);
        con[a-1][b-1] = con[b-1][a-1] = true;
    }
    
    float Max = numeric_limits<float>::infinity();
    
    for(int i=0; i<N; ++i)
    {
        for(int k=i+1; k < N; ++k)
        {
            for(int j=k+1; j<N; ++j)
            {
                if (con[i][k] && con[i][j] && con[k][j]) Max = min(Max, (float)(price[i]+price[j]+price[k]));
            }
        }
    }
    
    cout << (Max == numeric_limits<float>::infinity() ? -1:(int)Max);
    
    return 0;
}
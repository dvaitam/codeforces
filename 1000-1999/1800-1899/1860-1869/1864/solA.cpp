#include <iostream>
#include <string>
#include <vector>
#include <cmath>
#include <set>
#include <utility>
#include <algorithm>
#include <map>


void solve()
{
    int x, y, n;
    std::cin >> x >> y >> n;
    std::vector<int> b(n);
    b[0] = x;
    b[n-1] = y;
    std::reverse(b.begin(), b.end());
    b[1] = b[0] - 1;
    
    for (int i = 2; i < n - 1; ++i)
    {
        b[i] = b[i - 1] - i;
    }
    std::reverse(b.begin(), b.end());
    int last = b[1] - b[0];
    for (int i = 2; i < n; ++i)
        if ((last <= b[i] - b[i - 1]) || (b[i] == b[i -1]))
        {
            std::cout << -1 << '\n';
            return;
        }
    
    for (int i = 0; i < n; ++i)
        std::cout << b[i] << ' ';
    std::cout << '\n';
    
}



int main()
{   
    std::ios::sync_with_stdio(false);
    std::cin.tie(nullptr);
    int tt; std::cin >> tt;
    
    while (tt--)
    {   
        solve();
    }
    return 0;
}
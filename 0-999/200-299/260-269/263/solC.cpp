#include <iostream>
#include <cstdlib>
#include <cstdio>
#include <algorithm>
#include <cstring>
#include <cmath>
#include <fstream>
#include <algorithm>
#define rep(i,j,k) for (int i=j;i<=k;++i)
#define rrep(i,j,k) for (int i=j;i>=k;--i)

using namespace std;

int n;
int f[111111][10];
int a[111111],num[5];
bool use[111111],found;

bool getRem()
{
    rep(i,3,n-1)
	{
	    int u,v;
	    found = false;
	    rep(p,1,4)
		{
		    u = f[a[i-1]][p];
		    if (use[u]) continue;
		    rep(q,1,4) if (f[a[i-2]][q] == u)
			{
			    use[u] = true;
			    a[i] = u;
			    found = true;
			}
		    if (found) break;
		}
	    if (!found)
		{
		    return false;
		}
	}
    return true;
}

bool check()
{
    bool t1,t2,t3,t4,t5;
    t1 = t2 = t3 = t4 = t5 = false;
    rep(i,1,4) 
	{
	    if (f[a[n]][i] == a[1]) t1 = true;
	    if (f[a[n]][i] == a[2]) t2 = true;
	    if (f[a[n]][i] == a[n-1]) t3 = true;
	    if (f[a[n]][i] == a[n-2]) t4 = true;
	}
    if (t1 == true && t2 == true && t3 == true && t4 == true) 
	t5 = true;
    return t5;
}

void print()
{
    rep(i,1,n-1) cout << a[i] << " ";
    cout << a[n] << endl;
    exit(0);
}

int main()
{
    /*
    freopen("test.in","r",stdin);
    freopen("test.out","w",stdout);
    */
    ios::sync_with_stdio(false);
    cin >> n;
    rep(i,1,2*n)
	{
	    int u,v;
	    cin >> u >> v;
	    f[u][++f[u][0]] = v;
	    f[v][++f[v][0]] = u;
	}

    a[1] = 1;found = false;
    rep(p,1,4) 
	{
	    rep(q,1,4) if (p!=q)
		{
		    int u = f[1][p];
		    int v = f[1][q];
		    rep(r,1,4) if (f[u][r] == v)
			{
			    memset(use,false,sizeof(use));
			    use[1] = true;
			    a[n] = a[0] = u,a[2] = v;
			    use[u] = use[v] = true;
			    if (getRem() && check()) 
				{
				    print();
				    return 0;
				}
			}
		}
	}
    cout << -1 << endl;
    return 0;
}
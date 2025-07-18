#include <cstdio>
#include <cstdlib>
#include <cstring>

#include <algorithm>
#include <iostream>
#include <fstream>
#include <string>
#include <vector>
#include <map>
#include <stack>
#include <set>
#include <queue>

using namespace std;

#define MAX     150005
#define URL_HEAD 7

typedef long long ll;

int n, m, r, p;
char name[15];

struct student {
    string name;
    int points;
};

student top[10005][3];

int main()
{
    scanf("%d%d", &n, &m);
    
    for (int i=1; i<=10000; i++) {
        top[i][0].points = -1;
        top[i][1].points = -1;
    }
    
    for (int i=0; i<n; i++) {
        scanf("%s%d%d", name, &r, &p);
        
        if (p >= top[r][0].points) {
            top[r][2] = top[r][1];
            top[r][1] = top[r][0];
            
            top[r][0].points = p;
            top[r][0].name = name;
        }
        else if (p >= top[r][1].points) {
            top[r][2] = top[r][1];
            
            top[r][1].points = p;
            top[r][1].name = name;
        }
        else if (p >= top[r][2].points) {
            top[r][2].points = p;
            top[r][2].name = name;
        }
    }
    
    for (int i=1; i<=m; i++) {
        if (top[i][1].points == top[i][2].points) {
            printf("?\n");
        }
        else {
            printf("%s %s\n", top[i][0].name.c_str(), top[i][1].name.c_str());
        }
    }
    
    return 0;
}
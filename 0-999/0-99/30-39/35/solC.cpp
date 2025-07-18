#include<stdio.h>
#include<queue>
#include<map>
using namespace std;
queue<pair<int , int > > cola;
int p , a,b,n,m;
bool C[3000][3000];
int main()
{
FILE * fi = fopen("input.txt","r");
FILE * fo = fopen("output.txt","w");

 fscanf(fi,"%d%d",&n,&m);
 fscanf(fi,"%d",&p);
  while(p--)
  {
   fscanf(fi,"%d%d",&a,&b);
   cola.push(make_pair(a,b));
   C[a][b] = 1;
  }
   pair<int ,int> l = make_pair(1,1);
   while(cola.size())
   {
     l = cola.front();cola.pop();
     if(l.first != 1 && !C[l.first-1][l.second] )
     {
        C[l.first-1][l.second] = 1;
        cola.push(make_pair(l.first-1,l.second));
     }
     if(l.first != n && !C[l.first+1][l.second] )
     {
        C[l.first+1][l.second] = 1;
        cola.push(make_pair(l.first+1,l.second));
     }

     if(l.second != 1 && !C[l.first][l.second-1] )
     {
        C[l.first][l.second-1] = 1;
        cola.push(make_pair(l.first,l.second-1));
     }
     if(l.second != m && !C[l.first][l.second+1] )
     {
        C[l.first][l.second+1] = 1;
        cola.push(make_pair(l.first,l.second+1));
     }
   }
   fprintf(fo,"%d %d\n",l.first,l.second);
 return 0;
}
#include <bits/stdc++.h>
using namespace std ;

typedef long long LL ;

#define clr( a , x ) memset ( a , x , sizeof a )

const int MAXN = 755 ;
const int MAXS = 1e7 + 5 ;
const int INF = 0x3f3f3f3f ;

struct Node {
 Node* nxt[2] ;
 Node* fail ;
 int end ;
} node[MAXS] , *cur , *root , *q[MAXS] ;

int Q[MAXN] , head , tail ;
int Lx[MAXN] , Ly[MAXN] ;
int dx[MAXN] , dy[MAXN] ;
int vis[MAXN] , Time ;
int mdis ;

bitset < MAXN > G[MAXN] ;
int usex[MAXN] , usey[MAXN] ;
char s[MAXS] ;
int st[MAXN] , len[MAXN] ;
int n ;

Node* newnode () {
 cur->nxt[0] = cur->nxt[1] = cur->fail = NULL ;
 cur->end = -1 ;
 return cur ++ ;
}

void init () {
 for ( int i = 0 ; i < n ; ++ i ) {
  G[i].reset () ;
 }
 cur = node ;
 root = newnode () ;
}

void insert ( char s[] , int n , int idx ) {
 Node* now = root ;
 for ( int i = 0 ; i < n ; ++ i ) {
  int x = s[i] - 'a' ;
  if ( now->nxt[x] == NULL ) now->nxt[x] = newnode () ;
  now = now->nxt[x] ;
 }
 if ( ~now->end ) G[idx][now->end] = G[now->end][idx] = 1 ;
 now->end = idx ;
}

void build () {
 head = tail = 0 ;
 root->fail = root ;
 for ( int i = 0 ; i < 2 ; ++ i ) {
  if ( root->nxt[i] == NULL ) {
   root->nxt[i] = root ;
  } else {
   root->nxt[i]->fail = root ;
   q[tail ++] = root->nxt[i] ;
  }
 }
 while ( head != tail ) {
  Node* now = q[head ++] ;
  if ( now->end == -1 ) now->end = now->fail->end ;
  for ( int i = 0 ; i < 2 ; ++ i ) {
   if ( now->nxt[i] == NULL ) {
    now->nxt[i] = now->fail->nxt[i] ;
   } else {
    now->nxt[i]->fail = now->fail->nxt[i] ;
    q[tail ++] = now->nxt[i] ;
   }
  }
 }
}

void query ( char s[] , int n , int idx ) {
 Node* now = root ;
 for ( int i = 0 ; i < n ; ++ i ) {
  now = now->nxt[s[i] - 'a'] ;
  int x = now->fail->end ;
  int y = now->end ;
  if ( i == n - 1 && x != -1 ) G[x][idx] = 1 ;
  if ( i <  n - 1 && y != -1 ) G[y][idx] = 1 ;
 }
}

void floyd () {
 for ( int k = 0 ; k < n ; ++ k ) {
  for ( int i = 0 ; i < n ; ++ i ) if ( G[i][k] ) {
   G[i] |= G[k] ;
  }
 }
}

bool bfs () {
 mdis = INF ;
 clr ( dx , -1 ) ;
 clr ( dy , -1 ) ;
 head = tail = 0 ;
 for ( int i = 0 ; i < n ; ++ i ) {
  if ( Lx[i] == -1 ) {
   dx[i] = 0 ;
   Q[tail ++] = i ;
  }
 }
 while ( head != tail ) {
  int u = Q[head ++] ;
  if ( dx[u] >= mdis ) continue ;
  for ( int v = 0 ; v < n ; ++ v ) if ( G[u][v] ) {
   if ( dy[v] == -1 ) {
    dy[v] = dx[u] + 1 ;
    if ( ~Ly[v] ) {
     dx[Ly[v]] = dy[v] + 1 ;
     Q[tail ++] = Ly[v] ;
    } else mdis = dy[v] ;
   }
  }
 }
 return mdis != INF ;
}

bool dfs ( int u ) {
 for ( int v = 0 ; v < n ; ++ v ) if ( G[u][v] ) {
  if ( vis[v] != Time && dy[v] == dx[u] + 1 ) {
   vis[v] = Time ;
   if ( ~Ly[v] && dy[v] == mdis ) continue ;
   if ( Ly[v] == -1 || dfs ( Ly[v] ) ) {
    Ly[v] = u ;
    Lx[u] = v ;
    return 1 ;
   }
  }
 }
 return 0 ;
}

int HopCroft_Krap () {
 int ans = 0 ;
 clr ( Lx , -1 ) ;
 clr ( Ly , -1 ) ;
 clr ( vis , 0 ) ;
 Time = 0 ;
 while ( bfs () ) {
  ++ Time ;
  for ( int i = 0 ; i < n ; ++ i ) {
   if ( Lx[i] == -1 ) ans += dfs ( i ) ;
  }
 }
 return ans ;
}

void solve () {
 init () ;
 st[0] = 0 ;
 for ( int i = 0 ; i < n ; ++ i ) {
  scanf ( "%s" , s + st[i] ) ;
  len[i] = strlen ( s + st[i] ) ;
  insert ( s + st[i] , len[i] , i ) ;
  st[i + 1] = st[i] + len[i] ;
 }
 build () ;
 for ( int i = 0 ; i < n ; ++ i ) {
  query ( s + st[i] , len[i] , i ) ;
 }
 floyd () ;
 HopCroft_Krap () ;
 clr ( usex , 0 ) ;
 clr ( usey , 0 ) ;
 clr ( vis , 0 ) ;
 for ( int i = 0 ; i < n ; ++ i ) {
  if ( ~Lx[i] ) usex[i] = 1 ;
 }
 while ( 1 ) {
  bool flag = 0 ;
  for ( int i = 0 ; i < n ; ++ i ) if ( !usex[i] && !vis[i] ) {
   flag = 1 ;
   vis[i] = 1 ;
   for ( int j = 0 ; j < n ; ++ j ) if ( G[i][j] ) {
    if ( usey[j] == 0 ) {
     usey[j] = 1 ;
     usex[Ly[j]] = 0 ;
    }
   }
  }
  if ( !flag ) break ;
 }
 int ans = 0 , flag = 0 ;
 for ( int i = 0 ; i < n ; ++ i ) {
  if ( !usex[i] && !usey[i] ) ++ ans ;
 }
 printf ( "%d\n" , ans ) ;
 for ( int i = 0 ; i < n ; ++ i ) {
  if ( !usex[i] && !usey[i] ) {
   if ( flag ) printf ( " " ) ;
   flag = 1 ;
   printf ( "%d" , i + 1 ) ;
  }
 }
 puts ( "" ) ;
}

int main () {
 while ( ~scanf ( "%d" , &n ) ) solve () ;
 return 0 ;
}
/*
Submission by Kohlic
This is a solution to Codeforces Round 880 Division 1 Problem F.
https://codeforces.com/contest/1835/problem/F

This is a graph theory problem with a beautiful mathematical proof.
This problem is extremely difficult, with an ELO rating of 3500.
Only 6 people in the world have ELO rating greater than 3500 on Codeforces.
(Codeforces has 153978 users in total.)
During the contest, only 3 participants solved this problem
(out of 12932 participants in Round 880 Divisions 1 and 2).


Problem Statement:
Notation: For any set S, we denote |S| to be the cardinality (size) of S.

You are given a bipartite graph G with 2*n vertices.
The bipartite graph has n vertices on the left side, and n vertices on the
right side.
Let's call the vertex set on the left part as L.
Let's call the vertex set on the right part as R.
We have m edges connecting these two sets.
We know that |L| = |R| = n.

Definition of N(S):
For any subset S of L, let N(S) denote the set of all neighbors of vertices in
S.

Definition of a tight subset:
We say that a subset S of L in graph G is tight if |S| = |N(S)|.

Definition of a good graph:
We say that graph G is good if for all subsets S of L, |S| is less or equal to
|N(S)|.

Your task is to verify whether the graph is good and, if so, to optimize it.
If the graph is not good, find a subset S of L such that |S| > |N(S)|.
However, if the graph is good, your task is to find a good bipartite graph H
such that:
(1) H has the same set of vertices as G (L and R), and,
(2) For all subsets S of L, S is tight in G if and only if S is tight in H.
If there are multiple such graphs H, choose one with the smallest possible
number of edges.
If there are still multiple such graphs, print any.


Constraints:
1 <= n <= 1000
0 <= m <= n*n


Input format:

The first line of the input contains two integers n and m separated by a single
 space.
The number n denotes the number of vertices in each of the sets L and R, and
 the number m
denotes the number of edges between them.
The following m lines describe the edges.
Each of them contains two integers l and r (1 <= l <= n, n+1 <= r <= 2*n),
 separated by a
single space, indicating that there is an edge from vertex l (in L) to vertex
 r (in R).


Output format:

If the graph G given in the input is not good, output one word "NO" in the
 first line of the output.
In the second line of the output, output the number k, and in the third line,
 output k numbers l1,l2,...,lk
, separated by single spaces. These numbers should indicate that for the set
 S={l1,l2,...,lk}, |S|>|N(S)|.

However, if the graph G given in the input is good, output one word "YES" in
 the first line of the output.
In the second line of the output, output the number e, indicating the number
 of edges in the new, also good graph H.
Then, in the following e lines, output the edges of the graph H in the same
 format as given in the input.



Graph theory definitions:
Definition of matching:
A matching is a subset of edges such that any 2 edges in this subset do not
 share a common vertex.

Definition of maximum matching:
A maximum matching is a matching containing the maximum number of edges.

Definition of perfect matching:
A perfect matching is a matching that covers all vertices.



Solution:
We split the solution into 3 phases.



Phase (A): Determining whether a graph G is good.
Using Hall's Marriage Theorem for Graphs, we know that:
A finite bipartite graph G has a perfect matching if and only if all subsets
 S of L is tight in G.
We use Edmond-Karp max-flow algorithm to find a perfect matching.
Graph G is good if and only if there exists a perfect matching.
If no perfect matching exists, proceed to Phase (B) below.
If a perfect matching exists, proceed to Phase (C) below.



Phase (B): Suppose no perfect matching exists:
Then, we know the input graph G is not good.
Our goal: Find a subset W of L such that |W| > |N(W)|, and print out W.

We can do this procedure by following the proof of Hall's Marriage Theorem:
Step (1): Find a maximum matching M.
Step (2): Choose any vertex u in L that does not belong to the matching M.
Step (3): Define an alternating path to be a path that alternately uses edges
 outside and inside M.
Step (4): Define P to be all alternating paths starting from vertex u, using
 breadth-first-search.
Step (5): Define W to be the set of vertices in the alternating paths P that
 belong to L.
I claim that |W| > |N(W)|. Thus, we can output W as our answer.
Let's prove that |W| > |N(W)|.

Define Z be the set of vertices in the alternating paths P that belong to R.

Lemma 1: I claim that every vertex Z is matched by M to a vertex in W.
Proof of Lemma 1:
This is because, otherwise, an alternating path to an unmatched vertex could
have been used to increase the maximum matching M. This would contradict the
fact that M is a maximum matching.
This finishes the proof of Lemma 1.

Lemma 2: |W| >= |Z| + 1
Proof of Lemma 2:
Using Lemma 1: The size of W is at least the number |Z| of these matched
neighbots of Z, plus one
for the unmatched vertex u. Thus, |W| >= |Z| + 1.
This finishes the proof of Lemma 2.

Lemma 3: Z = N(W). In other words, Z is the set of neighbots of W.
Proof of Lemma 3:
You can find an alternating path from v to w either by:
- Removing the matched edge vw from the alternating path to v, or,
- By adding the unmatched edge vw from the alternating path to v.
Thus, for every vertex v in W, every neighbor w belongs to Z.
This finishes the proof of Lemma 3.

Using Lemmas 2 and 3, we have: |W| >= N(W) + 1, which completes the proof.



Phase (C): Suppose a perfect matching M exists:

Definition of match(v):
Given a vertex v, define match(v) to be the vertex matched with v in the
perfect matching M.

Lemma 4:
Let S be a subset of L. Then, |S| = |N(S)| if and only if:
For every edge (l,r) not in M: If S contains l, then S contains match(r).

Proof of Lemma 4, Part 1, Forward implication (=>):
Let's assume |S| = |N(S)|
We must prove for every edge (l,r) not in M, S contains l implies that S
contains match(r).

Proof by contradiction:
Suppose there exists an edge (l,r) not in M, such that S contains l, but S does
 not contain match(r).
Let's count the number of vertices in N(S) carefully:
First, every vertex x in S is matched with a unique vertex match(x) in N(S).
This is a total of |S| vertices.
Second, vertex l has an edge to vertex r.
Since match(r) is not in S, r cannot coincide with the |S| vertices we added
 previously. Thus, we must add 1 more vertex to N(S).
Therefore, |N(S)|>=|S|+1.
This contradicts the initial assumption that |S| = |N(S)|.


Proof of Lemma 4, Part 2, Backwards implication (<=):
Let's assume for every edge (l,r) not in M, S contains l implies that S
 contains match(r).
We must prove that |S| = |N(S)|. Let's split the proof into 2 parts:

Proof that |S| <= |N(S)|:
In Phase (A), we proved using Hall's theorem that if the perfect matching M
 exists, then G is good.
By definition, since G is good, |S| <= |N(S)|.

Proof that |S| >= |N(S)|:
Every vertex r in N(S) is matched with a unique vertex match(r). By the
 assumption, match(r) is in S.
Therefore, |S| >= |N(S)|.

Lemma 4 is very useful because it reduces the tightness criteria into a list
 of implications.
We borrow an idea from 2-SAT algorithm: Convert this list of implications into
 a Directed Acyclic Graph.

We construct a Directed Acyclic Graph D where:
- The set of vertices is equal to L.
- Each edge x -> y represents that "x is in S" implies "y is in S".
We also call D an "implication graph".
Our goal is to find a graph D' with the minimum number of edges, such that the
 transitive closures of D and D' are equal.

Definition of SCC:
An SCC (Strongly Connected Component) is a subset of vertices such that for
 any 2 vertices x,y in this subset, there is a path from x to y, and from y
 to x.

Let us observe which unneeded edges we can remove from this graph, but still
 preserve the list of implications.

There are 2 groups of unneeded edges:
Group (I): Within a SCC of size k in graph D, we can remove all edges except
 the k edges which form a cycle.

Group (II): For an edge connecting SCC X to SCC Y, we can remove this edge if
 there is another path from SCC X to SCC Y.

Proof of why you can remove Group (I) and Group (II) edges:
If you remove a Group (I) edge x->y within an SCC, you can travel from x to y
using the remaining cycle.
If you remove a Group (II) edge x->y, by definition there is another path
from x to y.
Since implication (=>) is a transitive relation, we preserve the list of
 implications.

Proof of why you cannot remove more than Group (I) and Group (II) edges:
For edges within an SCC: We can prove by induction that an SCC with k vertices
 must contain at least k edges.
For edges between SCC-s: If we remove an edge not in Group (II), by definition
 that means we remove an edge from SCC X to SCC Y, but there is no other path
 from SCC X to Y. This breaks the chain of implications.


To implement Phase (C), we do the following:

Step 1: Construct the Directed Acyclic Graph D ("Implication Graph").
- The set of vertices is equal to L.
- For every edge (l,r) in G and not in M: Add the edge l -> match(r) to D.

Step 2: Run Kosaraju's (or Tarjan's) SCC algorithm on the graph D.

Step 3: Within each SCC, keep only the edges which form a cycle.
This removes unneeded edges from Group (I).

Step 4: Contract each SCC into a vertex.
The resulting graph is called a "Condensation Graph" C.

Step 5: Run topological sort on Condensation Graph C.

Step 6: In reverse topological order, for each vertex v, maintain a bitset
which represents which vertices are reachable from v. Using these bitsets,
remove unneeded edges from Group (II).

This completes the implementation of Phase (C).


Final notes:
My solution is slightly different from the official solution.
The official solution uses the concept of submodularity.
You can view the official solution at: https://codeforces.com/blog/entry/117394
*/


#include <algorithm>
#include <bitset>
#include <cstring>
#include <iostream>
#include <queue>
#include <tuple>
#include <vector>

using namespace std;

/*
 * Data structure to represent a graph.
 */
struct Graph {
  Graph() {}

  /*
   * Constructor for a new graph
   * Parameters:
   *   num_vertices: The number of vertices in this graph
   */
  Graph(int num_vertices) {
    adjacency_list.resize(num_vertices);
  }

  /*
   * Adds a new edge to the graph
   * Parameters:
   *   from_vertex: The vertex this edge is pointing from
   *   to_vertex: The vertex this edge is pointing towards
   */
  void add_edge(int from_vertex, int to_vertex) {
    adjacency_list[from_vertex].push_back(to_vertex);
  }

  /*
   * Resizes the graph to the new number off vertices
   */
  void resize(int num_vertices) {
    adjacency_list.resize(num_vertices);
  }

  /*
   * Returns a modifiable reference to the adjacency list of the input vertex.
   */
  vector<int>& operator[](size_t vertex) {
    return adjacency_list[vertex];
  }

  /*
   * Returns the adjacency list of the input vertex.
   */
  vector<int> operator[](size_t vertex) const {
    return adjacency_list[vertex];
  }

  /*
   * Returns the number of vertices in the graph.
   */
  int get_num_vertices() const {
    return static_cast<int>(adjacency_list.size());
  }

  /*
   * Reverses the direction of every edge, then returns the new reversed graph.
   */
  Graph get_reversed_graph() const {
    Graph reversed_graph(get_num_vertices());
    for (int from_vertex = 0; from_vertex < get_num_vertices(); from_vertex++) {
      for (int to_vertex: adjacency_list[from_vertex]) {
        reversed_graph.add_edge(to_vertex, from_vertex);
      }
    }
    return reversed_graph;
  }

  // For any vertex x, adjancency_list[x] contains the neighbors of x.
  vector<vector<int> > adjacency_list;
};


/*
 * Struct to represent an edge in a graph.
 */
struct Edge {

  // Vertex where the edge begins
  int from_vertex = -1;
  // Vertex where the edge ends
  int to_vertex = -1;

  Edge() {}

  /*
   * Constructor for a new edge
   * Parameters:
   *   from_vertex: The vertex this edge is pointing from
   *   to_vertex: The vertex this edge is pointing towards
   */
  Edge(int _from_vertex, int _to_vertex):
    from_vertex(_from_vertex), to_vertex(_to_vertex) {}
  
  // Equal comparator
  bool operator==(const Edge & rhs) {
    return tie(from_vertex, to_vertex) == tie(rhs.from_vertex, rhs.to_vertex);
  }

  // Less-than comparator used for sorting
  bool operator<(const Edge & rhs) {
    return tie(from_vertex, to_vertex) < tie(rhs.from_vertex, rhs.to_vertex);
  }
};

/*
 * Given an input vector of elements, we want to remove duplicates in-place.
 * As a side-effect, we will also sort the elements.
 */
template<typename T> void remove_duplicates(vector<T> & input_vector) {
  sort(input_vector.begin(),input_vector.end());
  auto it=unique(input_vector.begin(),input_vector.end());
  input_vector.resize(distance(input_vector.begin(),it));
}

/*
 * Class to topological sort a given input directed graph using
 * Depth-First Search
 */
class TopologicalSorter {
  // Keeps track which vertices we have already seen in our DFS
  vector<bool> seen;

  // Fills the vector "topologically_sorted_vertices" with vertices in
  // increasing order of their DFS exit times
  void topological_sort_dfs(
    int current_vertex,
    Graph const& graph,
    vector<int> & topologically_sorted_vertices
  ) {
    // Mark vertex as seen, so we don't traverse the same vertex again
    seen[current_vertex] = true;
    for (int child_vertex: graph[current_vertex]) {
      if (!seen[child_vertex]) {
        // Recursively DFS into child vertex
        topological_sort_dfs(child_vertex, graph, topologically_sorted_vertices);
      }
    }
    // As we exit this vertex, we store it into topologically_sorted_vertices
    topologically_sorted_vertices.push_back(current_vertex);
  }

public:
  /*
   * Given an input graph, find a topological order of the vertices, then return them.
   * The vertices are returned in reversed topological order.
   */
  vector<int> reversed_topological_sort(Graph const& graph) {
    vector<int> topologically_sorted_vertices;
    const int num_vertices = graph.get_num_vertices();
    seen.resize(num_vertices);
    seen.assign(num_vertices, false);
    for (int vertex = 0; vertex < num_vertices; vertex++) {
      if (!seen[vertex]) {
        topological_sort_dfs(vertex, graph, topologically_sorted_vertices);
      }
    }
    seen.clear();
    return topologically_sorted_vertices;
  }

  /*
   * Given an input graph, find a topological order of the vertices, then return them.
   * The vertices are returned in topological order.
   */
  vector<int> topological_sort(Graph const& graph) {
    vector<int> topologically_sorted_vertices = reversed_topological_sort(graph);
    reverse(topologically_sorted_vertices.begin(), topologically_sorted_vertices.end());
    return topologically_sorted_vertices;
  }
};

/* Kosaraju's algorithm for Strongly Connected Components.
 * 
 * An SCC (Strongly Connected Component) is a subset of vertices such that for
 *  any 2 vertices x,y in this subset, there is a path from x to y, and from y
 *  to x.
 * 
 * Given an input directed graph, Kosaraju's algorithm will find all SCC-s
 * in the graph.
 * 
 * Kosaraju's algorithm consists of 2 phases:
 * Phase 1: Topological sort on the vertices
 * Phase 2: In reverse order of the topological sort, depth-first-search the
 * edges in reverse direction. Each depth-first-search will give you one
 * strongly connected component.
 */
class KosarajuSCC {
  // Keeps track of whether a vertex has been seen before in the DFS
  vector<bool> seen;

  // After we find the Strongly Connected Components, we can contract each
  // component into a condensation graph.
  // condensed is the adjacency list of the representative_nodes.
  // We can now traverse on condensed as our condensation graph, using only
  // those nodes which belong to representative_nodes.
  Graph condensation_graph;
  // components_list[i] contains the vertex IDs of this component i.
  vector<vector<int> > components_list;
  // representatives[v] indicates the representative node for the SCC to which
  // node v belongs.
  vector<int> representatives;
  // representative_nodes is the list of all representative nodes (one per
  // component) in the condensation graph.
  vector<int> representative_nodes;

  /*
   * This runs "Phase 2" of Kosaraju's algorithm.
   * Given an input vertex v, we will store all vertices in the same strongly
   * connected component as v into the output vector "component".
   */
  void find_scc_dfs(
    int current_vertex,
    vector<int> & component,
    Graph const& reversed_graph
    ) {
    seen[current_vertex] = true;
    component.push_back(current_vertex);
    // Depth-first-search on the reversed graph to find other vertices which
    // belong to this strongly connected component.
    for (int parent_vertex : reversed_graph[current_vertex]) {
      if (!seen[parent_vertex]) {
        find_scc_dfs(parent_vertex, component, reversed_graph);
      }
    }
  }

  /*
   * Constructs the strongly-connected component graph using Kosaraju's
   * algorithm.
   * Parameters:
   *   graph: The input directed graph we want to find SCC-s for.
   */
  void makeSCCHelper(Graph const& graph) {
    // Get the number of vertices in the graph
    const int num_vertices = static_cast<int>(graph.get_num_vertices());

    // Clear and reset all instance variables.
    seen.resize(num_vertices);
    seen.assign(num_vertices, false);

    // Phase 1: Topological sort on the vertices using DFS
    TopologicalSorter topological_sorter;
    vector<int> topologically_sorted_vertices = topological_sorter.topological_sort(graph);

    seen.assign(num_vertices, false);
    representatives.resize(num_vertices);
    condensation_graph.resize(num_vertices);

    // Construct the same graph with the same set of vertices and edges, but
    // each edge's direction is reversed.
    Graph reversed_graph = graph.get_reversed_graph();

    // Phase 2: In topological sorted order, depth-first-search using
    // the edges in reverse direction. Each depth-first-search will give you
    // one strongly connected component.
    for (int current_vertex : topologically_sorted_vertices) {
      if (!seen[current_vertex]) {
        // Vertex v belongs to a new strongly connected component that we
        // have not seen before.
        // Function find_scc_dfs stores all reached vertices in list component,
        // that is going to store next strongly connected component after each
        // run.
        vector<int> component;
        find_scc_dfs(current_vertex, component, reversed_graph);
        // Vector "component" now contains all vertices in this strongly
        // connected component.
        // Designate one vertex as the representative of this component.
        int representative = component.front();
        for (int vertex_in_component : component) {
          representatives[vertex_in_component] = representative;
        }
        representative_nodes.push_back(representative);
        components_list.push_back(component);
      }
    }

    // Construct the condensation graph.
    // It is formed by contracting each strongly connected component into a
    // single vertex.
    for (int current_vertex = 0; current_vertex < num_vertices; current_vertex++) {
      for (int neighbor : graph[current_vertex]) {
        int current_representative = representatives[current_vertex];
        int neighbor_representative = representatives[neighbor];
        if (neighbor_representative != current_representative) {
          condensation_graph.add_edge(current_representative, neighbor_representative);
        }
      }
    }

    // Remove duplicate edges from the adjacency list of the condensed graph.
    for (int vertex = 0; vertex < num_vertices; vertex++) {
      remove_duplicates(condensation_graph[vertex]);
    }
  }

public:
  /*
   * Given an input graph, finds all SCC-s in the graph.
   */
  KosarajuSCC(Graph const& graph) {
    makeSCCHelper(graph);
  }

  /*
   * Gets the list of SCC-s (strongly connected components).
   * We return a vector<vector<int> >
   * The size of the outer vector is equal to the number of SCC-s.
   * Each inner vector is a vector of vertices which belong to this SCC.
   */ 
  const vector<vector<int> > & get_components_list() const {
    return components_list;
  }

  /*
   * If you contract each SCC into a vertex, the resulting graph is called a
   * "Condensation Graph".
   * This function returns the condensation graph.
   */
  const Graph & get_condensation_graph() const {
    return condensation_graph;
  }
};


/*
 * Hopcroft-Karp algorithm for finding maximum matching in a bipartite graph.
 * Input: Bipartite graph with vertex sets L (vertices on the left) and R
 * (vertices on the right), and edge set E.
 * Output: Maximum matching in the bipartite graph.
 * 
 * Definition of an Augmented Path:
 * An "augmented path" is a path in the graph that alternates between matched
 * and unmatched edges.
 * ie. The first edge does not belong to the matching, the second edge belongs
 * to the matching,
 * the third edge belongs to the matching, ...
 * 
 * Brief description:
 * We start with an empty matching M.
 * Then, we repeatedly do the following phase of 3 steps:
 *   (1) Find P, which is the maximal set of vertex-disjoint shortest augmented
 *       paths P1, P2, ..., Pk.
 *   (2) If P is empty, terminate.
 *   (3) Set the new M' to be the XOR of M and the paths P. By XOR, we mean:
 *      The new matching M' will contain an edge E if and only if:
 *      E is in the old M and is not in P, or,
 *      E is not in the old M and is in P
 * 
 * Detailed description:
 * As mentioned, we start with an empty matching M.
 * Then, we repeatedly do phases until we are unable to increase our matching M.
 * Each phase consists of the following steps:
 *
 * Step 1) BFS (Breadth First Search) to find augmented paths.
 * Perform a Breadth First Search that partitions the vertices of the graph
 * into layers.
 * We initialize the queue of the Breadth First Search with unmatched vertices
 * in L.
 * The first level of the search must only contain unmatched edges.
 * At subsequent levels of the search, we alternate between matched and
 * unmatched edges.
 * The BFS terminates at the first layer k, where 1 or more unmatched vertex
 * in L is reached.
 *
 * Step 2) DFS (Depth First Search) to find maximal set of vertex-disjoint
 * augmented paths of length k.
 * Define F to be the set of unmatched vertices (in L) at layer k in Step (1).
 * Start your DFS from vertices in F, to unmatched vertices in R.
 * At layer i, the DFS must follow edges that lead to an unused vertex in
 * layer i-1.
 * Once an augmented path is found, we use the path to enlarge the matching
 * M using the XOR method described above.
 * Then, we continue DFS from the next starting vertex in F.
 * 
 * If no augmented paths are found, we have found a maximum matching, so
 * terminate the algorithm.
 */
class HopcroftKarp {

  // Constant to indicate an unmatched vertex
  const int UNMATCHED = -1;

  // Number of vertices in the left partition
  int num_left_vertices = 0;

  // Number of vertices in the right partition
  int num_right_vertices = 0;

  // Input bipartite graph containing edges from the left partition to the
  // right partition
  Graph graph_from_left;

  // Input bipartite graph containing edges from the right partition to the
  // left partition
  Graph graph_from_right;

  // Maximal matching we have found so far
  // For any left vertex x, match_from_left[x] is the vertex matched with x
  // If x is unmatched, then match_from_left[x] is -1 (UNMATCHED).
  vector<int> match_from_left;

  // Maximal matching we have found so far
  // For any right vertex x, match_from_right[x] is the vertex matched with x
  // If x is unmatched, then match_from_right[x] is -1 (UNMATCHED).
  vector<int> match_from_right;

  // Distance in terms of the "BFS step" of the algorithm
  vector<int> dist;

  // Indicates whether we have run the matching algorithm
  bool got_max_matching = false;

 /* 
  * This function performs Step 1 ("BFS step") of each phase in the
  * Hopcroft-Karp algorithm:
  * Perform a Breadth First Search that partitions the vertices of the graph
  *  into layers.
  * We initialize the queue of the Breadth First Search with unmatched vertices
  *  in L.
  * The first level of the search must only contain unmatched edges.
  * At subsequent levels of the search, we alternate between matched and
  *  unmatched edges.
  * The BFS terminates at the first layer k, where 1 or more unmatched vertex
  *  in L is reached.
  * If no augmented paths are found, we have found a maximum matching, so
  *  terminate the algorithm.
  */
  void bfs() {
    queue<int> vertex_queue;
    for (int left_vertex = 0; left_vertex < num_left_vertices; ++left_vertex) {
      if (match_from_left[left_vertex] == UNMATCHED) {
        vertex_queue.push(left_vertex);
        dist[left_vertex] = 0;
      }
      else {
        dist[left_vertex] = -1;
      }
    }
    while (!vertex_queue.empty()) {
      const int left_vertex = vertex_queue.front();
      vertex_queue.pop();
      for (int right_vertex : graph_from_left[left_vertex])
        if (match_from_right[right_vertex] != UNMATCHED && dist[match_from_right[right_vertex]] == UNMATCHED) {
          dist[match_from_right[right_vertex]] = dist[left_vertex] + 1;
          vertex_queue.push(match_from_right[right_vertex]);
        }
    }
  }

 /*
  * This function performs Step 2 ("DFS step") of each phase in the
  * Hopcroft-Karp algorithm:
  * Define F to be the set of unmatched vertices (in L) at layer k in Step (1).
  * Start your DFS from vertices in F, to unmatched vertices in R.
  * At layer i, the DFS must follow edges that lead to an unused vertex in
  * layer i-1.
  * Once an augmented path P is found, we use the path to enlarge the matching
  * M using the "XOR" method.
  * By "XOR", we mean that:
  *      The new matching M' will contain an edge E if and only if:
  *      E is in the old M and is not in P, or,
  *      E is not in the old M and is in P
  */
  bool dfs(int left_vertex) {
    for (int right_vertex : graph_from_left[left_vertex])
      if (match_from_right[right_vertex] == UNMATCHED) {
        match_from_left[left_vertex] = right_vertex;
        match_from_right[right_vertex] = left_vertex;
        return true;
      }
    for (int right_vertex : graph_from_left[left_vertex])
      if (
        dist[match_from_right[right_vertex]] == dist[left_vertex] + 1 &&
        dfs(match_from_right[right_vertex])
        ) {
        match_from_left[left_vertex] = right_vertex;
        match_from_right[right_vertex] = left_vertex;
        return true;
      }
    return false;
  }

public:

  /*
   * Constructor:
   * Parameters:
   *   input_num_left_vertices:
   *     Number of vertices in the left partition of your bipartite graph.
   *   input_num_right_vertices:
   *     Number of vertices in the right partition of your bipartite graph.
   */
  HopcroftKarp(int input_num_left_vertices, int input_num_right_vertices)
    : num_left_vertices(input_num_left_vertices),
      num_right_vertices(input_num_right_vertices),
      graph_from_left(input_num_left_vertices),
      graph_from_right(input_num_right_vertices),
      match_from_left(input_num_left_vertices, UNMATCHED),
      match_from_right(input_num_right_vertices, UNMATCHED),
      dist(input_num_left_vertices) {}

  /*
   * Add an edge from the left partition to the right partition in the
   * bipartite graph.
   * Parameters:
   *   left_vertex: Vertex on the left partition
   *   right_vertex: Vertex on the right partition
   */
  void add_edge(int left_vertex, int right_vertex) {
    graph_from_left[left_vertex].push_back(right_vertex);
    graph_from_right[right_vertex].push_back(left_vertex);
  }

  /*
   * Finds the maximum matching in the bipartite graph.
   * Returns the size of the matching.
   */
  int get_max_matching() {
    int total_flow = 0;
    while (true) {
      // Call Step 1 (The "BFS" step)
      bfs();

      // Call Step 2 (The "DFS" step)
      // Define F to be the set of unmatched vertices (in L) at layer k in
      // Step (1).
      // Start your DFS from vertices in F, to unmatched vertices in R.
      int augmented_flow = 0;
      for (int u = 0; u < num_left_vertices; ++u)
        if (match_from_left[u] == UNMATCHED) {
          augmented_flow += dfs(u);
        }
      
      // If no augmented paths found, our matching cannot be enlarged further,
      // terminate the algorithm.
      if (!augmented_flow) {
        break;
      }
      total_flow += augmented_flow;
    }
    got_max_matching = true;
    return total_flow;
  }

  /*
   * Find a subset S of the left partition such that |S| > |N(S)|.
   * 
   * Step (1): Assume we have found a non-perfect maximum matching M.
   * Step (2): Choose any vertex u in L that does not belong to the matching M.
   * Step (3): Define an alternating path to be a path that alternately uses
   *  edges outside and inside M.
   * Step (4): Define P to be all alternating paths starting from vertex u,
   *  using breadth-first-search.
   * Step (5): Define W to be the set of vertices in the alternating paths P
   *  that belong to L.
   */
  vector<int> find_nontight_set() {
    if (!got_max_matching) {
      // If we have not run the matching algorithm, then run it first
      get_max_matching();
    }

    vector<bool> seen(num_left_vertices + num_right_vertices);

    // Choose any vertex u in L that does not belong to the matching M.
    int u = UNMATCHED;
    for (int l = 0; l < num_left_vertices; l++) {
      if (match_from_left[l] == UNMATCHED) {
        u = l;
        break;
      }
    }
    if (u == UNMATCHED) {
      // We were not able to find a non-tight set because there is a L-perfect
      // matching.
      return {};
    }

    // Find P, which is all alternating paths starting from vertex u, using
    // breadth-first-search.
    queue<int> vertex_queue;
    vertex_queue.push(u);
    seen[u] = true;
    while(!vertex_queue.empty()) {
      int current_vertex = vertex_queue.front();
      vertex_queue.pop();
      if (current_vertex < num_left_vertices) {
        // If this is a vertex on the left partition, BFS all neighbors using
        // edges not in the matching.
        for (int neighbor: graph_from_left[current_vertex]) {
          if (
            !seen[neighbor + num_left_vertices] &&
            match_from_left[current_vertex] != neighbor
            )  {
            seen[neighbor + num_left_vertices] = true;
            vertex_queue.push(neighbor + num_left_vertices);
          }
        }
      } else {
        // If this is a vertex on the right partition, BFS only using edges
        // in the matching.
        int neighbor = match_from_right[current_vertex - num_left_vertices];
        if (neighbor != -1 && !seen[neighbor]) {
          seen[neighbor] = true;
          vertex_queue.push(neighbor);
        }
      }
    }

    // Nontight set is the set of vertices in the alternating paths P that
    // belong to L.
    vector<int> nontight_set;
    for (int x = 0; x < num_left_vertices; x++) {
      if (seen[x]) {
        nontight_set.push_back(x);
      }
    }
    return nontight_set;
  }

  /*
   * Gets the matching.
   * We return a vector<int> match where:
   * For any vertex x, match[x] is the vertex that is matched with x.
   */
  vector<int> get_matching() {
    if (!got_max_matching) {
      // If we have not run the matching algorithm, then run it first
      get_max_matching();
    }
    vector<int> match(num_left_vertices + num_right_vertices, -1);
    for (int l = 0; l < num_left_vertices; l++) {
      match[l] = match_from_left[l] + num_left_vertices;
    }
    for (int r = 0; r < num_right_vertices; r++) {
      match[r + num_left_vertices] = match_from_right[r];
    }
    return match;
  }
};


/*
 * Given a directed acyclic graph D, we want to find the minimal graph D' such
 * that the transitive closures of D and D' are equal.
 * This allows us the remove "Group (II) unneeded edges".
 * Recall the definition of a Group (II) edge:
 * For an edge connecting SCC X to SCC Y, we can remove this edge if there is
 *  another path from SCC X to SCC Y.
 * 
 * To implement this, we do the following:
 * In reverse topological sort order, for each vertex x, we will compute a
 * bitset reachable[x] such that:
 * reachable[x][y] is true if vertex y is reachable from vertex x via a path
 * of at least 2 edges.
 * 
 * Then, an edge x -> y is unneeded if and only if reachable[x][y] is true.
 * 
 * Parameters:
 *   graph: The input directed graph we want to minimize
 */
template<size_t MaxVertices> vector<Edge> minimize_edges_in_transitive_graph(Graph const& graph) {
  // Get number of vertices in the graph
  int num_vertices = static_cast<int>(graph.get_num_vertices());

  // Run topological sort on the graph
  TopologicalSorter topological_sorter;
  vector<int> reversed_topological_sorted_vertices = topological_sorter.reversed_topological_sort(graph);

  // Compute the reachable bitsets
  vector<bitset<MaxVertices> > reachable(num_vertices);
  vector<Edge> minimal_edges;
  for (const int current_vertex: reversed_topological_sorted_vertices) {
    for (int child_vertex: graph[current_vertex]) {
      // If vertex x is reachable by the child vertex, then by transitivity, x
      // is also reachable by the current vertex.
      reachable[current_vertex] |= reachable[child_vertex];
    }

    for (int child_vertex: graph[current_vertex]) {
      if (!reachable[current_vertex][child_vertex]) {
        // If we cannot reach the child vertex from the current vertex via a
        // longer path, then this edge is needed.
        minimal_edges.push_back(Edge(current_vertex, child_vertex));
      }
    }

    for (int child_vertex: graph[current_vertex]) {
      // Indicate that we can reach the child vertex from the current vertex
      reachable[current_vertex][child_vertex] = true;
    }
  }
  return minimal_edges;
}


int main() {
  // Speed up cin input and cout output by disabling syncing with stdio functions
  ios_base::sync_with_stdio(false);
  cin.tie(0);

  // Number of vertices in the left (and right) partition in bipartite graph G
  int num_vertices;

  // Number of edges in bipartite graph G
  int num_edges;

  // Read from standard input
  cin >> num_vertices >> num_edges;

  HopcroftKarp hopcroft_karp(num_vertices, num_vertices);
  vector<Edge> edges;
  for (int edge_index = 0; edge_index < num_edges; edge_index++) {
    // Read the edge
    int left_vertex, right_vertex;
    cin >> left_vertex >> right_vertex;

    // Subtract 1 to make vertex ID zero-indexed instead of one-indexed.
    --left_vertex;
    --right_vertex;

    // Add edge from vertex l to vertex r
    hopcroft_karp.add_edge(left_vertex, right_vertex - num_vertices);
    edges.push_back(Edge(left_vertex, right_vertex));
  }

  // This is Phase (A) of our solution.
  // Run Hopcroft-Karp bipartite matching algorithm
  const int maxflow = hopcroft_karp.get_max_matching();

  // Check if perfect matching exists
  if (maxflow < num_vertices) {
    // This is Phase (B) of our solution.
    // If no perfect matching, find a set of vertices which is not tight.
    vector<int> nontight_set = hopcroft_karp.find_nontight_set();
    cout << "NO\n";
    cout << nontight_set.size() << '\n';
    for (int x: nontight_set) {
      cout << x + 1 << ' ';
    }
    cout << endl;
  } else {
    // This is Phase (C) of our solution.
    // We found a perfect matching M.
    // For any vertex v, match[v] is the vertex that is matched with v by the matching M.
    const vector<int> match = hopcroft_karp.get_matching();

    // Construct the Directed Acyclic Graph ("Implication Graph").
    // - The set of vertices is equal to L.
    // - For every edge (l,r) in G and not in M: Add the edge l -> match(r).
    Graph implication_graph(num_vertices);
    for (Edge & edge: edges) {
      int left_vertex = edge.from_vertex, right_vertex = edge.to_vertex;
      if (match[left_vertex] != right_vertex) {
        implication_graph.add_edge(left_vertex, match[right_vertex]);
      }
    }

    // List of edges in the minimal graph H we want to output
    vector<Edge> minimal_edge_list;

    // Run Kosaraju's Strongly Connected Components algorithm on the implication graph.
    KosarajuSCC scc(implication_graph);
    for (int vertex = 0; vertex < num_vertices; vertex++) {
      minimal_edge_list.push_back(Edge(vertex, match[vertex]));
    }

    // Try to minimize number of edges within each SCC (strongly-connected component)
    for (const vector<int> & comp: scc.get_components_list()) {
      const int component_size = static_cast<int>(comp.size());
      if (component_size > 1) {
        // We have a nontrivial SCC with many edges.
        // Our goal: Use the minimum number of edges, but still keep the
        // component an SCC.
        // Within each SCC, we will only add the edges which form a cycle.
        for (int current_index = 0; current_index < component_size; current_index++) {
          int next_index = (current_index + 1) % component_size;
          minimal_edge_list.push_back(Edge(comp[current_index], match[comp[next_index]]));
        }
      }
    }

    // For an edge connecting SCC X to SCC Y, we can remove this edge if
    // there is another path from SCC X to SCC Y.
    // In reverse topological order, for each vertex v, maintain a bitset
    // which represents which vertices are reachable from v. Using these
    // bitsets, remove unneeded edges.
    const int MAX_VERTICES = 1000;
    const vector<Edge> minimal_condensation_edges = \
      minimize_edges_in_transitive_graph<MAX_VERTICES>(
        scc.get_condensation_graph());

    // Add the minimal edges between SCC-s to the minimal graph H
    for (Edge const& edge: minimal_condensation_edges) {
      minimal_edge_list.push_back(Edge(edge.from_vertex, match[edge.to_vertex]));
    }

    // Remove duplicate edges
    remove_duplicates(minimal_edge_list);

    // Print YES to indicate that the input graph is good.
    cout << "YES\n";

    // Print the number of edges.
    cout << minimal_edge_list.size() << '\n';

    for (const Edge & edge: minimal_edge_list) {
      // Print each edge.
      cout << edge.from_vertex + 1 << ' ' << edge.to_vertex + 1 << '\n';
    }
    cout << endl;
  }
}
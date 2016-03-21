# Dijkstra performance

Dijkstra's algorithm is widely shown to have O(n log n) asymtotic performance.
To test performance random graphs are constructed with graph.EuclideanLabeled and random start nodes are chosen.  As shown in the example here, these graphs can have sinks and even isolated nodes.  Start nodes are are chosen to not be sink nodes but are otherwise random.

![Random graph](https://cdn.rawgit.com/soniakeys/graph/svg-v0/bench/img/Euclidean.svg)

The Benchmark function of the testing package is useful for accurate timings on a given graph from a given start node.  Graphs of various sizes are constructed with order N a power of 2 ranging from 2^4 through 2^22.  On the machine used for testing these graphs fit easily in RAM and don't take terribly long to test.  Each graph is tested with five randomly chosen start nodes.  The minimum and maximum of the five times for each graph are shown here with error bars.  Also an n log n curve is fit and shown with the solid line

![Dijkstra performance](https://cdn.rawgit.com/soniakeys/graph/svg-v0/bench/img/DijkstraAllPaths.svg)

The plot shows the implementation performing near the expected O(n log n).  Inspection of the code shows that there is certainly an O(n) term but the plot indicates this term must have little significance.  Also complications of memory management and memory architecture can be expected to add additional terms but these too appear to have little significance over the range tested.
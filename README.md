### Monte Carlo Tree Search
This is a simple Go script written to run a Monte Carlo tree search (MCTS) on a Gomoku board. A future goal is to use this script in conjunction with a reinforcement learning model as a means of policy improvement.  

Because MCTS is a highly repetitive, non-matrix based algorithm I chose to implement it in a language that is faster than Python. Moreover, as there are ways to parallelize MCTS, I wanted to use a language with easy-to-use concurrency features. Go seemed liked a natural fit for this project. 

Currently, this script allows for basic parallelization by generating many trees from the same root in parallel and then averaging the results. This leads to rapid memory growth that can terminate the program prematurely. I'm currently working to fix this. 
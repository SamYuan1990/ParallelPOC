# Trun for normal nodes

# Producer
receive txs array
1. Convert txs array into nodes array
2. While the nodes array is not empty, put into channel of nodes

# Processor
1. Read from channel
2. Process

# Turn for supports U
# Producer
receive txs array
1. Convert txs array into nodes array
2. While the nodes array is not empty, put into channel of nodes
3. Wait Group

# Processor
1. Read from channel
2. Process
3. Wait Group
4. Tx channels as copy and return terms
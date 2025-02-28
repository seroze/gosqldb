# gosqldb

I'll be posting my about's while writing an SQL db,
i'm following `How to write your own database from scratch in Go` book

Book summary:

Chapter-1: (Persistance/How to Update to disk-data)

- it's better to append then to rewrite some data.
- avoid inplace updates
- take care of fsycning
- renaming is not atomic, you have to do fsync in parent dir(in case there's powerloss)
- log is not an index

Chapter-2: (Index structures)

- There are 3 types of SQL queries (point query, range query, full scan)
- we don't need index for full scan queries
- If we only care about point query, hash-table is good as an in-memory data structure. (Re-hashing is expensive sometimes if load factor is exceeded)
- we can start with sorted lists as in-memory ds, but insertion is costly
- insertion cost can be kept down if we do something like square-root decomposition
- btree is nothing but multi-level sorted lists
- Important points about IO
  - disk can only perform limited no of IO's per second (IOPS)
  - each btree look up is a disk io so the shorter the tree the better
  - basic unit of disk-io is sectors(which are contigous 512 byte blocks on disk on old devices)
  - disk sectors is not concern for app kernel caches pages (usually 4K byte blocks)
- B+ tree is a B-tree but values are only stored in last level
- log structured merge tree (if db works by appending queries to a log and uses it as structure then it's LSM like it only works with logs)
- There are two kinds of indexes BTree and LSM trees depending on your use case

Chapter-3: (B+Tree recovery and crash)
- In a B+Tree values are present in leaves and keys are duplicated in internal nodes
- B+Tree invariants
  - all leaves are at same level
  - node size never grows beyond a constant
  - no node is empty
- BTree on disk
  - one can implement B+Tree following above rules, but how do we ensure that the size is within the block limit ?
  - that's why we should create B+Tree nodes with same size and have freelist of nodes. (There's no malloc/free we have to do everything )//I didn't get this point clearly
- Copy on write B-tree's for safe updates
  - we need to update the node till the root following it's ancestors it's a O(LogN) operation
  - there's way to improve about which is called double-write // read more on this
-B+Tree node insertion and deletion
  -

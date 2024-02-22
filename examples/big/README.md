# gree Example - big
This example illustrates some of the more complex use cases available from the package. 

It creates nodes with children, examining their initial depth and then shows the recalculation of depth after thier parent is added to another node. 

It also demonstrats retrieving nodes by generation and index and setting their colors.

Finally, it shows all of the DrawInput options available. 

Text output:
```
┌─────────────────────┐
│ root                │
│ ├─ child0           │
│ │  ├─ grandchild0   │
│ │  │  ├─ apple0     │
│ │  │  │  └─ oranges0│
│ │  │  └─ bob        │
│ │  ├─ grandchild1   │
│ │  │  └─ apple1     │
│ │  │     └─ oranges1│
│ │  └─ grandchild2   │
│ │     └─ apple2     │
│ │        └─ oranges2│
│ ├─ child1           │
│ │  ├─ grandchild0   │
│ │  │  └─ apple3     │
│ │  │     └─ oranges3│
│ │  ├─ grandchild1   │
│ │  │  └─ apple4     │
│ │  │     └─ oranges4│
│ │  └─ grandchild2   │
│ │     └─ apple5     │
│ │        └─ oranges5│
│ └─ one              │
│    └─ two           │
│       ├─ three      │
│       │  ├─ four    │
│       │  │  └─ five │
│       │  └─ oranges6│
│       └─ apple6     │
│          └─ oranges7│
└─────────────────────┘

|....|....|....|....|..
0    5    10   15   20
```

Screenshot to show colors

![](./img/big.png)
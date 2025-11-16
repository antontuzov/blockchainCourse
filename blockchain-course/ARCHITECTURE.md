# Blockchain System Architecture

## Overview
This document provides a visual representation of the blockchain system architecture implemented in this course.

## Core Components Architecture

```mermaid
graph TB
    A[Blockchain Client] --> B[CLI Interface]
    A --> C[REST API]
    A --> D[WebSocket API]
    
    B --> E[Blockchain Core]
    C --> E
    D --> E
    
    E --> F[Block Structure]
    E --> G[Consensus Engine]
    E --> H[Transaction Processor]
    E --> I[Wallet Manager]
    
    F --> J[Persistence Layer]
    G --> K[Proof of Work]
    G --> L[Proof of Stake]
    H --> M[UTXO Model]
    H --> N[Script Engine]
    I --> O[Key Store]
    
    J --> P[LevelDB]
    K --> Q[Difficulty Adjustment]
    L --> R[Staking Mechanism]
    M --> S[Transaction Validation]
    N --> T[Smart Contracts]
    O --> U[Encryption]
    
    E --> V[P2P Network]
    V --> W[Node Discovery]
    V --> X[Message Broadcasting]
    V --> Y[Block Sync]
    
    W --> Z[Peer Management]
    X --> AA[Transaction Propagation]
    Y --> AB[Chain Sync]
    
    subgraph "External Integrations"
        AC[Web Dashboard]
        AD[Mobile App]
        AE[Third-party Services]
    end
    
    AC --> C
    AD --> C
    AE --> C
```

## Data Flow Diagram

```mermaid
graph LR
    A[User] --> B[Create Transaction]
    B --> C[Sign Transaction]
    C --> D[Submit to Network]
    D --> E[Broadcast to Peers]
    E --> F[Validate Transaction]
    F --> G[Add to Mempool]
    G --> H[Mine New Block]
    H --> I[Add Block to Chain]
    I --> J[Propagate Block]
    J --> K[Update Wallet Balances]
```

## Module Dependencies

```mermaid
graph TD
    A[Module 1 - Fundamentals] --> B[Module 2 - Transactions]
    B --> C[Module 3 - Networking]
    C --> D[Module 4 - Enterprise]
    D --> E[Module 5 - Production]
    
    A --> F[Week 1 - Cryptography]
    A --> G[Week 2 - Blockchain]
    B --> H[Week 3 - UTXO Model]
    B --> I[Week 4 - Smart Contracts]
    C --> J[Week 5 - P2P Network]
    C --> K[Week 6 - Consensus]
    D --> L[Week 7 - Permissioned]
    D --> M[Week 8 - Scalability]
    E --> N[Week 9 - Security]
    E --> O[Week 10 - Deployment]
```

## Capstone Project Options

```mermaid
graph TD
    A[Capstone Projects] --> B[Cryptocurrency]
    A --> C[Enterprise Platform]
    A --> D[DeFi Protocol]
    
    B --> E[Wallet System]
    B --> F[Mining System]
    B --> G[P2P Network]
    
    C --> H[Access Control]
    C --> I[Smart Contracts]
    C --> J[Multi-tenant]
    
    D --> K[DEX]
    D --> L[Lending/Borrowing]
    D --> M[Liquidity Pools]
```
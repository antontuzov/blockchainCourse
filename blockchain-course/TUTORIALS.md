# Interactive Coding Tutorials

## Overview
This document provides step-by-step interactive tutorials to enhance your learning experience in blockchain development.

## Tutorial 1: Building Your First Block

### Objective
Create a simple block structure and understand its components.

### Steps

1. **Create the Block Structure**
   ```go
   type Block struct {
       Index         int64
       Timestamp     int64
       Data          []byte
       PrevBlockHash []byte
       Hash          []byte
       Nonce         int
   }
   ```

2. **Implement the NewBlock Function**
   ```go
   func NewBlock(data string, prevBlockHash []byte) *Block {
       block := &Block{
           Index:         0,
           Timestamp:     time.Now().Unix(),
           Data:          []byte(data),
           PrevBlockHash: prevBlockHash,
           Hash:          []byte{},
           Nonce:         0,
       }
       block.SetHash()
       return block
   }
   ```

3. **Implement the SetHash Function**
   ```go
   func (b *Block) SetHash() {
       headers := [][]byte{
           b.PrevBlockHash,
           b.Data,
           IntToHex(b.Timestamp),
           IntToHex(int64(b.Nonce)),
       }
       header := bytes.Join(headers, []byte{})
       hash := sha256.Sum256(header)
       b.Hash = hash[:]
   }
   ```

4. **Test Your Implementation**
   ```go
   func main() {
       block := NewBlock("Hello, Blockchain!", []byte{})
       fmt.Printf("Block Hash: %x\n", block.Hash)
   }
   ```

### Challenge
Modify the block structure to include a Merkle root of transactions.

## Tutorial 2: Implementing Proof of Work

### Objective
Implement a Proof of Work consensus mechanism for your blockchain.

### Steps

1. **Create the ProofOfWork Structure**
   ```go
   type ProofOfWork struct {
       block  *Block
       target *big.Int
   }
   ```

2. **Implement NewProofOfWork Function**
   ```go
   func NewProofOfWork(block *Block) *ProofOfWork {
       target := big.NewInt(1)
       target.Lsh(target, uint(256-targetBits))
       
       pow := &ProofOfWork{block: block, target: target}
       return pow
   }
   ```

3. **Implement the Run Method**
   ```go
   func (pow *ProofOfWork) Run() (int, []byte) {
       var hashInt big.Int
       var hash [32]byte
       var nonce int = 0
       
       fmt.Printf("Mining a new block...")
       for nonce < math.MaxInt64 {
           data := pow.prepareData(nonce)
           hash = sha256.Sum256(data)
           fmt.Printf("\r%x", hash)
           hashInt.SetBytes(hash[:])
           
           if hashInt.Cmp(pow.target) == -1 {
               break
           } else {
               nonce++
           }
       }
       fmt.Print("\n\n")
       
       return nonce, hash[:]
   }
   ```

4. **Implement prepareData Method**
   ```go
   func (pow *ProofOfWork) prepareData(nonce int) []byte {
       data := bytes.Join(
           [][]byte{
               pow.block.PrevBlockHash,
               pow.block.Data,
               IntToHex(pow.block.Timestamp),
               IntToHex(int64(targetBits)),
               IntToHex(int64(nonce)),
           },
           []byte{},
       )
       
       return data
   }
   ```

5. **Test Your Implementation**
   ```go
   func main() {
       block := NewBlock("Test Block", []byte{})
       pow := NewProofOfWork(block)
       nonce, hash := pow.Run()
       
       block.Hash = hash[:]
       block.Nonce = nonce
       
       fmt.Printf("Block mined: %x\n", block.Hash)
   }
   ```

### Challenge
Implement difficulty adjustment based on mining time.

## Tutorial 3: Creating a Simple Blockchain

### Objective
Build a complete blockchain with multiple blocks and validation.

### Steps

1. **Create the Blockchain Structure**
   ```go
   type Blockchain struct {
       Blocks []*Block
   }
   ```

2. **Implement NewBlockchain Function**
   ```go
   func NewBlockchain() *Blockchain {
       return &Blockchain{Blocks: []*Block{GenesisBlock()}}
   }
   ```

3. **Implement GenesisBlock Function**
   ```go
   func GenesisBlock() *Block {
       return NewBlock("Genesis Block", []byte{})
   }
   ```

4. **Implement AddBlock Method**
   ```go
   func (bc *Blockchain) AddBlock(data string) {
       prevBlock := bc.Blocks[len(bc.Blocks)-1]
       newBlock := NewBlock(data, prevBlock.Hash)
       
       // Mine the block
       pow := NewProofOfWork(newBlock)
       nonce, hash := pow.Run()
       
       newBlock.Hash = hash[:]
       newBlock.Nonce = nonce
       
       bc.Blocks = append(bc.Blocks, newBlock)
   }
   ```

5. **Implement IsValid Method**
   ```go
   func (bc *Blockchain) IsValid() bool {
       for i := 1; i < len(bc.Blocks); i++ {
           currentBlock := bc.Blocks[i]
           prevBlock := bc.Blocks[i-1]
           
           // Check if the current block's hash is valid with proof of work
           pow := NewProofOfWork(currentBlock)
           if !pow.Validate() {
               return false
           }
           
           // Check if the previous block hash matches
           if !bytes.Equal(currentBlock.PrevBlockHash, prevBlock.Hash) {
               return false
           }
       }
       return true
   }
   ```

6. **Test Your Implementation**
   ```go
   func main() {
       bc := NewBlockchain()
       
       bc.AddBlock("First transaction")
       bc.AddBlock("Second transaction")
       bc.AddBlock("Third transaction")
       
       for i, block := range bc.Blocks {
           fmt.Printf("Block %d:\n", i)
           fmt.Printf("  Data: %s\n", block.Data)
           fmt.Printf("  Hash: %x\n", block.Hash)
           fmt.Printf("  Previous Hash: %x\n", block.PrevBlockHash)
           fmt.Printf("  Nonce: %d\n", block.Nonce)
           fmt.Println()
       }
       
       if bc.IsValid() {
           fmt.Println("Blockchain is valid!")
       } else {
           fmt.Println("Blockchain is invalid!")
       }
   }
   ```

### Challenge
Add a CLI interface to interact with your blockchain.

## Tutorial 4: Implementing Transactions

### Objective
Create a UTXO-based transaction system.

### Steps

1. **Create Transaction Structures**
   ```go
   type Transaction struct {
       ID   []byte
       Vin  []TXInput
       Vout []TXOutput
   }
   
   type TXInput struct {
       Txid      []byte
       Vout      int
       Signature []byte
       PubKey    []byte
   }
   
   type TXOutput struct {
       Value      int
       PubKeyHash []byte
   }
   ```

2. **Implement NewCoinbaseTX Function**
   ```go
   func NewCoinbaseTX(to, data string) *Transaction {
       if data == "" {
           data = fmt.Sprintf("Reward to \"%s\"", to)
       }
       
       txin := TXInput{[]byte{}, -1, nil, []byte(data)}
       txout := NewTXOutput(10, to)
       tx := Transaction{nil, []TXInput{txin}, []TXOutput{*txout}}
       tx.ID = tx.Hash()
       
       return &tx
   }
   ```

3. **Implement Hash Method**
   ```go
   func (tx *Transaction) Hash() []byte {
       var hash [32]byte
       
       txCopy := *tx
       txCopy.ID = []byte{}
       
       hash = sha256.Sum256(txCopy.Serialize())
       return hash[:]
   }
   ```

4. **Test Your Implementation**
   ```go
   func main() {
       coinbase := NewCoinbaseTX("Alice", "Genesis block reward")
       fmt.Printf("Coinbase Transaction ID: %x\n", coinbase.ID)
   }
   ```

### Challenge
Implement transaction signing and verification.

## Tutorial 5: Building a P2P Network

### Objective
Create a simple peer-to-peer network for blockchain synchronization.

### Steps

1. **Create Node Structure**
   ```go
   type Node struct {
       Address    string
       Port       int
       Peers      map[string]*Peer
       Blockchain *Blockchain
   }
   
   type Peer struct {
       Address string
       Port    int
       Conn    net.Conn
   }
   ```

2. **Implement NewNode Function**
   ```go
   func NewNode(address string, port int, blockchain *Blockchain) *Node {
       return &Node{
           Address:    address,
           Port:       port,
           Peers:      make(map[string]*Peer),
           Blockchain: blockchain,
       }
   }
   ```

3. **Implement AddPeer Method**
   ```go
   func (n *Node) AddPeer(address string, port int) error {
       peerAddress := fmt.Sprintf("%s:%d", address, port)
       
       // Connect to peer
       conn, err := net.Dial("tcp", peerAddress)
       if err != nil {
           return err
       }
       
       // Create peer
       peer := &Peer{
           Address: address,
           Port:    port,
           Conn:    conn,
       }
       
       // Add peer to map
       n.Peers[peerAddress] = peer
       
       fmt.Printf("Added peer: %s\n", peerAddress)
       return nil
   }
   ```

4. **Test Your Implementation**
   ```go
   func main() {
       bc := NewBlockchain()
       node := NewNode("localhost", 3000, bc)
       
       // Add a peer
       err := node.AddPeer("localhost", 3001)
       if err != nil {
           fmt.Printf("Error adding peer: %s\n", err)
       }
   }
   ```

### Challenge
Implement block broadcasting between nodes.
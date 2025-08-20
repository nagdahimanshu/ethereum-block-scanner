# Backend Interview – Crypto

## Mandatory task

Given a list of 500,000 Ethereum addresses, each associated to a `userId`, create a microservice in Golang that monitors the Ethereum blockchain for any transactions involving those addresses. In summary, the service should:

1. Connect to the Ethereum blockchain via **native JSON-RPC methods** (e.g. using Alchemy, QuickNode, or any other free RPC provider).  
   Do not use third-party APIs or indexing services — only the RPC methods exposed by Ethereum nodes.

2. Detect all transactions that involve the specified addresses.  

3. For the filtered transactions, output the following information:
   - `userId`
   - `from`
   - `to`
   - `amount`
   - `hash`
   - `blockNumber`

The service should be designed for scalability, capable of processing blocks in real time. Assume you have 500,000 users (therefore 500,000 unique addresses).

We value modular, simple, testable code. Showcase how testable it is by testing it :)

---

## Bonus task (not mandatory)

- Add a **Mermaid diagram** to illustrate your solution.  
- Add **Kafka integration** to publish the output as events.  
- Explain (no need to code) how you would handle edge cases such as retry situations, block reorganization, and how to not lose any transactions in a downtime scenario of the blockchain node.

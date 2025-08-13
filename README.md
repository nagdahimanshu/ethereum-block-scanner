# backend-interview-crypto

## Mandatory task
Given a list of 500,000 Ethereum addresses, each associated to a `userId`, create a microservice in Golang that monitors the Ethereum blockchain for any transactions involving those addresses. In summary, the service should:

1. Connect via RPC to the Ethereum blockchain using Blockdaemon/Alchemy/Infura/Quicknode (feel free to use another provider if you prefer!).
 
2. Consume the appropriate information to detect all incoming transactions that involve the specified addresses.
 
3. For the filtered transactions, process the payload and output the following information:
- Source
- Destination
- Amount
- Fees

This output should be in the form of an event emitted by a Kafka producer.

The service should be designed for scalability, capable of processing blocks in real time. Assume you have 500,000 users (therefore 500,000 unique addresses).

We value modular, simple, testable code. Showcase how testable it is by testing it :)

## Bonus task 
Not mandatory, but appreciated:

Bonus 1: Add a Mermaid diagram to illustrate your solution. 
Bonus 2: Explain (no need to code) how you would handle edge cases like retry situations, block reorganization, how to not lose any txs in a 1h downtime scenario of the blockchain node, and any other scenarios you want to showcase.

## RPC Docs

You can use any RPC provider you want. If you need an example or reference

[RPC ethereum](https://docs.blockdaemon.com/reference/how-to-access-ethereum-api)

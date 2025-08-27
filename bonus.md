## What We Can Add for Block Reorganization

### Current state:
- We maintain a checkpoint (lastProcessedBlock) and process blocks only after N confirmations, which protects against short reorgs.
However, there is no logic to re-fetch or reconcile blocks between the checkpoint and the current node height. This means that if the scanner was down or a longer reorg occurs, some blocks and transactions could be skipped.

### Extension to handle reorgs:
- On startup, fetch all blocks from lastProcessedBlock + 1 up to currentHeight - N confirmations.
- Reprocess these blocks through the same confirmation queue logic.
- Continue publishing only confirmed transactions, ensuring no loss or duplication.
This way, checkpoint + confirmation queue + re-fetch logic together make the scanner resilient to downtime and chain reorganizations.

## What We Can Add for Downtime Recovery

### Current State:
- tryReconnect ensures the scanner reconnects to the Ethereum node automatically if the WebSocket or RPC connection fails.
- Checkpoint file keeps track of the last successfully processed block.

### What we can add:
- Catch-up Processing: On restart, fetch and process all blocks from lastProcessedBlock + 1 to the current node height minus confirmations. This prevents skipping transactions during downtime.
- Retry with Backoff: Implement exponential backoff for failed block fetches or RPC calls to handle transient node issues.
- Optional Persistent Buffer:
  - Maintain an in-memory queue to temporarily store blocks or transactions.
  - This helps to resume processing gracefully if the scanner crashes in the middle of a batch.

## What We Can Add for Retry Situations
### Current Implementation:
- We have tryReconnect to handle WebSocket or RPC connection failures with the Ethereum node.

### What can be extended:
- Exponential Backoff for Retries:
  - Instead of retrying immediately, introduce exponential backoff (e.g., 1s → 2s → 4s → 8s) when fetching blocks or submitting transactions.
  - This prevents overwhelming the node during transient failures.
- Retry on Temporary RPC Errors:
  - Catch errors like timeouts, rate limits, or network hiccups.
  - Retry the operation a configurable number of times before marking it as failed.
- Idempotent Processing:
  - Ensure that retrying the same block or transaction does not result in duplicates being published or logged.
  - Already partially handled via checkpointing and transaction filtering with the Bloom filter.

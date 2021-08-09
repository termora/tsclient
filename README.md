# tsclient

`tsclient` is an alternative Go client for [Typesense](https://typesense.org/).

**Why not [the official client](https://github.com/typesense/typesense-go)?** We didn't like it, that's all.
This client is a lot more similar in usage to Arikawa, the Discord library used in Termora.

**WARNING:** this is a *heavy* work in progress! The API may change at any time.

## Todo

- [X] Collection endpoints
- [ ] Document endpoints
  - [X] Inserting documents
  - [X] Upserting documents
  - [X] Retrieving documents
  - [X] Updating documents
  - [X] Deleting documents
  - [X] Bulk deleting documents
  - [X] Bulk importing documents
  - [ ] Search
  - [ ] Group by searches
- [ ] API key endpoints
- [ ] Override endpoints
- [ ] Synonym endpoints
- [ ] Cluster operations
  - [ ] Create snapshot
  - [ ] Re-elect leader
  - [ ] Toggle slow request log
  - [ ] Cluster metrics
  - [ ] API stats
  - [X] Health
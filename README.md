# tsclient

[![Go Reference](https://pkg.go.dev/badge/github.com/termora/tsclient.svg)](https://pkg.go.dev/github.com/termora/tsclient)

`tsclient` is an alternative Go client for [Typesense](https://typesense.org/).

**Why not [the official client](https://github.com/typesense/typesense-go)?** We didn't like it, that's all.
This client is a lot more similar in usage to Arikawa, the Discord library used in Termora.

> **Note:** this is still a work in progress! The current API is unlikely to change, but no promises about that.

## Differences to the official client

- No `interface{}`: all methods with variable responses (such as inserting/deleting documents) allow passing in a pointer to unmarshal to.
  - In the few places where the API returns variable responses in otherwise structured data, raw byte slices are used instead of the empty interface.
- Cleaner layout: no method chaining needed to make requests. For example, `client.Collection("companies").Document("123").Delete()` in the official client becomes `client.DeleteDocument("companies", "123")`.
- No generated code: the official client uses a lot of generated code, which, while perfectly passable, isn't very intuitive to use. Everything in this client is written by hand.

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
  - [X] Search
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

## License

BSD 3-Clause License

Copyright (c) 2021, Starshine System  
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice, this
  list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright notice,
  this list of conditions and the following disclaimer in the documentation
  and/or other materials provided with the distribution.

3. Neither the name of the copyright holder nor the names of its
  contributors may be used to endorse or promote products derived from
  this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.


<h1 align="center">
  <a href="libp2p.io"><img width="250" src="https://github.com/libp2p/libp2p/blob/master/logo/black-bg-2.png?raw=true" alt="libp2p hex logo" /></a>
</h1>

<h3 align="center">The Go implementation of the libp2p Networking Stack.</h3>

<p align="center">
  <a href="http://ipn.io"><img src="https://img.shields.io/badge/made%20by-Protocol%20Labs-blue.svg?style=flat-square" /></a>
  <a href="http://libp2p.io/"><img src="https://img.shields.io/badge/project-libp2p-blue.svg?style=flat-square" /></a>
  <a href="http://webchat.freenode.net/?channels=%23ipfs"><img src="https://img.shields.io/badge/freenode-%23ipfs-blue.svg?style=flat-square" /></a>
  <a href="https://waffle.io/libp2p/libp2p"><img src="https://img.shields.io/badge/pm-waffle-blue.svg?style=flat-square" /></a>
</p>

<p align="center">
  <a href="https://travis-ci.org/libp2p/go-libp2p"><img src="https://travis-ci.org/libp2p/go-libp2p.svg?branch=master" /></a>
  <!--<a href="https://circleci.com/gh/libp2p/go-libp2p"><img src="https://circleci.com/gh/libp2p/go-libp2p.svg?style=svg" /></a>-->
  <!--<a href="https://coveralls.io/github/libp2p/go-libp2p?branch=master"><img src="https://coveralls.io/repos/github/libp2p/go-libp2p/badge.svg?branch=master"></a>-->
  <br>
  <a href="https://github.com/RichardLitt/standard-readme"><img src="https://img.shields.io/badge/standard--readme-OK-green.svg?style=flat-square" /></a>
  <a href="https://godoc.org/github.com/libp2p/go-libp2p"><img src="https://godoc.org/github.com/ipfs/go-libp2p?status.svg" /></a>
  <a href=""><img src="https://img.shields.io/badge/golang-%3E%3D1.8.0-orange.svg?style=flat-square" /></a>
  <br>
</p>

# Project status

[![Throughput Graph](https://graphs.waffle.io/libp2p/go-libp2p/throughput.svg)](https://waffle.io/libp2p/go-libp2p/metrics/throughput)

[**`Weekly Core Dev Calls`**](https://github.com/ipfs/pm/issues/674)

# Table of Contents

- [Background](#background)
- [Bundles](#bundles)
- [Usage](#usage)
  - [Install](#install)
  - [API](#api)
  - [Examples](#examples)
- [Development](#development)
  - [Tests](#tests)
  - [Packages](#packages)
- [Contribute](#contribute)
- [License](#license)

## Background

[libp2p](https://github.com/libp2p/specs) is a networking stack and library modularized out of [The IPFS Project](https://github.com/ipfs/ipfs), and bundled separately for other tools to use.
>
libp2p is the product of a long, and arduous quest of understanding -- a deep dive into the internet's network stack, and plentiful peer-to-peer protocols from the past. Building large scale peer-to-peer systems has been complex and difficult in the last 15 years, and libp2p is a way to fix that. It is a "network stack" -- a protocol suite -- that cleanly separates concerns, and enables sophisticated applications to only use the protocols they absolutely need, without giving up interoperability and upgradeability. libp2p grew out of IPFS, but it is built so that lots of people can use it, for lots of different projects.
>
> We will be writing a set of docs, posts, tutorials, and talks to explain what p2p is, why it is tremendously useful, and how it can help your existing and new projects. But in the meantime, check out
>
> - [**The libp2p Specification**](https://github.com/libp2p/specs)
> - [**go-libp2p implementation**](https://github.com/libp2p/go-libp2p)
> - [**js-libp2p implementation**](https://github.com/libp2p/js-libp2p)


## Bundles

There is currently only one bundle of `go-libp2p`, this package. This bundle is used by [`go-ipfs`](https://github.com/ipfs/go-ipfs).

## Usage

`go-libp2p` repo is a place holder for the list of Go modules that compose Go libp2p, as well as its entry point.

### Install

```bash
> go get -u -d github.com/libp2p/go-libp2p/...
> cd $GOPATH/src/github.com/libp2p/go-libp2p
> make
> make deps
```

### API

[![GoDoc](https://godoc.org/github.com/ipfs/go-libp2p?status.svg)](https://godoc.org/github.com/libp2p/go-libp2p)

### Examples

Examples can be found in the [examples repo](https://github.com/libp2p/go-libp2p-examples).

## Development

### Dependencies

While developing, you need to use [gx to install and link your dependencies](https://github.com/whyrusleeping/gx#dependencies), to do that, run:

```sh
> make deps
```

Before commiting and pushing to Github, make sure to rewind the gx'ify of dependencies. You can do that with:

```sh
> make publish
```

### Tests

Running of individual tests is done through `gx test <path to test>`

```bash
$ cd $GOPATH/src/github.com/libp2p/go-libp2p
$ make deps
$ gx test ./p2p/<path of module you want to run tests for>
```

### Packages

> This table is generated using the module [`package-table`](https://github.com/ipfs-shipyard/package-table) with `package-table --data=package-list.json`.

List of packages currently in existence for libp2p:

| Name | CI/Travis | CI/Jenkins | Coverage |
| ---------|---------|---------|--------- |
| **Libp2p** |
| [`interface-libp2p`](//github.com/libp2p/interface-libp2p) |  |  |  |
| [`libp2p`](//github.com/libp2p/go-libp2p) | [![Travis CI](https://travis-ci.org/libp2p/go-libp2p.svg?branch=master)](https://travis-ci.org/libp2p/go-libp2p) | [![jenkins](https://ci.ipfs.team/buildStatus/icon?job=libp2p/go-libp2p/master)](https://ci.ipfs.team/job/libp2p/job/go-libp2p/job/master/) | [![codecov](https://codecov.io/gh/libp2p/go-libp2p/branch/master/graph/badge.svg)](https://codecov.io/gh/libp2p/go-libp2p) |
| **Connection** |
| [`interface-connection`](//github.com/libp2p/interface-connection) |  |  |  |
| [`go-libp2p-net`](//github.com/libp2p/go-libp2p-net) | [![Travis CI](https://travis-ci.org/libp2p/go-libp2p-net.svg?branch=master)](https://travis-ci.org/libp2p/go-libp2p-net) | [![jenkins](https://ci.ipfs.team/buildStatus/icon?job=libp2p/go-libp2p-net/master)](https://ci.ipfs.team/job/libp2p/job/go-libp2p-net/job/master/) | [![codecov](https://codecov.io/gh/libp2p/go-libp2p-net/branch/master/graph/badge.svg)](https://codecov.io/gh/libp2p/go-libp2p-net) |
| **Transport** |
| [`interface-transport`](//github.com/libp2p/interface-transport) |  |  |  |
| [`go-ws-transport`](//github.com/libp2p/go-ws-transport) | [![Travis CI](https://travis-ci.org/libp2p/go-ws-transport.svg?branch=master)](https://travis-ci.org/libp2p/go-ws-transport) | [![jenkins](https://ci.ipfs.team/buildStatus/icon?job=libp2p/go-ws-transport/master)](https://ci.ipfs.team/job/libp2p/job/go-ws-transport/job/master/) | [![codecov](https://codecov.io/gh/libp2p/go-ws-transport/branch/master/graph/badge.svg)](https://codecov.io/gh/libp2p/go-ws-transport) |
| [`go-libp2p-transport`](//github.com/libp2p/go-libp2p-transport) | [![Travis CI](https://travis-ci.org/libp2p/go-libp2p-transport.svg?branch=master)](https://travis-ci.org/libp2p/go-libp2p-transport) | [![jenkins](https://ci.ipfs.team/buildStatus/icon?job=libp2p/go-libp2p-transport/master)](https://ci.ipfs.team/job/libp2p/job/go-libp2p-transport/job/master/) | [![codecov](https://codecov.io/gh/libp2p/go-libp2p-transport/branch/master/graph/badge.svg)](https://codecov.io/gh/libp2p/go-libp2p-transport) |
| [`go-tcp-transport`](//github.com/libp2p/go-tcp-transport) | [![Travis CI](https://travis-ci.org/libp2p/go-tcp-transport.svg?branch=master)](https://travis-ci.org/libp2p/go-tcp-transport) | [![jenkins](https://ci.ipfs.team/buildStatus/icon?job=libp2p/go-tcp-transport/master)](https://ci.ipfs.team/job/libp2p/job/go-tcp-transport/job/master/) | [![codecov](https://codecov.io/gh/libp2p/go-tcp-transport/branch/master/graph/badge.svg)](https://codecov.io/gh/libp2p/go-tcp-transport) |
| [`go-libp2p-transport-upgrader`](//github.com/libp2p/go-libp2p-transport-upgrader) | [![Travis CI](https://travis-ci.org/libp2p/go-libp2p-transport-upgrader.svg?branch=master)](https://travis-ci.org/libp2p/go-libp2p-transport-upgrader) | [![jenkins](https://ci.ipfs.team/buildStatus/icon?job=libp2p/go-libp2p-transport-upgrader/master)](https://ci.ipfs.team/job/libp2p/job/go-libp2p-transport-upgrader/job/master/) | [![codecov](https://codecov.io/gh/libp2p/go-libp2p-transport-upgrader/branch/master/graph/badge.svg)](https://codecov.io/gh/libp2p/go-libp2p-transport-upgrader) |
| **Crypto Channels** |
| [`go-libp2p-secio`](//github.com/libp2p/go-libp2p-secio) | [![Travis CI](https://travis-ci.org/libp2p/go-libp2p-secio.svg?branch=master)](https://travis-ci.org/libp2p/go-libp2p-secio) | [![jenkins](https://ci.ipfs.team/buildStatus/icon?job=libp2p/go-libp2p-secio/master)](https://ci.ipfs.team/job/libp2p/job/go-libp2p-secio/job/master/) | [![codecov](https://codecov.io/gh/libp2p/go-libp2p-secio/branch/master/graph/badge.svg)](https://codecov.io/gh/libp2p/go-libp2p-secio) |
| **Stream Muxers** |
| [`go-stream-muxer`](//github.com/libp2p/go-stream-muxer) | [![Travis CI](https://travis-ci.org/libp2p/go-stream-muxer.svg?branch=master)](https://travis-ci.org/libp2p/go-stream-muxer) | [![jenkins](https://ci.ipfs.team/buildStatus/icon?job=libp2p/go-stream-muxer/master)](https://ci.ipfs.team/job/libp2p/job/go-stream-muxer/job/master/) | [![codecov](https://codecov.io/gh/libp2p/go-stream-muxer/branch/master/graph/badge.svg)](https://codecov.io/gh/libp2p/go-stream-muxer) |
| **Discovery** |
| [`interface-peer-discovery`](//github.com/libp2p/interface-peer-discovery) |  |  |  |
| **NAT Traversal** |
| [`go-libp2p-circuit`](//github.com/libp2p/go-libp2p-circuit) | [![Travis CI](https://travis-ci.org/libp2p/go-libp2p-circuit.svg?branch=master)](https://travis-ci.org/libp2p/go-libp2p-circuit) | [![jenkins](https://ci.ipfs.team/buildStatus/icon?job=libp2p/go-libp2p-circuit/master)](https://ci.ipfs.team/job/libp2p/job/go-libp2p-circuit/job/master/) | [![codecov](https://codecov.io/gh/libp2p/go-libp2p-circuit/branch/master/graph/badge.svg)](https://codecov.io/gh/libp2p/go-libp2p-circuit) |
| [`go-libp2p-nat`](//github.com/libp2p/go-libp2p-nat) | [![Travis CI](https://travis-ci.org/libp2p/go-libp2p-nat.svg?branch=master)](https://travis-ci.org/libp2p/go-libp2p-nat) | [![jenkins](https://ci.ipfs.team/buildStatus/icon?job=libp2p/go-libp2p-nat/master)](https://ci.ipfs.team/job/libp2p/job/go-libp2p-nat/job/master/) | [![codecov](https://codecov.io/gh/libp2p/go-libp2p-nat/branch/master/graph/badge.svg)](https://codecov.io/gh/libp2p/go-libp2p-nat) |
| **Data Types** |
| [`go-libp2p-peer`](//github.com/libp2p/go-libp2p-peer) | [![Travis CI](https://travis-ci.org/libp2p/go-libp2p-peer.svg?branch=master)](https://travis-ci.org/libp2p/go-libp2p-peer) | [![jenkins](https://ci.ipfs.team/buildStatus/icon?job=libp2p/go-libp2p-peer/master)](https://ci.ipfs.team/job/libp2p/job/go-libp2p-peer/job/master/) | [![codecov](https://codecov.io/gh/libp2p/go-libp2p-peer/branch/master/graph/badge.svg)](https://codecov.io/gh/libp2p/go-libp2p-peer) |
| [`go-libp2p-peerstore`](//github.com/libp2p/go-libp2p-peerstore) | [![Travis CI](https://travis-ci.org/libp2p/go-libp2p-peerstore.svg?branch=master)](https://travis-ci.org/libp2p/go-libp2p-peerstore) | [![jenkins](https://ci.ipfs.team/buildStatus/icon?job=libp2p/go-libp2p-peerstore/master)](https://ci.ipfs.team/job/libp2p/job/go-libp2p-peerstore/job/master/) | [![codecov](https://codecov.io/gh/libp2p/go-libp2p-peerstore/branch/master/graph/badge.svg)](https://codecov.io/gh/libp2p/go-libp2p-peerstore) |
| [`go-libp2p-protocol`](//github.com/libp2p/go-libp2p-protocol) | [![Travis CI](https://travis-ci.org/libp2p/go-libp2p-protocol.svg?branch=master)](https://travis-ci.org/libp2p/go-libp2p-protocol) | [![jenkins](https://ci.ipfs.team/buildStatus/icon?job=libp2p/go-libp2p-protocol/master)](https://ci.ipfs.team/job/libp2p/job/go-libp2p-protocol/job/master/) | [![codecov](https://codecov.io/gh/libp2p/go-libp2p-protocol/branch/master/graph/badge.svg)](https://codecov.io/gh/libp2p/go-libp2p-protocol) |
| **Content Routing** |
| [`interface-content-routing`](//github.com/libp2p/interface-content-routing) |  |  |  |
| **Peer Routing** |
| [`interface-peer-routing`](//github.com/libp2p/interface-peer-routing) |  |  |  |
| **Record Store** |
| [`interface-record-store`](//github.com/libp2p/interface-record-store) |  |  |  |
| **Miscellaneous** |
| [`go-libp2p-crypto`](//github.com/libp2p/go-libp2p-crypto) | [![Travis CI](https://travis-ci.org/libp2p/go-libp2p-crypto.svg?branch=master)](https://travis-ci.org/libp2p/go-libp2p-crypto) | [![jenkins](https://ci.ipfs.team/buildStatus/icon?job=libp2p/go-libp2p-crypto/master)](https://ci.ipfs.team/job/libp2p/job/go-libp2p-crypto/job/master/) | [![codecov](https://codecov.io/gh/libp2p/go-libp2p-crypto/branch/master/graph/badge.svg)](https://codecov.io/gh/libp2p/go-libp2p-crypto) |
| [`go-libp2p-interface-connmgr`](//github.com/libp2p/go-libp2p-interface-connmgr) | [![Travis CI](https://travis-ci.org/libp2p/go-libp2p-interface-connmgr.svg?branch=master)](https://travis-ci.org/libp2p/go-libp2p-interface-connmgr) | [![jenkins](https://ci.ipfs.team/buildStatus/icon?job=libp2p/go-libp2p-interface-connmgr/master)](https://ci.ipfs.team/job/libp2p/job/go-libp2p-interface-connmgr/job/master/) | [![codecov](https://codecov.io/gh/libp2p/go-libp2p-interface-connmgr/branch/master/graph/badge.svg)](https://codecov.io/gh/libp2p/go-libp2p-interface-connmgr) |
| [`go-libp2p-swarm`](//github.com/libp2p/go-libp2p-swarm) | [![Travis CI](https://travis-ci.org/libp2p/go-libp2p-swarm.svg?branch=master)](https://travis-ci.org/libp2p/go-libp2p-swarm) | [![jenkins](https://ci.ipfs.team/buildStatus/icon?job=libp2p/go-libp2p-swarm/master)](https://ci.ipfs.team/job/libp2p/job/go-libp2p-swarm/job/master/) | [![codecov](https://codecov.io/gh/libp2p/go-libp2p-swarm/branch/master/graph/badge.svg)](https://codecov.io/gh/libp2p/go-libp2p-swarm) |
| [`go-libp2p-host`](//github.com/libp2p/go-libp2p-host) | [![Travis CI](https://travis-ci.org/libp2p/go-libp2p-host.svg?branch=master)](https://travis-ci.org/libp2p/go-libp2p-host) | [![jenkins](https://ci.ipfs.team/buildStatus/icon?job=libp2p/go-libp2p-host/master)](https://ci.ipfs.team/job/libp2p/job/go-libp2p-host/job/master/) | [![codecov](https://codecov.io/gh/libp2p/go-libp2p-host/branch/master/graph/badge.svg)](https://codecov.io/gh/libp2p/go-libp2p-host) |
| [`go-libp2p-blankhost`](//github.com/libp2p/go-libp2p-blankhost) | [![Travis CI](https://travis-ci.org/libp2p/go-libp2p-blankhost.svg?branch=master)](https://travis-ci.org/libp2p/go-libp2p-blankhost) | [![jenkins](https://ci.ipfs.team/buildStatus/icon?job=libp2p/go-libp2p-blankhost/master)](https://ci.ipfs.team/job/libp2p/job/go-libp2p-blankhost/job/master/) | [![codecov](https://codecov.io/gh/libp2p/go-libp2p-blankhost/branch/master/graph/badge.svg)](https://codecov.io/gh/libp2p/go-libp2p-blankhost) |
| [`go-conn-security-multistream`](//github.com/libp2p/go-conn-security-multistream) | [![Travis CI](https://travis-ci.org/libp2p/go-conn-security-multistream.svg?branch=master)](https://travis-ci.org/libp2p/go-conn-security-multistream) | [![jenkins](https://ci.ipfs.team/buildStatus/icon?job=libp2p/go-conn-security-multistream/master)](https://ci.ipfs.team/job/libp2p/job/go-conn-security-multistream/job/master/) | [![codecov](https://codecov.io/gh/libp2p/go-conn-security-multistream/branch/master/graph/badge.svg)](https://codecov.io/gh/libp2p/go-conn-security-multistream) |
| [`go-conn-security`](//github.com/libp2p/go-conn-security) | [![Travis CI](https://travis-ci.org/libp2p/go-conn-security.svg?branch=master)](https://travis-ci.org/libp2p/go-conn-security) | [![jenkins](https://ci.ipfs.team/buildStatus/icon?job=libp2p/go-conn-security/master)](https://ci.ipfs.team/job/libp2p/job/go-conn-security/job/master/) | [![codecov](https://codecov.io/gh/libp2p/go-conn-security/branch/master/graph/badge.svg)](https://codecov.io/gh/libp2p/go-conn-security) |
| [`go-libp2p-interface-pnet`](//github.com/libp2p/go-libp2p-interface-pnet) | [![Travis CI](https://travis-ci.org/libp2p/go-libp2p-interface-pnet.svg?branch=master)](https://travis-ci.org/libp2p/go-libp2p-interface-pnet) | [![jenkins](https://ci.ipfs.team/buildStatus/icon?job=libp2p/go-libp2p-interface-pnet/master)](https://ci.ipfs.team/job/libp2p/job/go-libp2p-interface-pnet/job/master/) | [![codecov](https://codecov.io/gh/libp2p/go-libp2p-interface-pnet/branch/master/graph/badge.svg)](https://codecov.io/gh/libp2p/go-libp2p-interface-pnet) |
| [`go-libp2p-metrics`](//github.com/libp2p/go-libp2p-metrics) | [![Travis CI](https://travis-ci.org/libp2p/go-libp2p-metrics.svg?branch=master)](https://travis-ci.org/libp2p/go-libp2p-metrics) | [![jenkins](https://ci.ipfs.team/buildStatus/icon?job=libp2p/go-libp2p-metrics/master)](https://ci.ipfs.team/job/libp2p/job/go-libp2p-metrics/job/master/) | [![codecov](https://codecov.io/gh/libp2p/go-libp2p-metrics/branch/master/graph/badge.svg)](https://codecov.io/gh/libp2p/go-libp2p-metrics) |
| **Utilities** |
| [`go-libp2p-loggables`](//github.com/libp2p/go-libp2p-loggables) | [![Travis CI](https://travis-ci.org/libp2p/go-libp2p-loggables.svg?branch=master)](https://travis-ci.org/libp2p/go-libp2p-loggables) | [![jenkins](https://ci.ipfs.team/buildStatus/icon?job=libp2p/go-libp2p-loggables/master)](https://ci.ipfs.team/job/libp2p/job/go-libp2p-loggables/job/master/) | [![codecov](https://codecov.io/gh/libp2p/go-libp2p-loggables/branch/master/graph/badge.svg)](https://codecov.io/gh/libp2p/go-libp2p-loggables) |
| [`go-maddr-filter`](//github.com/libp2p/go-maddr-filter) | [![Travis CI](https://travis-ci.org/libp2p/go-maddr-filter.svg?branch=master)](https://travis-ci.org/libp2p/go-maddr-filter) | [![jenkins](https://ci.ipfs.team/buildStatus/icon?job=libp2p/go-maddr-filter/master)](https://ci.ipfs.team/job/libp2p/job/go-maddr-filter/job/master/) | [![codecov](https://codecov.io/gh/libp2p/go-maddr-filter/branch/master/graph/badge.svg)](https://codecov.io/gh/libp2p/go-maddr-filter) |
| [`go-libp2p-netutil`](//github.com/libp2p/go-libp2p-netutil) | [![Travis CI](https://travis-ci.org/libp2p/go-libp2p-netutil.svg?branch=master)](https://travis-ci.org/libp2p/go-libp2p-netutil) | [![jenkins](https://ci.ipfs.team/buildStatus/icon?job=libp2p/go-libp2p-netutil/master)](https://ci.ipfs.team/job/libp2p/job/go-libp2p-netutil/job/master/) | [![codecov](https://codecov.io/gh/libp2p/go-libp2p-netutil/branch/master/graph/badge.svg)](https://codecov.io/gh/libp2p/go-libp2p-netutil) |
| [`go-testutil`](//github.com/libp2p/go-testutil) | [![Travis CI](https://travis-ci.org/libp2p/go-testutil.svg?branch=master)](https://travis-ci.org/libp2p/go-testutil) | [![jenkins](https://ci.ipfs.team/buildStatus/icon?job=libp2p/go-testutil/master)](https://ci.ipfs.team/job/libp2p/job/go-testutil/job/master/) | [![codecov](https://codecov.io/gh/libp2p/go-testutil/branch/master/graph/badge.svg)](https://codecov.io/gh/libp2p/go-testutil) |

# Contribute

go-libp2p is part of [The IPFS Project](https://github.com/ipfs/ipfs), and is MIT licensed open source software. We welcome contributions big and small! Take a look at the [community contributing notes](https://github.com/ipfs/community/blob/master/contributing.md). Please make sure to check the [issues](https://github.com/ipfs/go-libp2p/issues). Search the closed ones before reporting things, and help us with the open ones.

Guidelines:

- read the [libp2p spec](https://github.com/libp2p/specs)
- please make branches + pull-request, even if working on the main repository
- ask questions or talk about things in [Issues](https://github.com/libp2p/go-libp2p/issues) or #ipfs on freenode.
- ensure you are able to contribute (no legal issues please-- we use the DCO)
- run `go fmt` before pushing any code
- run `golint` and `go vet` too -- some things (like protobuf files) are expected to fail.
- get in touch with @jbenet and @diasdavid about how best to contribute
- have fun!

There's a few things you can do right now to help out:
 - Go through the modules below and **check out existing issues**. This would be especially useful for modules in active development. Some knowledge of IPFS/libp2p may be required, as well as the infrasture behind it - for instance, you may need to read up on p2p and more complex operations like muxing to be able to help technically.
 - **Perform code reviews**.
 - **Add tests**. There can never be enough tests.

## Modularizing go-libp2p

We have currently a work in progress of modularizing go-libp2p from a repo monolith to several packages in different repos that can be reused for other projects of swapped for custom libp2p builds.

We want to maintain history, so we'll use git-subtree for extracting packages. Find instructions below:

```sh
# 1) create the extracted tree (has the directory specified as -P as its root)
> cd go-libp2p/
> git subtree split -P p2p/crypto/secio/ -b libp2p-secio
62b0a5c21574bcbe06c422785cd5feff378ae5bd
# important to delete the tree now, so that outdated imports fail in step 5
> git rm -r p2p/crypto/secio/
> git commit
> cd ../

# 2) make the new repo
> mkdir go-libp2p-secio
> cd go-libp2p-secio/
> git init && git commit --allow-empty

# 3) fetch the extracted tree from the previous repo
> git remote add libp2p ../go-libp2p
> git fetch libp2p
> git reset --hard libp2p/libp2p-secio

# 4) update self import paths
> sed -someflagsidontknow 'go-libp2p/p2p/crypto/secio' 'golibp2p-secio'
> git commit

# 5) create package.json and check all imports are correct
> vim package.json
> gx --verbose install --global
> gx-go rewrite
> go test ./...
> gx-go rewrite --undo
> git commit

# 4) make the package ready
> vim README.md LICENSE
> git commit

# 5) bump the version separately
> vim package.json
> gx publish
> git add package.json .gx/
> git commit -m 'Publish 1.2.3'

# 6) clean up and push
> git remote rm libp2p
> git push origin master
```

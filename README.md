# ipfsgit

Tools to simplify publishing Git repositories on IPFS using IPNS including
self-hosted pinning

## In brief:

For some reason this isn't already done, even though it is simple enough to
follow the manual to do it manually, this project is simply automation with
some small refinements to make it easy to publish your software via IPFS,
and to help you ensure it's available even when your device is offline by
keeping the current version updated via IPNS.

With the resurgence of government attempting to stifle speech once again, in
an echo of the PGP and RSA encryption battle, with the outlawing of Tornado
Cash, it is only a matter of time before governments attempt to attack both
Bitcoin developers and the hosting of the code that runs it and especially
code which protects privacy.

So, this is our multi-pronged solution:

- simplify tagging and committing changes to a Git repository onto an IPFS IPNS
  based filesystem root to fully decentralise hosting of software
- simplify pinning of these repositories on budget storage VPS services (in
  foreign and perhaps multiple jurisdictions) to prevent censorship of the
  software
- provide simple tools that enable easy setup and configuration of an IPFS
  node and configuring it in a consistent way that enables interoperability
  with vanilla Golang configurations as well as Git in general.

IPFS is essentially a generalised, programmable variant of the Bittorrent
protocol, which has become the primary method of peer to peer file sharing.
It identifies content by computed hashes similar to Magnet links, and
provides a web proxy that can be used to access this data via any HTTP aware
application. Its main distinguishing feature is the P2P network which 
enables seacrhing peers for uncached content in order to cache it.

Unfortunately, due to the nature of the culture of the individuals and
projects that currently make the most use of IPFS, ie, NFT and ethereum
shitcoins, it seems that actual privacy protecting and censorship prevention
uses of the protocol have been neglected, and information on how to apply
the technology this way are hard to find. Instead, only how to publish
stupid JPEGs of degenerate looking monkeys. Pathetic.

This project will aim to fix this situation, and will itself be hosted on
IPFS not long after it is fully functional.

A vanity address generator, which enables you to create IPNS keys with
public keys that contain a fragment of a recognisable word, will also be
added. Maybe in the future if you own a Bitcoin mining rig an interface to
accelerate minting these nice names will be added.

## Installation

If you don't have Go installed, this is the first thing to do. It's the only
tool you need. You can install the `golang-go` package in Ubuntu 22, if you
prefer, this is a good option. Similarly for Mac users it is possible to do
with homebrew, but the developer of this project hates apple and would never
buy an overpriced mac so donations are welcome to provide him but otherwise
DYOR and GFY.

Or you can do this, real simple

    cd
    mkdir bin
    sudo apt install -y wget curl git build-essential
    wget https://go.dev/dl/go1.18.linux-amd64.tar.gz
    tar xvf go1.18.linux-amd64.tar.gz

Then add this to the end of your `~/.bashrc`:

    export GOBIN=$HOME/bin
    export GOPATH=$HOME
    export GOROOT=$GOPATH/go
    export PATH=$HOME/go/bin:$HOME/.local/bin:$GOBIN:$PATH

And voila, you have Go 1.18.5 which is the recommended version to run for
the time being. 1.19 versions may give you issues in general, it is
recommended to just go with the 1.18.5 for now. In the future likely you can
just install the latest but this is the stable version most supported by
most projects.

With these prerequisites out of the way, now you can install IPFS:

    go install -v -x -a github.com/ipfs/kubo/cmd/ipfs@latest

You are free to choose to use some other language version but kubo (formerly
go-ipfs) is the original and everything in the last screenful of text is all
you need to get it, so why bother with any other garbage "expressive" language.

Last step is to set up ipfs as a service in your own user account privileges
using your home folder to locate the cache storage.

We are assuming Ubuntu 22 or derivative such as Pop OS 22, or the mother
distribution Debian Sid. With the state of especially "hybrid graphics"
laptops these days, which are becoming the de facto standard, Pop OS 22 is
the best choice by far. After long battles in the trenches we will just
laugh at anyone trying to do it with other linux distributions. If you don't
use something based on Debian Sid you already know what you are doing anyway
smart pants, gfy.

To create a systemd service to automatically start up your own personal IPFS
node and not bother with this smarty pants "IPFS desktop", which you are
also free to use instead, but it's not gonna run auto for you, which kinda
defeats the purpose of hosting your own stuff unless you also add it to your
auto-start apps for login. But this option I will show you doesn't need you
to log in, it will start up by itself and already be running when the login
streen appears.

Put this following content into `/etc/systemd/system/ipfs.service`:

    [Unit]
    Description=InterPlanetary File System (IPFS) daemon
    Documentation=https://docs.ipfs.tech/
    After=network.target
    
    [Service]
    
    # enable for 1-1024 port listening
    AmbientCapabilities=CAP_NET_BIND_SERVICE 
    # enable to specify a custom path see docs/environment-variables.md for further documentations
    #Environment=IPFS_PATH=/custom/ipfs/path
    # enable to specify a higher limit for open files/connections
    #LimitNOFILE=1000000
    
    #don't use swap
    MemorySwapMax=0
    
    # Don't timeout on startup. Opening the IPFS repo can take a long time in some cases (e.g., when
    # badger is recovering) and migrations can delay startup.
    #
    # Ideally, we'd be a bit smarter about this but there's no good way to do that without hooking
    # systemd dependencies deeper into go-ipfs.
    TimeoutStartSec=infinity
    
    Type=notify
    User=loki
    Group=loki
    StateDirectory=ipfs
    Environment=IPFS_PATH="/home/loki/.ipfs"
    ExecStart=/home/loki/bin/ipfs daemon --init --migrate
    Restart=on-failure
    KillSignal=SIGINT
    
    [Install]
    WantedBy=default.target

Replace every instance of the word `loki` in here with your username.

Once that is in place, fire it up:

    sudo systemctl enable ipfs
    sudo systemctl start ipfs

If you are running Brave browser, you want to have the IPFS Companion
enabled, set your "ipfs public gateway address" to `http://localhost:8080`
and set as your default method to resolve IPFS resources as "gateway". By
default this will be a public gateway, but you don't need that anymore, you
have your own node now.

## Pull

Pull is the name for the IPFS proxy. So named because it causes data to be 
downloaded from IPFS p2p network.

The application in [./cmd/pull](cmd/ipfsproxy/) is a HTTP proxy that you 
can set as your system proxy, which will redirect the 'ipfs', 'ipns' and 
'ipld' "hostnames" and point them at a default configured IPFS installation, 
such as what the foregoing section explained how to install and set to run 
as a service.

## Push




#### TODO: finish the rest once ipfsgit tools are written
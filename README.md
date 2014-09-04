partyline
=========

Partyline is an experimental layer 2 ethernet bridge.

**PARTYLINE IS NOT SECURE**


Description
-----------

Partyline creates a 90's style ethernet hub over UDP using VXLAN encapsulation.

Installation
------------

    go get -u -v github.com/davecheney/partyline/...

Usage
----

Partyline requires a `partylined` server to be reachable by all partyline clients. The server process is non priviledged and communicates with partyline clients on the UDP port you nominate. Start one like so

    partylined -e $YOURHOST:$PORT

For testing, a `partylined` can be started on your localhost.

    partylined -e 127.0.0.100:9001

Once the server is running it will spew information as frames pass through it.

Create some partyline clients, these do need the permission to talk to the `/dev/tun` device on your machine

    [sudo] partyline -e $YOURPARTYLINEDHOST:$PORT

This will report the name of the `tap` device assigned to this partyline client. `tap` devices are layer 2 ethernet devices. Tap devices start in down state, so you will need to bring the interface up before you can do anything with it

    [sudo] ip link set dev tap0 up

At this point your operating system will spring into life and start to send some packets as it sniffs around.

Use cases
---------

Having another layer 2 interface on your machine is interesting, sort of, you can talk to other partyline participants if you can agree on ip details on your tap devices.

A more interesting use case would be to configure a bridge interface and bridge your partyline tap device onto your local ethernet segment.





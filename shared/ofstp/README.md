# Object File State Transfer Protocol

Transfer object files with state over TCP. This module is a helper for encoding
and decoing `OFSTP` packets.

## Stateless packets

Stateless packets do not contain the state of the `stack` or the `process flag`.
More expicitly they only transfer the `.data` and `.text` sections. This packet
type is only used in the initial transfer from client to router to server. This
packet will then be loaded into the virtual machine on the server and all
subsequent packets will either be [Statefull packets](#statefull-packets) or
they will be [Return packets](#return-packets).

These packets will be 192 bytes (1 + 16 + 175) long. 1 byte for the packet type,
16 bytes for the `.data` state, and 175 bytes for the `.text` state.

```
0000 0001 # packet header / packet type
16 Bytes  # .data section
.
.
175 Bytes # .text section
.
.
```

## Statefull packets

Statefull packets will contain all four registers in order
(`R0`, `R1`, `SP`, `PC`) followed by the entire memory of the virtual machine.
The actual information from the packet is up to the receiver to derive. An
example would be checking the upper 5 bits of `Mem(79)` to see if the process
is in a halted state or a sleeping state (these are the only valid states for
a process to be sent from the server back to the router).

```
0000 0010 # packet header / packet type
1 Byte    # R0 State
1 Byte    # R1 State
1 Byte    # SP State
1 Byte    # PC State
16 Bytes  # .data section State
64 Bytes  # Stack State
1 Byte    # Process Flag State
175 Bytes # .text section
```

Notice: the `.text` section is not statefull, this is intentional. The reason
for keeping this data in the statefull packet is such that if a process is not
assigned to its preferred machine, it will need to reconstruct it's instructions
in the memory of the new virtual machine. ***More information about preferred
machines in the router documentation***

## Return packets

Return packets are created from [Statefull packets](#statefull-packets). They
contain the output from a process and its exit code OR if it has errored, it
will have the reason the process has had to exit and its exit status. Packet
size is within the bounds of 2 Bytes to 1500 Bytes

```
0000 0011  # packet header / packet type
1 Byte     # exit status of the program (we know it is 0-255 so 1 Byte is fine)
1498 Bytes # program output: (See note below) TCP_MAX_SIZE - header - exit status
```

Notice: This stackoverflow page is where I derived 1498 Bytes for the max packet
size: [Totally credible Max TCP Packet](https://stackoverflow.com/a/2614188).

#  Tcp chat
I decided to learn more about golang and tcp. So I wrote the application for a tcp server that broadcasts messages and user connecting/disconnecting notifications. The end of the stream is marked by a new line byte while the rest is utf8 bytes. There is no dedicated client. For now, GNU netcat can be used as a client.

Later I might expand this application by adding:
- a dedicated cli client
- limit message size and cap concurrent user count
- authentication and/or tls handshake

# hexagen

Hexagen generates coloured triangles in a hexagon from a value looking like a scuttlebot id (both public keys and hashes).

```
git clone ssb://%1jSkm2ziiZ9FbO/kyhRyd3gtn9UtbJQLzYf13HgRO4E=.sha256 hexagen
cd hexagen
go build
./hexagen $(sbot whoami | grep -o @.*\.ed25519sbot whoami | grep -o @.*\.ed25519)
```

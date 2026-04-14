# gpg-bridge
A bridge connects OpenSSH and GnuPG on Windows.

## Build

Build the Go executable on Windows:

```powershell
go build -o gpg-bridge.exe .
```

## Usage

1. Make sure you have set up gpg agent forwarding following [the guide](https://wiki.gnupg.org/AgentForwarding).

2. Directly using the socket provided by GnuPG will not work on Windows, so change the local socket to a TCP port instead.

   ```sshconfig
   RemoteForward <socket_on_remote_box> 127.0.0.1:4321
   ```

   You are free to use any port that has not been taken. `4321` is just an example.

3. Build a bridge between the TCP port and the GnuPG extra socket.

   ```powershell
   .\gpg-bridge.exe --extra 127.0.0.1:4321
   ```

   If you have customized the extra socket location, you can set the path using `--extra-socket`.

4. To run the SSH agent bridge, ensure `enable-putty-support` is configured for gpg-agent.

   ```text
   enable-putty-support
   ```

   Then start the SSH bridge:

   ```powershell
   .\gpg-bridge.exe --ssh \\.\pipe\gpg-bridge-ssh
   ```

   This can also be used with the extra socket bridge at the same time.

   ```powershell
   .\gpg-bridge.exe --extra 127.0.0.1:4321 --ssh \\.\pipe\gpg-bridge-ssh
   ```

5. Let OpenSSH use gpg-agent by setting `SSH_AUTH_SOCK` to `\\.\pipe\gpg-bridge-ssh`.

## Notes

There are several gotchas if you try to forward the gpg-agent socket on Windows without a bridge. See PowerShell/Win32-OpenSSH#1564.

1. Specifying the remote forward local socket path in OpenSSH can be tricky.
2. OpenSSH on Windows still does not handle Unix domain socket forwarding for this GnuPG scenario cleanly.
3. GnuPG on Windows simulates a Unix domain socket via a TCP stream socket with a custom connect step, so an adapter is required.

# gpg-bridge
转发 Windows 上 GnuPG 的到TCP网络的桥接工具。

## 构建

在 Windows 上构建 Go 可执行文件：

```powershell
go build -o gpg-bridge.exe .
```

## 使用方式

1. 请先按照[官方指南](https://wiki.gnupg.org/AgentForwarding)完成 gpg-agent 转发配置。

2. 在 Windows 上不能直接使用 GnuPG 提供的 socket，因此需要把本地 socket 改为 TCP 端口转发。

   ```sshconfig
   RemoteForward <socket_on_remote_box> 127.0.0.1:4321
   ```

   你可以使用任意未被占用的端口，`4321` 只是示例。

3. 在 TCP 端口与 GnuPG extra socket 之间启动桥接。

   ```powershell
   .\gpg-bridge.exe --extra 127.0.0.1:4321
   ```

   如果你自定义了 extra socket 的位置，可以通过 `--extra-socket` 指定路径。

4. 如果要启用 SSH agent 桥接，请先确认 `gpg-agent` 已开启 `enable-putty-support`。

   ```text
   enable-putty-support
   ```

   然后启动 SSH 桥接：

   ```powershell
   .\gpg-bridge.exe --ssh \\.\pipe\gpg-bridge-ssh
   ```

   它也可以和 extra socket 桥接同时使用。

   ```powershell
   .\gpg-bridge.exe --extra 127.0.0.1:4321 --ssh \\.\pipe\gpg-bridge-ssh
   ```

5. 将 `SSH_AUTH_SOCK` 设置为 `\\.\pipe\gpg-bridge-ssh`，让 OpenSSH 通过 gpg-agent 提供 SSH agent 能力。

## 说明

如果你尝试在 Windows 上不借助桥接直接转发 gpg-agent socket，会遇到一些已知问题，可参考 PowerShell/Win32-OpenSSH#1564。

1. 在 OpenSSH 中配置远端转发的本地 socket 路径比较麻烦。
2. 在这个 GnuPG 场景下，Windows 上的 OpenSSH 仍然无法很好地处理 Unix domain socket 转发。
3. Windows 上的 GnuPG 是通过带有自定义连接步骤的 TCP 流 socket 来模拟 Unix domain socket，因此需要额外的适配层。

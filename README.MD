# WT-SSH-OPEN
A simple stand alone helper program writen in GO to parse SSH URI's and open an SSH session in [WindowsTerminal](https://aka.ms/terminal) using openSSH.

ToDo:
* use `golang.org/x/sys/windows/registry` to create a user friendly install/uninstall flag
* Add additional OpenSSH flags to C-Params Parsing

URI based off of: https://datatracker.ietf.org/doc/html/draft-ietf-secsh-scp-sftp-ssh-uri

```shell
Usage: wtssh.exe [-v] ssh://[<user>[;ConnectTimeout=<Timeout>]@]<host>[:<port>]
```

Connection parameters are mapped to the `-o` options for OpenSSH see the [OpenSSH manpage](https://manpages.org/ssh) for further details.
**Note: Invalid options are ignored and Option variables are not validated** 

## Install as SSH protocol handler

You can use the following install.reg to add wtssh.exe as the ssh handler

```
Windows Registry Editor Version 5.00

[HKEY_CLASSES_ROOT\ssh]
@="URL:SSH Protocol"
"EditFlags"=dword:00210000
"URL Protocol"=""

; Change to wtssh location
[HKEY_CLASSES_ROOT\wtssh\shell\open\command]
@="\"C:\\tools\\wtssh\\wtssh.exe\" \"%1\""

[HKEY_LOCAL_MACHINE\SOFTWARE\RegisteredApplications]
"WT SSH handler"="Software\\Classes\\wtssh\\Capabilities"

[HKEY_CLASSES_ROOT\wtssh\Capabilities\UrlAssociations]
"ssh"="wtssh"
```

### Uninstall


```
Windows Registry Editor Version 5.00

[HKEY_CLASSES_ROOT\ssh]
@="URL:SSH Protocol"
"EditFlags"=dword:00210000
"URL Protocol"=""

; Change to wtssh location
[HKEY_CLASSES_ROOT\wtssh\shell\open\command]
@="\"C:\\scripts\\wtssh.exe\" \"%1\""

[HKEY_LOCAL_MACHINE\SOFTWARE\RegisteredApplications]
"WT SSH handler"="Software\\Classes\\wtssh\\Capabilities"

[HKEY_CLASSES_ROOT\wtssh\Capabilities\UrlAssociations]
"ssh"="wtssh"
```

# Build

`go build -o wtssh.exe -buildvcs=false -ldflags -H=windowsgui`
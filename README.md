## Nextwrk for i3

Nextwrk is a simple script that moves the current window to the next available workspace in i3. It is useful for quickly moving windows to a new workspace without needing to look at the workspaces in use.

### Usage:

1. Download from releases
2. Make the script executable: `chmod +x nextwrk`
3. Edit i3 config: `bindsym $mod+Shift+n exec /path/to/nextwrk`
4. Reload i3: `$mod+Shift+r`


---

```
Usage: nextwrk [--switch] [--renumber]
	[no args]		Move focused container to the next free workspace
	--switch		Move container to next free workspace and switch to it
	--renumber		Renumber all workspaces to remove gaps
	[any other args]	Show this help message
------------------------------
Built with Go for i3WM. MIT License. github.com/arrayindex-dev/nextwrk
```
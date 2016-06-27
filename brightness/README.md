# Brightness

In all likelihood, it'll make no sense for you to use this tool, and you'll
want to use `xbacklight` or some other reasonable tool. I was running into
strange issues with xbacklight where my screen would flash and everything would
freeze for a few seconds. That kind of terrified me, but I noticed editing
`/sys/class/backlight/intel_backlight` manually worked fine for me, so that's
what this program does.

# Game Boy Zero config menu

This is a simple config menu to be used at [Game Boy Zero projects](http://www.sudomod.com/wiki/index.php?title=Game_Boy_Zero)

It's a Work In Progress (still only at visual part =p)

## Sample menu definition (yaml file)
```
---
name: "Yaya's power optimization"

options:
  disablesamba:
    desc: "Disable Samba"

  disableromservice:
    desc: "Disable ROMService (for USB keys)"

  disabletv:
    desc:    "Disable TV and HDMI signals"
    cmd:     "echo '/usr/bin/tvservice -o' >> /etc/rc.local && /usr/bin/tvservice -o"
    undocmd: "sed -in -e '/tvservice -o/d' /etc/rc.local"
    check:   "grep -iq '/usr/bin/tvservice -o' /etc/rc.local"

  disablepileds:
    desc:    "Disable Pi status LEDs"
    cmd:     "echo -e '# Disable pi status LED\ndtparam=act_led_trigger=none\ndtparam=act_led_activelow=on' >> /boot/config.txt"
    undocmd: "sed -in -e '/Disable pi status LED/d;/dtparam=act_led_trigger=none/d;/dtparam=act_led_activelow=on/d' /boot/config.txt"
    check:   "grep -iq 'dtparam=act_led_trigger=none' /boot/config.txt"

  improvesdperf:
    desc:    "Improve SD performance"
    cmd:     "sed -in -e 's/elevator=deadline/elevator=noop/' /boot/cmdline.txt"
    undocmd: "sed -in -e 's/elevator=noop/elevator=deadline/' /boot/cmdline.txt"
    check:   "grep -iq 'elevator=deadline'"

  adjustswap:
    desc:    "Adjust swap"
    cmd:     "echo 'vm.swappiness = 1' >> /etc/sysctl.conf"
    undocmd: "sed -in -e '/vm.swappiness = 1/d' /etc/sysctl.conf"
    check:   "grep -iq 'vm.swappiness = 1' /etc/sysctl.conf"
    
  tweakfs:
    desc:    "Tweak filesystem for fewer writes"
    cmd:     "sed -i -re 's/(.*boot.*) defaults (.*)/\\1 defaults,noatime \\2/' /etc/fstab"
    undocmd: "sed -i -re 's/(.*boot.*) defaults,noatime (.*)/\\1 defaults \\2/' /etc/fstab"
    check:   "grep -q 'boot.* defaults,noatime' /etc/fstab"

  disableavahi:
    desc:    "Disable avahi service"
    cmd:     "apt-get purge avahi-daemon -y"
    undocmd: "apt-get install avahi-daemon -y"
    check:   "! dpkg-query -s avahi-daemon"

  disablewifi:
    desc:    "Disable WiFi"
    cmd:     "echo -e '# Disable wifi\ndtoverlay=pi3-disable-wifi' >> /boot/config.txt && ifdown wlan0"
    undocmd: "sed -i '/# Disable wifi/d;/dtoverlay=pi3-disable-wifi/d' /boot/config.txt && ifup wlan0"
    check:   "grep -q 'dtoverlay=pi3-disable-wifi' /boot/config.txt"

  disablewifi:
    desc:    "Disable WiFi"
    cmd:     "echo -e '# Disable bluetooth\ndtoverlay=pi3-disable-bt' >> /boot/config.txt"
    undocmd: "sed -i '/# Disable bluetooth/d;/dtoverlay=pi3-disable-bt/d' /boot/config.txt"
    check:   "grep -q 'dtoverlay=pi3-disable-bt' /boot/config.txt"
```

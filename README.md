# Game Boy Zero config menu

This is a simple config menu to be used at [Game Boy Zero projects](http://www.sudomod.com/wiki/index.php?title=Game_Boy_Zero)

It's a Work In Progress (still only at visual part =p)

## Sample menu definition (yaml file)
```
---
name: "Yaya's power optimization"

options:
  samba:
    desc: "Samba - Enabled"
    desc2: "Samba - Disabled"
    cmd: "/home/retropie/RetroPie-Setup/retropie_packages.sh samba depends remove"
    undocmd: "/home/retropie/RetroPie-Setup/retropie_packages.sh samba depends install && /home/retropie/RetroPie-Setup/retropie_packages.sh samba install_shares"
    check: "test -r /etc/samba/smb.conf"

  romservice:
    desc: "ROMService (for USB keys) - Enabled"
    desc2: "ROMService (for USB keys) - Disabled"
    cmd: "/home/retropie/RetroPie-Setup/retropie_packages.sh usbromservice enable"
    undocmd: "/home/retropie/RetroPie-Setup/retropie_packages.sh usbromservice disable"
    check: "test ! -r /opt/retropie/supplementary/usbromservice/disabled"

  tv:
    desc:    "TV and HDMI signals - Enabled"
    desc2:   "TV and HDMI signals - Disabled"
    cmd:     "echo '/usr/bin/tvservice -o' >> /etc/rc.local && /usr/bin/tvservice -o"
    undocmd: "sed -in -e '/tvservice -o/d' /etc/rc.local"
    check:   "grep -iq '/usr/bin/tvservice -o' /etc/rc.local"

  pileds:
    desc:    "Pi status LEDs - Enabled"
    desc2:   "Pi status LEDs - Disabled"
    cmd:     "echo -e '# Disable pi status LED\ndtparam=act_led_trigger=none\ndtparam=act_led_activelow=on' >> /boot/config.txt"
    undocmd: "sed -in -e '/# Disable pi status LED/,/dtparam=act_led_activelow=on/d' /boot/config.txt"
    check:   "grep -iq 'dtparam=act_led_trigger=none' /boot/config.txt"

  sdperf:
    desc:    "SD performance - Default"
    desc2:   "SD performance - Improved"
    cmd:     "sed -in -e 's/elevator=deadline/elevator=noop/' /boot/cmdline.txt"
    undocmd: "sed -in -e 's/elevator=noop/elevator=deadline/' /boot/cmdline.txt"
    check:   "grep -iq 'elevator=noop' /boot/cmdline.txt"

  swap:
    desc:    "Swap - Default"
    desc2:   "Swap - Adjusted"
    cmd:     "echo 'vm.swappiness = 1' >> /etc/sysctl.conf"
    undocmd: "sed -in -e '/vm.swappiness = 1/d' /etc/sysctl.conf"
    check:   "grep -iq 'vm.swappiness = 1' /etc/sysctl.conf"
    
  tweakfs:
    desc:    "Tweak filesystem for fewer writes - Disabled"
    desc2:   "Tweak filesystem for fewer writes - Enabled"
    cmd:     "sed -i -re 's/(.*boot.*) defaults (.*)/\\1 defaults,noatime \\2/' /etc/fstab"
    undocmd: "sed -i -re 's/(.*boot.*) defaults,noatime (.*)/\\1 defaults \\2/' /etc/fstab"
    check:   "grep -q 'boot.* defaults,noatime' /etc/fstab"

  avahi:
    desc:    "Avahi service - Enabled"
    desc2:   "Avahi service - Disabled"
    cmd:     "apt-get purge avahi-daemon -y"
    undocmd: "apt-get install avahi-daemon -y"
    check:   "! dpkg-query -s avahi-daemon"

  wifi:
    desc:    "WiFi - Enabled"
    desc2:   "WiFi - Disabled"
    cmd:     "echo -'# Disable wifi\ndtoverlay=pi3-disable-wifi' >> /boot/config.txt && ifdown wlan0"
    undocmd: "sed -i '/# Disable wifi/,/dtoverlay=pi3-disable-wifi/d' /boot/config.txt && ifup wlan0"
    check:   "grep -q 'dtoverlay=pi3-disable-wifi' /boot/config.txt"

  bluetooth:
    desc:    "Bluetooth - Enabled"
    desc2:   "Bluetooth - Disabled"
    cmd:     "echo -e '# Disable bluetooth\ndtoverlay=pi3-disable-bt' >> /boot/config.txt"
    undocmd: "sed -i '/# Disable bluetooth/,/dtoverlay=pi3-disable-bt/d' /boot/config.txt"
    check:   "grep -q 'dtoverlay=pi3-disable-bt' /boot/config.txt"
```

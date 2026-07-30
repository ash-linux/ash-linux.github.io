#include <tunables/global>

profile lsfs_daemon @{HOME}/.config/scripts/lsfs_daemon.py {
  #include <abstractions/base>
  #include <abstractions/python>
  #include <abstractions/nameservice>

  @{HOME}/.config/scripts/lsfs_daemon.py rm,
  owner @{HOME}/** r,
  owner @{HOME}/.config/scripts/** rw,
  /tmp/** rw,
  /var/log/ash/** rw,

  network tcp,
}

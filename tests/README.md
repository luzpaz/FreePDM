# Howto create a file server
The specs for a file server are minimal (a RPI4/5 should work). Right now I am using a desktop and it works for me, when you need to have a few guys then you can use a proper server. 

This are the things you need to do to make it work.

The program(s) run on Ubuntu Server or Debian. They support both the Intel and ARM platforms.
First you need to install Ubuntu Server (or Debian) and after that you need to write down the ip address of the server. When you need to have access outside of the LAN you are gonna need port forwarding and also a certificate (certbot for instance). This is the ip address of my server: 10.0.0.11 

To understand the IP address of your server, type in `ip a`

### Log into the server
`ssh ubuntu@10.0.0.11`

The password is `ubuntu`. The loginscreen asks you to modify the password.


### example usernames and passwords:

login name 1: user1, passwd: passwd1

login name 2: user2, passwd: passwd2

(Use different names and passwords, these are only as an example.)


The following commands are all in root (use `sudo su` or in debian `su -`):

```
groupadd vault

useradd -M -g vault vaultadmin
passwd vaultadmin

useradd -M -g vault user1
passwd user1

useradd -M -g vault user2
passwd user2

mkdir /vault
chown vault:vault /vault
chmod g+w /vault

exit
```

While still in ssh you can type `id user1`, `id user2` and `id vault` and write the number after `uid=` (but keep in mind to use different names). After you completed these tasks you can `exit` to exit the ssh. Next you can write down in the `FreePDM.toml` file that is stored (in Ubuntu) in the `.config/FreePDM` directory the following, at the bottom:

```
[Users]
vault = 1001
user1 = 1002
user2 = 1003
```

### Testing things out in Ubnutu Linux: 

Install pip modules:
```
pip install PySide2
pip install appdirs
pip install defusedxml
```

1. System update: `sudo apt update`
1. Install sshfs: `sudo apt install sshfs`
1. Make a directory : `sudo mkdir /mnt/test`
1. Mount your drive: `sudo sshfs -o allow_other user1@10.0.0.11:/vault /mnt/test` (Advise: Make a alias for that in your `.bashrc`)
1. Download FreePDM with this command: `git pull https://github.com/grd/FreePDM` 

### Creating a new Vault

In the FreePDM directory go to the directory `tests` and run `python3 create_new_vault.py`. 

### Testing a vault

In the FreePDM directory go to the directory `tests` and run `python3 fileserver.py` and if it runs then it is done.

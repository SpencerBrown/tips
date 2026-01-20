# How to install and run Linux on your Mac

You have a MacBook with Apple Silicon (m1 or higher), and you have kept up with macOS releases; hopefully you are on Tahoe. You want to run a Linux virtual machine for testing, learning, trying things out, reproducing environments for diagnosis of problems...

Here’s an easy step-by-step guide, You will end up with an Ubuntu Linux environment you can use via the Linux command line.

1. Check your prerequisites
2. Install UTM
3. Download a Linux boot image
4. Create a Virtual Machine (VM)
5. Set up Ubuntu Linux
6. Set up the SSH service
7. Set up sharing with your host macOS

# Check your prerequisites

* MacBook or Mac Mini with Apple silicon (m1 or higher)
* Kept up reasonably with macOS updates
* A reasonable amount of memory; 16MB minimum
* A reasonable amount of free disk space: I’d say 40-50 GB
* A decent network connection, you’ll be downloading a lot of software

# Download a Linux boot image

We suggest Ubuntu 24.04 and this document assumes that is what you are using,

* Open the “Ubuntu Server for ARM” [web page](https://ubuntu.com/download/server/arm).
* Click “Download 24.04.3 LTS” (or whatever the current LTS version is).
* A `.iso` file, approximately 3GB, is downloaded.

# Install UTM

[UTM](https://mac.getutm.app/) is a free Mac application that makes it easy to create and manage virtual machines on your Mac using Apple’s hypervisor framework. Go to [the web site](https://mac.getutm.app/) and click Download to get the .pkg file, or if you want to support them financially, click the Mac App Store button and buy it.

Double-click the downloaded file and follow the simple installation instructions (basically drag it to the applications folder),

# Create a VIrtual Machine

* Start UTM,
* Click “Create a new Virtual Machine”.
* click “Virtualize”.
* click “Linux”.
* adjust memory and CPU cores as needed. Suggest leaving as is for a first try.
  * We suggest using Apple Virtualization even though there’s a warning about it.
* Click “Continue”.
* click “Browse” and navigate to your downloaded boot image. Select it and click “Open”.
* click “Continue”.
* Adjust the virtual disk drive size as needed. The default of 64 GiB is reasonable; it is “thin provisioned”, meaning only the actually used space will be allocated on your macOS drive. Click “Continue”,
* Click “Continue” to bypass the setup of shared directories, this can be added later.
* Adjust the name as desired and click “Save”.

# Set up Ubuntu Linux

* In UTM, navigate to the VM you created and note that the `.iso` file is mounted in the virtual CD/DVD drive.
* You can right-click the VM name and click “Edit” to change any of the VM settings.
* Click the big triangle or “Run” to “power on” the virtual machine.
* You should see a console window, click “Try/Install Ubuntu” or just wait. You will see Ubuntu starting up and going to the install/configure screens
* Select “English” and hit Enter. Note the mouse won’t work here; you need to use the keyboard to navigate the windows.
* hit Enter for “Done” to bypass keyboard layout
* Hit Done to choose “Ubuntu Server”
* Hit “Done” to take the default network interface, assuming it’s DHCP. You could create a static IP here if you wanted to.
* Hit “Done” to skip setting a proxy
* Hit “Done” to accept the standard mirror location for downloading
* Hit “Done” to accept the default disk layout
* Hit “Done” to accept the storage layout summary
* Select the red “Continue” that warns about erasing the disk; it’s a brand new virtual disk so no worries.
* Fill in your name, your server name, your username and password. Hit “Done”.
* Do not select Ubuntu Pro (It’s a paid option), hit Continue
* Select the OpenSSH Server installation (use the space bar to select) (we will set this up later), hit “Done”
* Hit “Done”, don’t select any optional packages (unless you really want to\!)
* Now the installation runs and completes\!
* Hit “Reboot Now” which *should* unmount the install image from the virtual CD/DVD drive.
  * If it does not, power off the VM and clear the CD/DVD drive in the settings
* Login with the user/password you set up earlier.
* Update your packages:
  * `sudo apt update`
  * `sudo apt dist-upgrade`
* Enable the ssh service
  * `sudo systemctl ssh enable`
* Reboot Linux
  * `sudo systemctl reboot`

# Set up the SSH service

This will allow you to connect from a normal terminal window on your Mac instead of the console window.

## Initial connection

* Run `ip a` to get your IP address. Note it: mine is 192.168.64.5
* Open a Mac terminal window and run `ssh <username>@<ip-address>`
* Login with your password
* `exit` to leave your VM and go back to your terminal session

## Create and enable an SSH key

* On your Mac, in a terminal window, create the directory `.ssh` if it doesn’t already exist
* Run `ssh-keygen -f ~/.ssh/<username-ubuntu>` (you can pick a different filename if you want)
  * hit Enter twice for no passphrase
  * this will generate `<username-ubuntu>` and `<username-ubuntu>.pub` files in the `.ssh` directory
  * these are the SSH private and public keys, respectively
* Run `ssh-copy-id -i ~/.ssh/<username-ubuntu>.pub <username>@<ip-address>` and enter your password
  * this will copy your public key to the server’s `.ssh/authorizedKeys/` directory
  * and thus authorize your private SSH key to login to Linux as your username
* Run `ssh -i ~/.ssh/<username-ubuntu> <username>@<ip-address>` which should log you in to Linux without a password; you authenticated with your private key.

## Harden access to your VM

We want to remove the ability to login remotely via username/password, or to root. The only way to login is through the VM console, or via SSH to your username with your private key file. We assume you know now how to edit a file in Linux.

* Login to your VM via SSH.
* Run `sudo vi /etc/ssh/sshd_config` (or use an editor of your choice)
* Make sure the following items are set:
  * `PermitRootLogin no`
  * `PasswordAuthentication no`
* Save the file and restart Linux with `sudo systemctl reboot`

## Set up SSH configuration for easy connection

Instead of always typing (for example)
`ssh -i /Users/spencer/.ssh/spencer-ubuntu 'spencer@192.168.64.5’`
you can configure ssh in your mac to login with a nickname.

* On your local Mac terminal (**not** on your VM), edit the `~/.ssh/config` file and add a section like this:

`Host linux`
        `HostName 192.168.64.5`
        `IdentityFile ~/.ssh/spencer-ubuntu`
        `User spencer`

* Obviously, use your own username and IP address. The nickname for this connection is `linux`.
* Now, you can connect and login to your VM with just `ssh linux`

## Allow `sudo` without password

This is optional, but I find it annoying and of no security value to require re-entering your password to run sudo. To remove that requirement:

* login to LInux on your VM
* run `sudo visudo`
* using the editor, add this line to the bottom of the file and save it (obviously, use your own username instead of “spencer”:

`spencer ALL=(ALL) NOPASSWD: ALL`

* Now, you should be able to run `sudo` without requiring you to re-enter your password.

# Share local directories with Linux

* Power off your Linux VM with `sudo systemctl poweroff`
* Add one or more macOS directories to the shared directory list in the UTM VM settings
* Start the Linux VM and log in via SSH
* Create a mount point directory; for example, `sudo mkdir /mnt/utm`
* Mount the shared directory: `sudo mount -t virtiofs share /mnt/utm`
* You should now be able to see the shared directories under `/mnt/utm`


To make the mount permanent, add a line to `/etc/fstab:`

`share  /mnt/utm        virtiofs        rw,nofail       0       0`

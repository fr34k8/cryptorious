# [![CircleCI](https://circleci.com/gh/malnick/cryptorious.svg?style=svg)](https://circleci.com/gh/malnick/cryptorious)

Like 1Password but for the CLI. Stores your encrypted data in eyaml using generic SSH keys as the basis for encryption/decryption so you never have to type a password to get your passwords ever again.

## Download
### Linux
- AMD64 | [v1.2.1](https://dl.dropboxusercontent.com/u/77193293/tools/cryptorious_1.2.1)
- AMD64 | [v1.2.0](https://dl.dropboxusercontent.com/u/77193293/tools/cryptorious_1.2.0)
- AMD64 | [v1.1.0](https://dl.dropboxusercontent.com/u/77193293/tools/cryptorious_1.1.0)
- AMD64 | [v1.0.0](https://dl.dropboxusercontent.com/u/77193293/tools/cryptorious)

### Darwin (OSx)
- AMD64 | [v1.2.1](https://dl.dropboxusercontent.com/u/77193293/tools/cryptorious_1.2.1_darwin)

## Manpage
### Main Menu
```
NAME:
   
 _________                            __                   .__                        
 \_   ___ \ _______  ___.__.______  _/  |_   ____  _______ |__|  ____   __ __   ______
 /    \  \/ \_  __ \<   |  |\____ \ \   __\ /  _ \ \_  __ \|  | /  _ \ |  |  \ /  ___/
 \     \____ |  | \/ \___  ||  |_> > |  |  (  <_> ) |  | \/|  |(  <_> )|  |  / \___ \ 
  \______  / |__|    / ____||   __/  |__|   \____/  |__|   |__| \____/ |____/ /____  >
         \/          \/     |__|                                                   \/ 
 - CLI-based encryption for passwords and random data

USAGE:
   cryptorious [global options] command [command options] [arguments...]
   
VERSION:
   1.2.1
   
AUTHOR(S):
   Jeff Malnick <malnick@gmail.com> 
   
COMMANDS:
    rename	 Rename an entry in the vault
    rotate	 Rotate your cryptorious SSH keys and vault automatically
    delete	 Remove an entry from the cryptorious vault
    decrypt	 Decrypt a value in the vault `VALUE`
    encrypt	 Encrypt a value for the vault `VALUE`
    generate Generate a RSA keys or a secure password.	

GLOBAL OPTIONS:
   --vault-path, --vp "/home/malnick/.cryptorious/vault.yaml"         Path to vault.yaml
   --private-key, --priv "/home/malnick/.ssh/cryptorious_privatekey"  Path to private key
   --public-key, --pub "/home/malnick/.ssh/cryptorious_publickey"     Path to public key
   --debug                                                            Debug/Verbose log output
   --help, -h                                                         Show help
   --version, -v                                                      Print the version

```
### Decrypt Sub Menu
```   
NAME:
   cryptorious decrypt - Decrypt a value in the vault `VALUE`

USAGE:
   cryptorious decrypt [command options] [arguments...]

OPTIONS:
   --copy, -c           Copy decrypted password to clipboard automatically
   --goto, -g           Open your default browser to https://<key_name> and login automatically
   --timeout, -t "10"   Timeout in seconds for the decrypt session window to expire
```   
### Rename Sub Menu
```
NAME:
   cryptorious rename - Rename an entry in the vault

USAGE:
   cryptorious rename [command options] [arguments...]

OPTIONS:
   --old, -o    Name of old entry name [key] in vault
   --new, -n    Name of new entry name [key] in vault
```
### Generate Sub Menu
```
NAME:
 generate - 	Generate a RSA keys or a secure password 

USAGE:
  generate command [command options] [arguments...]

COMMANDS:
    keys	                 Generate SSH key pair for cryptorious
    password	[--[l]ength] Generate a random password

OPTIONS:
   --help, -h	show help

```

## Step 0: Build && Alias

Build it and install: `make install`

Add to your `.[bash | zsh | whatever]rc`: `alias cpt=cryptorious`

## Step 1: Generate keys

```
cryptorious generate keys 
```

Defaults to placing keys in ```$HOME/.ssh/cryptorious_privatekey``` and ```$HOME/.ssh/cryptorious_publickey```.

You can override this with ```--private-key``` and ```--public-key```:

```
cryptorious generate keys --private-key foo_priv --public-key foo_pub 
```

### Lock It Down
If you want to win extra security stars, lock down your keys with root ownership. By default they're already read/write by the user who ran the `cryptorious` command (0600), but you can increase this security more with `chmod root:root ~/.ssh/cryptorious_privatekey`. Now you'll have to run `cryptorious` with `sudo` and enter in your root password (ugh, passwords..) every time. 

## Step 2: Encrypt

```
cryptorious encrypt github  
```

Will open a ncurses window and prompt you for username, password and a secure note. All input is optional. 


## Step 3: Decrypt 

```
cryptorious decrypt thing
```

Will open a ncurses window with the decrypted vault entry. 

Forgo the the ncurses window and copy the decrypted password stright to the system clipboard? 
```
cryptorious decrypt -[c]opy thing
```
No printing, just a message that your decrypted password is now available in the paste buffer for your user. 

If you've saved your vault entries with the URI of the site they belong to (i.e., ran `cryptorious encrypt github.com`...) then you can use the `-[g]oto` flag to open your default browser to this URI. Pair it with `-[c]opy` and the shorthand for `[d]ecrypt` and you'll have a fast way of navigating directly to your desired, secure website (let's also assume you've aliased `cpt=cryptorious`):
```
cpt d -g -c github.com
```

## Step 4: Rotate Keys & Vault
Compromised your keys? Not a problem. 

```
cryptorious rotate
```

1. Backs up your old keys to `keyPath.bak`
1. Backs up your old vault to `vaultPath.bak`
1. Generates new keys to `keyPath`
1. Decrypts vault using `cryptorious_privatekey.bak` and encrypts vault in place with new `cryptorious_publickey`
1. Writes the vault back to disk at `vaultPath`

## Step 5: Generate Secure Password
The `generate` command also lets you generate random, secure passwords of `n` length:
```
cryptorious generate password --length 20
(yZkj,GX`w7T4x&TaYyw
```

This defaults to a length of 15 if you don't pass --[l]ength.

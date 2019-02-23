# gokpcli

Simple KeePass2 console client for linux using a simple command system
(with some commands beeing inspired on linux commands)

This project is heavily inspired on [http://kpcli.sourceforge.net/](http://kpcli.sourceforge.net/), 
my attempt is to recreate something similar using something more 'modern' than perl.
However while **kpcli** takes over the terminal creating a new screen, my approach does not

## Usage

The following flags are available when using `gokpcli`:

- `nbackup` If set when saving the database no backups will be done
- `pwfile` File where your database password is stored (instead of passing the plaintext password)
- `db` Database filepath 

### Example

```
./gokpcli -db=/mnt/raggaer_g/KeePass/Databases/raggaer_test.kdbx -pwfile=/home/raggaer/.kpcli-master
```

## Commands

You can use the `help` command while a database is opened to view the list of commands:

- `xp` Copies the password of an entry
- `xu` Copies the username of an entry
- `xw` Copies the URL (www) of an entry
- `ls` Lists all the groups and entries of the current group
- `cd` Changes the current working group
- `exit` Closes the application
- `save` Saves the database
- `mkdir` Shows and processes a form to create a new group
- `rmdir` Removes a group (sends the group to the recycle bin)
- `new` Shows and processes a form to create a new entry
- `edit` Modifies an entry
- `rm` Removes an entry from the current working group (sends the entry to the recycle bin)
- `show` Shows an entry from the current working group
- `search` Searches entries (by title) from the current working group
- `save` Saves the database to disk
- `xx` Clears the clipboard

## Deleting groups and entries

When an entry or a group is deleted we move it to the `Recycle Bin` group (this will be created if its missing).
You can delete delete the entry forever or just leave it there as some sort of backup folder

After deleting a backup of the database file is created (before the delete change) with the format `y-m-d_h:i:s_name.kdbx`

## Clipboard

Commands like `xu` and `xp` copy the content to the system clipboard, making use of [github.com/atotto/clipboard](https://github.com/atotto/clipboard).
You will need `xclip` or `xsel` installed

## License

**gokpcli** is licensed under the **GNU GPLv3**, basically you can do almost anything you want
with this project, except to distribute closed source versions

This application is mainly using [github.com/tobischo/gokeepasslib](http://github.com/tobischo/gokeepasslib) package
to modify the KeePass database

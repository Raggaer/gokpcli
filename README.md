# gokpcli

Simple KeePass2 console client for linux using a simple command system
(with some commands beeing inspired on linux commands)

This project is heavily inspired on [http://kpcli.sourceforge.net/](http://kpcli.sourceforge.net/), 
my attempt is to recreate something similar using something more 'modern' than perl.
However while **kpcli** takes over the terminal creating a new screen, my approach does not

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
- `new` Shows and processes a form to create a new entry
- `rm` Removes an entry from the current working group
- `show` Shows an entry from the current working group
- `search` Searches entries (by title) from the current working group
- `save` Saves the database to disk

## Clipboard

Commands like `xu` and `xp` copy the content to the system clipboard, making use of [https://github.com/atotto/clipboard](https://github.com/atotto/clipboard).
You will need `xclip` or `xsel` installed

## License

**gokpcli** is licensed under the **GNU GPLv3**, basically you can do almost anything you want
with this project, except to distribute closed source versions

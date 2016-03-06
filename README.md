# bigbin
Package Big Binary makes it easy to bundle all your go programs into a single big binary

## Usage

To generate a Big Binary using appa and appb in the following directory tree:
```
  ├── appa
  └── appb
```

Do:
```bash
  $ genbigbin --to mybigbin --apply ./appa ./appb
```  
That will:
* Turn appa and appb into applets of mybigbin
* Create also stand alone versions in appa/appa & appb/appb

Then run mybigbin in the binary folder of your choice to autogenerate all symlinks:
```bash
  $ mybigbin
  Rebuilding symlinks:
  {full path here}/mybigbin -> appa
  {full path here}/mybigbin -> appb
```

Finally you can invoke each app individually, from the same big bynary executable file, by using the appropriate symlink:

```bash
  $ ./appa
  App A running...
```

```bash
  $ ./appb
  App B doing its thing...
```

## Reverting changes

In case you changed your mind and want to go back to your code without bigbin support, do:
```bash
  $ genbigbin --restore --to mybigbin --apply ./appa ./appb
```

### General help

```bash
$ genbigbin
Usage:
genbigbin [flags] mainDir1 [mainDir2...]

Flags:
  -apply
    	Apply changes to the filesystem (false by default)
  -restore
    	Restore files to before the big binary changes intead (false by default)
  -to string
    	Directory in where to generate the big binary main (by default is empty and does not create a big binary main)
```

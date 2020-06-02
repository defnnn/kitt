## Install asdf-vm

See instructions at https://github.com/asdf-vm/asdf

## Install kitt asdf plugins

```
asdf plugin-add letfn--shflags https://github.com/defn/asdf-plugin
asdf plugin-add defn--sub https://github.com/defn/asdf-plugin
asdf plugin-add defn--kitt https://github.com/defn/asdf-plugin
asdf plugin-add defn--kiki https://github.com/defn/asdf-plugin
```

## Configure kitt versions

Specify the versions of each kitt project in `$HOME/.tool-versions` which is used by asdf to install kitt software.  The first column is the name of the plugin, the second column is the version.

```
letfn--shflags defn1
defn--sub defn3
defn--kitt defn3
defn--kiki defn3
```

## Install kitt software

`asdf install`

## Configure $PATH

In addition to the normal asdf setup, put this into your `PATH`: `$HOME/.asdf/install/bin`.  It has symlinks to the installed kitt software.

## Run kitt cli

```
kitt
kitt up --help
```

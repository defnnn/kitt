---
layout: blog
category: blog
published: false
title: Untitled
tags:
  - mangos
---
## Install asdf-vm

See instructions at https://github.com/asdf-vm/asdf

## Configure $PATH

In addition to the normal asdf setup, put this into your `PATH`: `$HOME/.asdf/installs/bin`.  It has symlinks to the installed kitt software.

## Install kitt asdf plugins

```
asdf plugin-add defn--asdf-plugin https://github.com/defn/asdf-plugin
asdf plugin-add shflags https://github.com/letfn/shflags
asdf plugin-add sub https://github.com/defn/sub
asdf plugin-add kitt https://github.com/defn/kit
asdf plugin-add kiki https://github.com/defn/kiki
```

## Configure kitt versions

Specify the versions of each kitt project in `$HOME/.tool-versions` which is used by asdf to install kitt software.  The first column is the name of the plugin, the second column is the version.

```
defn--asdf-plugin defn1
shflags defn7
sub defn7
kitt defn8
kiki defn7
```

## Install kitt software

`asdf install`

## Run kitt cli

```
kitt
kitt up --help
```

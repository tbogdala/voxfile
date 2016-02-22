VOXFILE
=======

Voxfile is a library to load .vox files which is a storage solution for
voxel models and scenes. Editors like [MagicaVoxel][1].

This currently supports file version 150.

Installation
------------

```
go get github.com/tbogdala/Voxfile
```

Usage
-----

A `.vox` file can be loaded with the following statement:

```
voxFile, err := DecodeFile("testdata/chr_sword.vox")
```


[1]: https://ephtracy.github.io/

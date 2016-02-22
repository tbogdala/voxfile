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

License
-------

Voxfile is released under the BSD license. See the [LICENSE][license-link] file for more details.


[1]: https://ephtracy.github.io/
[license-link]: https://raw.githubusercontent.com/tbogdala/voxfile/master/LICENSE

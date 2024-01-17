# Go Ethereum sidechain

Implementation of an Ethereum sidechain for Bitcoin, using Drivechain (BIP300/301). 

## Dependencies

The following things are needed before you can build this project. Obtaining them is left as an excercise to the reader. 

1. Rust (`cargo`). The project uses a [Drivechain library](drivechain/drivechain_linux.go) written in Rust.
2. C compiler. Needed for using the compiled Rust bindings. 
3. Go. 
4. `make`

## Getting started

```bash
$ make sidegeth

# Tweak these values as needed
$ ./build/bin/sidegeth \
      --main.host=localhost    --main.user=user \
      --main.password=password --main.port=18443
```


## Windows

If you're on Windows, things are more complicated. The first step here is 
to reevaluate your choices, and use a proper OS. If you still insist on 
using Windows, do this: 

1. Install dependencies using Chocolatey: `choco install golang make rust mingw`   
   In case you're wondering, `mingw` is a GNU distribution for Windows.  
2. Download [MSYS2](https://www.msys2.org/)  
   This gives you (among other things) a working `bash` shell for Windows. 
3. Clone, compile and install [`dlfcn-win32`](https://github.com/dlfcn-win32/dlfcn-win32). This 
   is a dependency for the Rust Drivechain library.  
   1. Clone the repo
   2. Open up a `bash` shell (**NOT** through WSL, but native Windows `bash` from MSYS2)
   3. Build the library: `./configure --prefix=/ --libdir=$PWD/libdir --incdir=$PWD/incdir && make`
   4. Install the library by placing `./libdir/libdl.a` somehere `ld` can find it. One such 
      location can be `C:\ProgramData\mingw64\mingw64\lib`, but who knows if this is a horrible
      idea. This guide was written by a Windows noob.
4. Build: `make sidegeth`
   

# Compiler for nop
Please refer to the contribution guidelines [here](CONTRIBUTING.md)

## Installation :


### Prerequisites :

You have to install antlr4 first:
- With your package manager <br>
  Make sure that `antlr4` works and has Cpp compatibility

- By Downloading latest complete version at https://www.antlr.org/download/index.html <br>
and moving the .jar file at /usr/local/bin/antlr4.jar

You have to install antlr4-runtime library:
- Download latest C++ runtime lib at https://www.antlr.org/download/index.html and extract it<br>
  To compile it :
  ```
  $ mkdir build && mkdir run && cd build
  $ cmake .. -DANTLR_JAR_LOCATION=full/path/to/antlr4-VERSION-complete.jar -DWITH_DEMO=True
  $ make
  $ DESTDIR=../runtime/Cpp/run make install
  ```
  You then have to move the created files to your local files :
  ```
  $ cd ../runtime/Cpp/run/
  $ cp usr/local/lib/* /usr/local/lib/
  $ cp -r usr/local/include/antlr4-runtime/ /usr/local/include/
  $ cp -r usr/local/share /usr/local/share/
  ```

### Building the nopc-bootstrap:

```./build.sh```

The binary of nopc is located at **./bin/nopc**
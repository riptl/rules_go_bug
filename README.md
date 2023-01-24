Minimal reproducible example of an annoying bug with `rules_go` and Cgo dependencies.

Causes Cgo dependencies to be built twice when only once is enough.

## Structure

- `//clib`: A `cc_library`
- `//golib`: A `go_library` and `go_test`

## Problem

It is sufficient to build `//clib` once.

`rules_go` emits transitions that cause `//clib` to be built twice.

```
$ bazel aquery --output=jsonproto 'deps(//golib/...)' | jq -r '.actions[] | select(.mnemonic == "CppCompile") | .arguments | join(" ")'
Loading:
Loading: 0 packages loaded
Analyzing: 2 targets (0 packages loaded, 0 targets configured)
INFO: Analyzed 2 targets (0 packages loaded, 0 targets configured).
INFO: Found 2 targets...
INFO: Elapsed time: 0.352s, Critical Path: 0.00s
INFO: 0 processes.
INFO: Build completed successfully, 0 total actions
/bin/gcc -U_FORTIFY_SOURCE -fstack-protector -Wall -Wunused-but-set-parameter -Wno-free-nonheap-object -fno-omit-frame-pointer -MD -MF bazel-out/k8-fastbuild/bin/clib/_objs/clib/lib.pic.d -frandom-seed=bazel-out/k8-fastbuild/bin/clib/_objs/clib/lib.pic.o -fPIC -DBAZEL_CURRENT_REPOSITORY="" -iquote . -iquote bazel-out/k8-fastbuild/bin -fno-canonical-system-headers -Wno-builtin-macro-redefined -D__DATE__="redacted" -D__TIMESTAMP__="redacted" -D__TIME__="redacted" -c clib/lib.c -o bazel-out/k8-fastbuild/bin/clib/_objs/clib/lib.pic.o
/bin/gcc -U_FORTIFY_SOURCE -fstack-protector -Wall -Wunused-but-set-parameter -Wno-free-nonheap-object -fno-omit-frame-pointer -MD -MF bazel-out/k8-fastbuild-ST-1665e0aa65f7/bin/clib/_objs/clib/lib.pic.d -frandom-seed=bazel-out/k8-fastbuild-ST-1665e0aa65f7/bin/clib/_objs/clib/lib.pic.o -fPIC -DBAZEL_CURRENT_REPOSITORY="" -iquote . -iquote bazel-out/k8-fastbuild-ST-1665e0aa65f7/bin -fno-canonical-system-headers -Wno-builtin-macro-redefined -D__DATE__="redacted" -D__TIMESTAMP__="redacted" -D__TIME__="redacted" -c clib/lib.c -o bazel-out/k8-fastbuild-ST-1665e0aa65f7/bin/clib/_objs/clib/lib.pic.o
/bin/gcc -U_FORTIFY_SOURCE -fstack-protector -Wall -Wunused-but-set-parameter -Wno-free-nonheap-object -fno-omit-frame-pointer -g0 -O2 -D_FORTIFY_SOURCE=1 -DNDEBUG -ffunction-sections -fdata-sections -std=c++0x -MD -MF bazel-out/k8-opt-exec-2B5CBBC6/bin/external/bazel_tools/src/tools/launcher/_objs/launcher/dummy.pic.d -frandom-seed=bazel-out/k8-opt-exec-2B5CBBC6/bin/external/bazel_tools/src/tools/launcher/_objs/launcher/dummy.pic.o -fPIC -DBAZEL_CURRENT_REPOSITORY="bazel_tools" -iquote external/bazel_tools -iquote bazel-out/k8-opt-exec-2B5CBBC6/bin/external/bazel_tools -g0 -g0 -fno-canonical-system-headers -Wno-builtin-macro-redefined -D__DATE__="redacted" -D__TIMESTAMP__="redacted" -D__TIME__="redacted" -c external/bazel_tools/src/tools/launcher/dummy.cc -o bazel-out/k8-opt-exec-2B5CBBC6/bin/external/bazel_tools/src/tools/launcher/_objs/launcher/dummy.pic.o
/bin/gcc -U_FORTIFY_SOURCE -fstack-protector -Wall -Wunused-but-set-parameter -Wno-free-nonheap-object -fno-omit-frame-pointer -g0 -O2 -D_FORTIFY_SOURCE=1 -DNDEBUG -ffunction-sections -fdata-sections -std=c++0x -MD -MF bazel-out/k8-opt-exec-2B5CBBC6/bin/external/bazel_tools/src/tools/launcher/_objs/launcher/dummy.d -frandom-seed=bazel-out/k8-opt-exec-2B5CBBC6/bin/external/bazel_tools/src/tools/launcher/_objs/launcher/dummy.o -DBAZEL_CURRENT_REPOSITORY="bazel_tools" -iquote external/bazel_tools -iquote bazel-out/k8-opt-exec-2B5CBBC6/bin/external/bazel_tools -g0 -g0 -fno-canonical-system-headers -Wno-builtin-macro-redefined -D__DATE__="redacted" -D__TIMESTAMP__="redacted" -D__TIME__="redacted" -c external/bazel_tools/src/tools/launcher/dummy.cc -o bazel-out/k8-opt-exec-2B5CBBC6/bin/external/bazel_tools/src/tools/launcher/_objs/launcher/dummy.o
```

`lib.pic.o` is created twice.

- `bazel-out/k8-fastbuild/bin/clib/_objs/clib/lib.pic.o`
- `bazel-out/k8-fastbuild-ST-1665e0aa65f7/bin/clib/_objs/clib/lib.pic.o`

```
$ cmp "bazel-out/k8-fastbuild/bin/clib/_objs/clib/lib.pic.o" "bazel-out/k8-fastbuild-ST-*/bin/clib/_objs/clib/lib.pic.o" && echo "same"
same
```

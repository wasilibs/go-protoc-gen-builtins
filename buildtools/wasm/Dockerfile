FROM ghcr.io/webassembly/wasi-sdk:wasi-sdk-24

RUN apt-get update && apt-get install -y binaryen git curl patch python3

WORKDIR /workspace
ADD buildtools/wasm/version.txt version.txt

RUN curl -L https://github.com/protocolbuffers/protobuf/archive/refs/tags/$(cat version.txt | awk '{$1=$1};1').tar.gz | tar -xz --strip-components 1 -C /workspace

ADD buildtools/wasm/patch-absl.txt patch-absl.txt
ADD buildtools/wasm/patch.txt patch.txt
RUN patch -p1 < patch.txt

# From https://github.com/protocolbuffers/protobuf/blob/v25.1/src/google/protobuf/compiler/cpp/main.cc
ADD buildtools/wasm/main.cc main.cc
ADD buildtools/wasm/create_plugin_entrypoints.py create_plugin_entrypoints.py
RUN python3 create_plugin_entrypoints.py

# libprotoc, the library that code generator plugins link against, includes the command line interface, including
# subprocess execution. Code generator plugins never need this, so we can remove the files to allow the project to
# build with wasi-sdk.
RUN rm src/google/protobuf/compiler/subprocess.* src/google/protobuf/compiler/command_line_interface.*
RUN sed -i '/src\/google\/protobuf\/compiler\/subprocess\./d' src/file_lists.cmake
RUN sed -i '/src\/google\/protobuf\/compiler\/command_line_interface\./d' src/file_lists.cmake

ENV CMAKE_TOOLCHAIN_FILE=/opt/wasi-sdk/share/cmake/wasi-sdk-pthread.cmake
ENV CFLAGS="-D_WASI_EMULATED_MMAN -D_WASI_EMULATED_PROCESS_CLOCKS -D_WASI_EMULATED_SIGNAL -DABSL_HAVE_MMAP -DABSL_FORCE_THREAD_IDENTITY_MODE=1"
ENV CXXFLAGS="$CFLAGS -fno-exceptions"
ENV LDFLAGS="-lwasi-emulated-process-clocks -lwasi-emulated-mman -lwasi-emulated-signal -Wl,--max-memory=4294967296 -Wl,--global-base=1024"

RUN cmake -DCMAKE_BUILD_TYPE=Release -DCMAKE_C_FLAGS="$CFLAGS" -DCMAKE_CXX_FLAGS="$CXXFLAGS" -DCMAKE_EXE_LINKER_FLAGS="$LDFLAGS" -Dprotobuf_BUILD_TESTS=off -S . -B build

# cmake automatically downloads abseil in previous step so we patch here before then building
RUN cd build/_deps/absl-src/ && patch -p1 < ../../../patch-absl.txt

RUN cmake --build build --target plugins protoc-gen-upb protoc-gen-upbdefs protoc-gen-upb_minitable -j 2

SHELL ["/bin/bash", "-c"]
RUN for f in build/protoc-gen-*[^0-9]; do wasm-opt -o "$f".wasm --low-memory-unused --flatten --rereloop --converge -O3 "$f"; done

CMD ["bash", "-c", "cp build/protoc-gen-*.wasm /out/"]

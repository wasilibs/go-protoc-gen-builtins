FROM ghcr.io/wasilibs/wasix-sdk:sha-fc94d60

RUN apt-get update && apt-get install -y binaryen git patch

ADD buildtools/wasm/version.txt version.txt
RUN git clone --recursive https://github.com/protocolbuffers/protobuf.git --branch $(cat version.txt | awk '{$1=$1};1') /workspace

ENV CFLAGS "${CFLAGS} -O3 -pthread"
ENV CXXFLAGS "${CXXFLAGS} -O3 -pthread -D_WASI_EMULATED_MMAN -D_WASI_EMULATED_PROCESS_CLOCKS -DABSL_HAVE_MMAP"
ENV LDFLAGS $LDFLAGS -lwasi-emulated-mman -lwasi-emulated-process-clocks -Wl,--max-memory=4294967296 -Wl,--export=__stack_pointer -Wl,--export=__heap_base

WORKDIR /workspace

ADD buildtools/wasm/patch-absl.txt patch-absl.txt
RUN cd third_party/abseil-cpp && patch -p1 < ../../patch-absl.txt

# From https://github.com/protocolbuffers/protobuf/blob/v25.1/src/google/protobuf/compiler/cpp/main.cc
ADD buildtools/wasm/main.cc main.cc
ADD buildtools/wasm/create_plugin_entrypoints.py create_plugin_entrypoints.py
RUN python3 create_plugin_entrypoints.py


RUN cmake -S . -B build
RUN cmake --build build --target plugins protoc-gen-upb protoc-gen-upbdefs protoc-gen-upb_minitable

SHELL ["/bin/bash", "-c"]
RUN for f in build/protoc-gen-*[^0-9]; do wasm-opt -o "$f".wasm --low-memory-unused --flatten --rereloop --converge -O3 "$f"; done

CMD ["bash", "-c", "cp build/protoc-gen-*.wasm /out/"]

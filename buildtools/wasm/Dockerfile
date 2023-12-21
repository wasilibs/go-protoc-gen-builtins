FROM ghcr.io/wasilibs/wasix-sdk:sha-fc94d60

RUN apt-get update && apt-get install -y binaryen git patch

RUN git clone --recursive https://github.com/protocolbuffers/protobuf.git --branch v25.1 /workspace

WORKDIR /workspace
# From https://github.com/protocolbuffers/protobuf/blob/v25.1/src/google/protobuf/compiler/cpp/main.cc
ADD buildtools/wasm/main.cc main.cc
ADD buildtools/wasm/create_plugin_entrypoints.py create_plugin_entrypoints.py
RUN python3 create_plugin_entrypoints.py

ENV CXXFLAGS "${CXXFLAGS} -O3 -pthread -D_WASI_EMULATED_MMAN -D_WASI_EMULATED_PROCESS_CLOCKS"
ENV LDFLAGS $LDFLAGS -lwasi-emulated-mman -lwasi-emulated-process-clocks -Wl,--max-memory=4294967296

RUN cmake -S . -B build
RUN cmake --build build --target plugins

SHELL ["/bin/bash", "-c"]
RUN for f in build/protoc-gen-*[^0-9]; do wasm-opt -o "$f".wasm --low-memory-unused --flatten --rereloop --converge -O3 "$f"; done

CMD ["bash", "-c", "cp build/protoc-gen-*.wasm /out/"]

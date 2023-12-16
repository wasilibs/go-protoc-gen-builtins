FROM ghcr.io/wasilibs/wasix-sdk:sha-fc94d60

RUN apt-get update && apt-get install -y binaryen git patch

RUN git clone --recursive https://github.com/protocolbuffers/protobuf.git --branch v25.1 /workspace

ENV CXXFLAGS "${CXXFLAGS} -O3 -pthread -D_WASI_EMULATED_MMAN -D_WASI_EMULATED_PROCESS_CLOCKS"
ENV LDFLAGS $LDFLAGS -lwasi-emulated-mman -lwasi-emulated-process-clocks -Wl,--max-memory=4294967296

WORKDIR /workspace
RUN cmake -S . -B build
RUN cmake --build build --target protoc

RUN wasm-opt -o build/protoc.wasm --low-memory-unused --flatten --rereloop --converge -O3 build/protoc-*

CMD ["bash", "-c", "cp build/protoc.wasm /out/"]

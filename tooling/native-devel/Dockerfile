FROM oreilly/rac-native-runtime:latest


# cuda-core-$CUDA_PKG_VERSION \
# cuda-misc-headers-$CUDA_PKG_VERSION \
# cuda-command-line-tools-$CUDA_PKG_VERSION \
# cuda-nvrtc-dev-$CUDA_PKG_VERSION \
# cuda-nvml-dev-$CUDA_PKG_VERSION \
# cuda-nvgraph-dev-$CUDA_PKG_VERSION \
# cuda-cusolver-dev-$CUDA_PKG_VERSION \
# cuda-cublas-dev-$CUDA_PKG_VERSION \
# cuda-cufft-dev-$CUDA_PKG_VERSION \
# cuda-curand-dev-$CUDA_PKG_VERSION \
# cuda-cusparse-dev-$CUDA_PKG_VERSION \
# cuda-npp-dev-$CUDA_PKG_VERSION \
# cuda-cudart-dev-$CUDA_PKG_VERSION \
# cuda-driver-dev-$CUDA_PKG_VERSION \

# RUN apt-get update && apt-get install -y --no-install-recommends \
#       protobuf-compiler \
#       && \
#      rm -rf /var/lib/apt/lists/*

ADD run.sh /usr/bin/run.sh

ENV LIBRARY_PATH /usr/local/cuda/lib64/stubs:${LIBRARY_PATH}

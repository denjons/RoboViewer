# Setup kafka dependencies

xcode-select --install

brew install pkg-config

export PKG_CONFIG_PATH=\$PKG_CONFIG_PATH:/usr/lib/pkgconfig

brew install librdkafka

go get gopkg.in/confluentinc/confluent-kafka-go.v1/kafka

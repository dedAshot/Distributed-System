kafka-metadata-quorum --bootstrap-server kafka-broker-1:19092 describe --status

kafka-topics --bootstrap-server kafka-broker-1:19092 --create --if-not-exists --topic messages --replication-factor 2 --partitions 3
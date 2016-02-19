# File Producer for Kafka

This tool is a Kafka producer to be used as a side-ckick container to push files into Kafka.

After building the tool, you can run it using the following command:

        FILE="file_to_be_sent" TOPIC="TOPIC_USED" PARTITION=0 BROKERS="localhost:9092" CLIENT_NAME="mygotest" ./file_producer


If you want to use the container, the command is:

        docker run -e FILE="myfile.zip" -e TOPIC="TEST_GO" \
            -e PARTITION=0 -e BROKERS="localhost:9092"  \
            -e CLIENT_NAME="ivantest2" \
            -v $(pwd)/myfile.zip:/myfile.zip \
            quay.io/ukhomeofficedigital/docker-file-producer


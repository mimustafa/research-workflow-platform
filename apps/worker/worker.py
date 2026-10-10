import json
import logging
import os
import signal
import sys

import pika


QUEUE_NAME = "workflow.jobs"
logging.basicConfig(level=logging.INFO, format="%(asctime)s %(levelname)s %(message)s")
logger = logging.getLogger(__name__)


def main() -> None:
    parameters = pika.URLParameters(
        os.getenv(
            "RABBITMQ_URL",
            "amqp://research:research_local_password@localhost:5672/",
        )
    )
    connection = pika.BlockingConnection(parameters)
    channel = connection.channel()
    channel.queue_declare(queue=QUEUE_NAME, durable=True)
    channel.basic_qos(prefetch_count=1)

    def handle_message(channel, method, _properties, body) -> None:
        try:
            job = json.loads(body)
            logger.info("received workflow job", extra={"workflow_id": job["id"]})
            channel.basic_ack(delivery_tag=method.delivery_tag)
        except (json.JSONDecodeError, KeyError) as error:
            logger.error("discarding malformed workflow job: %s", error)
            channel.basic_nack(delivery_tag=method.delivery_tag, requeue=False)

    channel.basic_consume(queue=QUEUE_NAME, on_message_callback=handle_message)
    logger.info("worker is ready", extra={"queue": QUEUE_NAME})
    channel.start_consuming()


if __name__ == "__main__":
    try:
        main()
    except KeyboardInterrupt:
        logger.info("worker stopped")
        sys.exit(0)

import json
import logging
import os
import signal
import sys

import pika
import psycopg


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
    database = psycopg.connect(
        os.getenv(
            "DATABASE_URL",
            "postgresql://research:research_local_password@localhost:5432/research_workflows",
        ),
        autocommit=True,
    )

    def handle_message(channel, method, _properties, body) -> None:
        try:
            job = json.loads(body)
            workflow_id = job["id"]
            with database.cursor() as cursor:
                cursor.execute(
                    "UPDATE workflows SET status = 'processing' WHERE id = %s",
                    (workflow_id,),
                )
                result = (
                    f"Accepted research brief with {len(job['brief'].split())} words "
                    "for deterministic analysis."
                )
                cursor.execute(
                    "UPDATE workflows SET status = 'completed', result = %s WHERE id = %s",
                    (result, workflow_id),
                )
            logger.info("completed workflow job", extra={"workflow_id": workflow_id})
            channel.basic_ack(delivery_tag=method.delivery_tag)
        except (json.JSONDecodeError, KeyError, psycopg.Error) as error:
            logger.error("discarding malformed workflow job: %s", error)
            if "job" in locals() and "id" in job:
                with database.cursor() as cursor:
                    cursor.execute(
                        "UPDATE workflows SET status = 'failed', result = %s WHERE id = %s",
                        (str(error), job["id"]),
                    )
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

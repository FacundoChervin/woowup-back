import {
  Queue,
  StackContext,
} from "sst/constructs";

import { Duration } from "aws-cdk-lib";


export function QueueStore({
  stack,
  app,
}: StackContext) {
  const queueStore = new Queue(
    stack,
    "Queue",
    {
      consumer: {
        function: {
          handler: "functions/lambda/sendqueue/main.go",
          functionName: "store_emails_queue",
          timeout: "120 seconds",
        },
        cdk: {
          eventSource: {
            batchSize: 1,
          },
        },
      },
      cdk: {
        queue: {
          fifo: true,
          queueName: process.env.SEND_SQS_QUEUE_NAME,
          visibilityTimeout: Duration.seconds(120),
          deliveryDelay: Duration.seconds(10),
          retentionPeriod: Duration.seconds(60),
        },
      },
    }
  );

  queueStore.attachPermissions("*");
  queueStore.attachPermissions([
    "sqs",
    "dynamodb",
  ]);

}

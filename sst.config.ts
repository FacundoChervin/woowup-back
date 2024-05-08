import { SSTConfig } from "sst";
import { HTTPStack } from "./stacks/HTTPStack";
import { QueueStore } from "./stacks/QueueStack";

export default {
  config(_input) {
    return {
      name: "email-sender-woowup",
      region: "us-west-2",
    };
  },
  stacks(app) {
    app.setDefaultFunctionProps({
      runtime: "go",
      environment: {
        "API_KEY": process.env.API_KEY || "",
        "SQL_DRIVER": process.env.SQL_DRIVER || "",
        "SQL_PASSWORD": process.env.SQL_PASSWORD || "",
        "SQL_USERNAME": process.env.SQL_USERNAME || "",
        "SQL_HOST": process.env.SQL_HOST || "",
        "SQL_PORT": process.env.SQL_PORT || "",
        "LAMBDA_AWS_ACCESS_KEY": process.env.LAMBDA_AWS_ACCESS_KEY || "",
        "LAMBDA_AWS_SECRET_KEY": process.env.LAMBDA_AWS_SECRET_KEY || "",
        "LAMBDA_AWS_REGION": process.env.LAMBDA_AWS_REGION || "",
        "LAMBDA_AWS_DYNAMO_TABLE": process.env.LAMBDA_AWS_DYNAMO_TABLE || "",
        "MAILGUN_API_KEY": process.env.MAILGUN_API_KEY || "",
        "SENDGRID_API_KEY": process.env.SENDGRID_API_KEY || "",
        "TWILIO_API_KEY": process.env.TWILIO_API_KEY || "",
        "SENDGRID_HOST": process.env.SENDGRID_HOST || "",
        "EMAIL_DOMAIN": process.env.EMAIL_DOMAIN || "",
        "SEND_SQS_QUEUE_NAME": process.env.SEND_SQS_QUEUE_NAME || "",
        "EMAIL_PROVIDERS_ORDER": process.env.EMAIL_PROVIDERS_ORDER || "",
      }
    });
    app.stack(HTTPStack);
    app.stack(QueueStore);
  },
} satisfies SSTConfig;

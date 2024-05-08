import { StackContext, Api, Queue } from "sst/constructs";
import { QueueStore } from "./QueueStack";
import { Duration } from "aws-cdk-lib";

export function HTTPStack({ stack }: StackContext) {
	const api = new Api(stack, "api", {
		routes: {
			"POST /emails/send": { function: { handler: "functions/lambda/send/main.go" } },
			"GET /emails/{id}": { function: { handler: "functions/lambda/get/main.go" } }
		},

	});

	api.attachPermissions("*");
	api.attachPermissions([
		"sqs",
		"dynamodb",
	]);

	stack.addOutputs({
		ApiEndpoint: api.url,
	});

}

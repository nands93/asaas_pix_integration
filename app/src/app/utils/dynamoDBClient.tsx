import { DynamoDBClient } from "@aws-sdk/client-dynamodb";
import { DynamoDBDocumentClient } from "@aws-sdk/lib-dynamodb";

export function getDynamoDBClient() {
    const client = new DynamoDBClient({ region: "sa-east-1" });
    return DynamoDBDocumentClient.from(client);
}
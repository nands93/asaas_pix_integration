import { NextResponse } from 'next/server';
import { DynamoDBClient } from "@aws-sdk/client-dynamodb";
import { DynamoDBDocumentClient, ScanCommand } from "@aws-sdk/lib-dynamodb";

const client = new DynamoDBClient({ region: "sa-east-1" });
const ddbDocClient = DynamoDBDocumentClient.from(client);

export async function GET() {
    try {
        const data = await ddbDocClient.send(new ScanCommand({
            TableName: "AsaasWebhooks"
        }));
        return NextResponse.json(data.Items);
    } catch (error) {
        console.error("Erro ao buscar webhooks do DynamoDB:", error);
        return NextResponse.json({ message: "Erro ao buscar webhooks" }, { status: 500 });
    }
}

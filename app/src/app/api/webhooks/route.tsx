import { NextRequest, NextResponse } from 'next/server';
import { DynamoDBClient } from "@aws-sdk/client-dynamodb";
import { DynamoDBDocumentClient, PutCommand } from "@aws-sdk/lib-dynamodb";
import dotenv from 'dotenv';

dotenv.config();

const client = new DynamoDBClient({ region: "sa-east-1" });
const ddbDocClient = DynamoDBDocumentClient.from(client);

type PaymentEvent = {
  id: string;
  event: string;
  payment: {
    id: string;
    amount: number;
  };
};

async function saveToDynamoDB(item: PaymentEvent) {
  await ddbDocClient.send(new PutCommand({
    TableName: "TransacoesAsaas",
    Item: item,
  }));
}

async function createPayment(payment: PaymentEvent['payment']) {
  console.log('Criando pagamento:', payment);
}

async function receivePayment(payment: PaymentEvent['payment']) {
  console.log('Recebendo pagamento:', payment);
}

export async function POST(request: NextRequest) {
  try {
    const body: PaymentEvent = await request.json();

    if (!body.id || !body.event) {
      return NextResponse.json({ message: 'Missing required fields' }, { status: 400 });
    }

    switch (body.event) {
      case 'PAYMENT_CREATED':
        await createPayment(body.payment);
        break;
      case 'PAYMENT_RECEIVED':
        await receivePayment(body.payment);
        break;
      default:
        return NextResponse.json({ message: 'Evento desconhecido' }, { status: 400 });
    }

    await saveToDynamoDB(body);

    return NextResponse.json({ message: 'Item successfully created', item: body });
  } catch (error) {
    console.error("Erro ao salvar webhook no DynamoDB:", error);
    return NextResponse.json({ message: "Erro ao salvar webhook" }, { status: 500 });
  }
}

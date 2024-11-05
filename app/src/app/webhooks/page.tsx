"use client";
import { useEffect, useState } from 'react';

interface Webhook {
    id: string;
    event: string;
    payment?: {
        value: number;
        status: string;
    };
}

export default function WebhooksPage() {
    const [webhooks, setWebhooks] = useState<Webhook[]>([]);

    useEffect(() => {
        async function fetchWebhooks() {
            const response = await fetch('/api/webhooks');
            const data = await response.json();
            setWebhooks(data);
        }
        fetchWebhooks();
    }, []);

    return (
        <div>
            <h1>Webhooks Recebidos</h1>
            <ul>
                {webhooks.map((webhook) => (
                    <li key={webhook.id}>
                        <p>Evento: {webhook.event}</p>
                        <p>Valor: {webhook.payment?.value}</p>
                        <p>Status: {webhook.payment?.status}</p>
                        <hr />
                    </li>
                ))}
            </ul>
        </div>
    );
}

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
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        async function fetchWebhooks() {
            try {
                const response = await fetch('/api/webhooks');
                if (!response.ok) {
                    throw new Error(`HTTP error! status: ${response.status}`);
                }
                const data = await response.json();
                setWebhooks(data.items);
            } catch (error) {
                setError((error as Error).message);
            }
        }
        fetchWebhooks();
    }, []);

    return (
        <div>
            <h1>Webhooks Recebidos</h1>
            {error ? (
                <p>Erro ao carregar webhooks: {error}</p>
            ) : (
                <ul>
                    {webhooks.map((webhook) => (
                        <li key={webhook.id}>
                            <p>Evento: {webhook.event}</p>
                            {webhook.payment && (
                                <>
                                    <p>Valor: {webhook.payment.value}</p>
                                    <p>Status: {webhook.payment.status}</p>
                                </>
                            )}
                            <hr />
                        </li>
                    ))}
                </ul>
            )}
        </div>
    );
}
import type {
	CreateMarketRequest,
	CreateMarketResponse,
	MarketCreateContextResponse
} from '$lib/api-types';
import { fetchBackend } from '$lib/backend-api';

export async function loadCreateMarketContext(): Promise<MarketCreateContextResponse> {
	const response = await fetchBackend('/api/v1/markets/context');

	if (!response.ok) {
		throw new Error('Unable to load market options right now.');
	}

	return (await response.json()) as MarketCreateContextResponse;
}

export async function submitCreateMarket(input: CreateMarketRequest): Promise<CreateMarketResponse> {
	const response = await fetchBackend('/api/v1/markets', {
		method: 'POST',
		headers: {
			'Content-Type': 'application/json'
		},
		body: JSON.stringify(input)
	});

	if (!response.ok) {
		const fallbackMessage = 'Unable to create the market right now.';

		try {
			const payload = (await response.json()) as { error?: string };
			throw new Error(payload.error ?? fallbackMessage);
		} catch (error) {
			if (error instanceof Error) {
				throw error;
			}

			throw new Error(fallbackMessage);
		}
	}

	return (await response.json()) as CreateMarketResponse;
}

import type {
	BalanceResponse,
	CreatePositionRequest,
	CreatePositionResponse,
	ListMarketsResponse,
	ListPositionsResponse,
	MarketResponse,
	PositionResponse
} from '$lib/api-types';
import { fetchBackend } from '$lib/backend-api';

export async function loadMarketDetail(marketID: string): Promise<MarketResponse> {
	const response = await fetchBackend(`/api/v1/markets/${marketID}`);

	if (!response.ok) {
		throw new Error(response.status === 404 ? 'That market could not be found.' : 'Unable to load this market right now.');
	}

	return (await response.json()) as MarketResponse;
}

export async function loadRelatedMarkets(marketID: string): Promise<MarketResponse[]> {
	const response = await fetchBackend('/api/v1/markets');

	if (!response.ok) {
		throw new Error('Unable to load related markets right now.');
	}

	const payload = (await response.json()) as ListMarketsResponse;

	return payload.active.filter((market) => market.id !== marketID).slice(0, 3);
}

export async function loadMarketAccountData(): Promise<{
	balance: BalanceResponse;
	positions: PositionResponse[];
}> {
	const [balanceResponse, positionsResponse] = await Promise.all([
		fetchBackend('/api/v1/players/me/balance'),
		fetchBackend('/api/v1/players/me/positions')
	]);

	if (balanceResponse.status === 401 || positionsResponse.status === 401) {
		throw new Error('Your guest session is not ready yet.');
	}

	if (!balanceResponse.ok) {
		throw new Error('Unable to load your balance right now.');
	}

	if (!positionsResponse.ok) {
		throw new Error('Unable to load your positions right now.');
	}

	return {
		balance: (await balanceResponse.json()) as BalanceResponse,
		positions: ((await positionsResponse.json()) as ListPositionsResponse).positions
	};
}

export async function submitPosition(
	marketID: string,
	input: CreatePositionRequest
): Promise<CreatePositionResponse> {
	const response = await fetchBackend(`/api/v1/markets/${marketID}/positions`, {
		method: 'POST',
		headers: {
			'Content-Type': 'application/json'
		},
		body: JSON.stringify(input)
	});

	if (!response.ok) {
		const fallbackMessage = 'Unable to place that position right now.';

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

	return (await response.json()) as CreatePositionResponse;
}

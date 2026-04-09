import { env } from '$env/dynamic/public';

export function buildBackendUrl(path: string): string {
	const trimmedBaseUrl = (env.PUBLIC_API_BASE_URL ?? '').trim();

	if (!trimmedBaseUrl) {
		throw new Error('PUBLIC_API_BASE_URL is not set.');
	}

	const normalizedBaseUrl = trimmedBaseUrl.endsWith('/') ? trimmedBaseUrl : `${trimmedBaseUrl}/`;
	const normalizedPath = path.startsWith('/') ? path.slice(1) : path;

	return new URL(normalizedPath, normalizedBaseUrl).toString();
}

export function fetchBackend(path: string, init: RequestInit = {}): Promise<Response> {
	return fetch(buildBackendUrl(path), {
		...init,
		credentials: 'include'
	});
}

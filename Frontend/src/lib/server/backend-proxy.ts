import { env as privateEnv } from '$env/dynamic/private';
import { env as publicEnv } from '$env/dynamic/public';

export function buildBackendProxyUrl(path: string, search: string): string {
	const baseUrl = (privateEnv.API_BASE_URL ?? publicEnv.PUBLIC_API_BASE_URL ?? '').trim();

	if (!baseUrl) {
		throw new Error('API_BASE_URL is not set.');
	}

	const normalizedBaseUrl = baseUrl.endsWith('/') ? baseUrl : `${baseUrl}/`;
	const normalizedPath = path.startsWith('/') ? path.slice(1) : path;
	const url = new URL(normalizedPath, normalizedBaseUrl);

	if (search) {
		url.search = search;
	}

	return url.toString();
}

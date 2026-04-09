export function buildApiPath(path: string): string {
	return path.startsWith('/') ? path : `/${path}`;
}

export function fetchBackend(path: string, init: RequestInit = {}): Promise<Response> {
	return fetch(buildApiPath(path), {
		...init,
		credentials: 'include'
	});
}

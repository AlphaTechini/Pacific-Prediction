import type { RequestHandler } from './$types';

import { buildBackendProxyUrl } from '$lib/server/backend-proxy';

const HOP_BY_HOP_HEADERS = new Set([
	'connection',
	'content-length',
	'host',
	'keep-alive',
	'proxy-authenticate',
	'proxy-authorization',
	'te',
	'trailer',
	'transfer-encoding',
	'upgrade'
]);

export const GET: RequestHandler = (event) => proxyRequest(event);
export const POST: RequestHandler = (event) => proxyRequest(event);
export const PUT: RequestHandler = (event) => proxyRequest(event);
export const PATCH: RequestHandler = (event) => proxyRequest(event);
export const DELETE: RequestHandler = (event) => proxyRequest(event);
export const OPTIONS: RequestHandler = (event) => proxyRequest(event);

async function proxyRequest({ fetch, params, request, url }: Parameters<RequestHandler>[0]) {
	const targetPath = params.path ? `/api/${params.path}` : '/api';
	const targetUrl = buildBackendProxyUrl(targetPath, url.search);
	const headers = new Headers(request.headers);

	for (const header of HOP_BY_HOP_HEADERS) {
		headers.delete(header);
	}

	const init: RequestInit = {
		method: request.method,
		headers
	};

	if (request.method !== 'GET' && request.method !== 'HEAD') {
		init.body = await request.arrayBuffer();
	}

	const upstreamResponse = await fetch(targetUrl, init);
	const responseHeaders = new Headers(upstreamResponse.headers);

	for (const header of HOP_BY_HOP_HEADERS) {
		responseHeaders.delete(header);
	}

	return new Response(upstreamResponse.body, {
		status: upstreamResponse.status,
		statusText: upstreamResponse.statusText,
		headers: responseHeaders
	});
}

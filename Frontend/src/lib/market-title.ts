export function formatMarketTitleDisplay(title: string): string {
	return title.replace(
		/(\d{4}-\d{2}-\d{2})T(\d{2}:\d{2}(?::\d{2})?(?:\.\d+)?(?:Z|[+-]\d{2}:\d{2})?)/g,
		'$1 $2'
	);
}

export function formatMarketQuestionTiming(value: string): string {
	return formatMarketTitleDisplay(value.trim());
}

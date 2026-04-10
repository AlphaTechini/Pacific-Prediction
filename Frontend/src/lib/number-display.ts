export function formatAmount(value?: string | null): string {
	if (value === null || value === undefined) {
		return '0';
	}

	const trimmed = String(value).trim();
	if (trimmed === '') {
		return '0';
	}

	const negative = trimmed.startsWith('-');
	const unsigned = negative ? trimmed.slice(1) : trimmed;
	const [wholePartRaw, fractionalPartRaw = ''] = unsigned.split('.', 2);
	const wholePart = wholePartRaw === '' ? '0' : wholePartRaw;
	const fractionalPart = fractionalPartRaw.replace(/0+$/, '');
	const normalized = fractionalPart === '' ? wholePart : `${wholePart}.${fractionalPart}`;

	return negative ? `-${normalized}` : normalized;
}

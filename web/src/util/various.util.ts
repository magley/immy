export const AddRangeNoDuplicates = <T>(target: T[], range: T[]): T[] => {
	return range.reduce((acc, item) => {return acc.includes(item) ? acc : [...acc, item] }, [...target]);
}

export const GetElipsisString = (input: string, lenIncludingDots: number = 10 + 3) => {
	if (!input) return input;

	if (input.length > lenIncludingDots) {
		return `${input.substring(0, lenIncludingDots)}...`;
	} else {
		return input;
	}
}

export const IsAlphaNumeric = (str: string) =>  {
	var code, i, len;

	for (i = 0, len = str.length; i < len; i++) {
		code = str.charCodeAt(i);
		if (!(code > 47 && code < 58) && !(code > 64 && code < 91) && !(code > 96 && code < 123)) {
			return false;
		}
	}
	return true;
}

export const StripSlashes = (str: string) => {
	if (str.startsWith("/") && str.endsWith("/")) {
		return str.slice(1, -1);
	}
	return str;
}
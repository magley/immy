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

export const GetPublicIdColorBackground = (user_id: string): string => {
	var hash = 0;
	if (user_id.length === 0) {
		return "#FFFFFF";
	}
	for (var i = 0; i < user_id.length; i++) {
		hash = user_id.charCodeAt(i) + ((hash << 5) - hash);
		hash = hash & hash;
	}

	var color = '#';
	for (var i = 0; i < 3; i++) {
		var value = (hash >> (i * 8)) & 255;
		color += ('00' + value.toString(16)).substr(-2);
	}

	return color;
}

export const GetPublicIdColorForeground = (user_id: string): string => {
	const bgIsLIght: boolean = IsHexColorLight(GetPublicIdColorBackground(user_id));

	if (bgIsLIght) {
		return "#000000";
	} else {
		return "#FFFFFF";
	}
}

export const IsHexColorLight = (color: string): boolean => {
	if (color.length == 6) {
		return IsHexColorLight("#" + color);
	}

	if (color.length == 7) {
		const rgb: number[] = [
			parseInt(color.substring(1, 3), 16),
			parseInt(color.substring(3, 5), 16),
			parseInt(color.substring(5), 16),
		];
		const luminance =
		(0.2126 * rgb[0]!) / 255 +
		(0.7152 * rgb[1]!) / 255 +
		(0.0722 * rgb[2]!) / 255;
		return luminance > 0.5;
	}
	return false
}